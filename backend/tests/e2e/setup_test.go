//go:build e2e

package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"

	"backend/cmd/app/config"
	adapterhttp "backend/internal/adapter/in/http"
	adapterpg "backend/internal/adapter/out/postgres"
	adapteryookassa "backend/internal/adapter/out/yookassa"
	"backend/internal/app/usecase"
	infrajwt "backend/internal/infrastructure/jwt"
	infrapass "backend/internal/infrastructure/password"
	infraqrcode "backend/internal/infrastructure/qrcode"
	"backend/internal/testhelpers"
	pkgpostgres "backend/pkg/postgres"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/stretchr/testify/require"
)

const (
	migrationVersion = 7 // Убедись, что тут актуальная версия
	apiVersion       = "v1"
)

type TestApp struct {
	Server   *httptest.Server
	Client   *http.Client
	Pg       *testhelpers.PostgresContainer
	DBClient *pkgpostgres.Client
	Cfg      *config.Config
}

var (
	appInstance *TestApp
	once        sync.Once
)

func parseLogLevel(level string) slog.Level {
	switch level {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func newLogger(level string) *slog.Logger {
	logLevel := parseLogLevel(level)
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))
}

func newPostgresClient(ctx context.Context, cfg *config.Config) (*pkgpostgres.Client, error) {
	pgConfig := pkgpostgres.NewConfig(
		cfg.DbHost, cfg.DbPort, cfg.DbUser, cfg.DbPassword,
		cfg.DBName, cfg.DbSSLMode, cfg.DbMaxConn,
		cfg.DbMinConn, cfg.DbMaxConnLifeTime, cfg.DbMaxConnIdleTime,
	)

	return pkgpostgres.NewClient(ctx, pgConfig)
}

func seedAdmins(
	ctx context.Context,
	cfg *config.Config,
	adminRepo *adapterpg.AdminRepository,
	passHasher *infrapass.Hasher,
) error {
	for _, seed := range cfg.GetAdminSeeds() {
		hash, err := passHasher.Hash(seed.Password)
		if err != nil {
			return fmt.Errorf("failed to hash password for login %s: %w", seed.Login, err)
		}

		// Поскольку тут UPSERT, метод безопасно отработает как на пустой, так и на заполненной базе
		err = adminRepo.EnsureAdmin(ctx, seed.Login, hash)
		if err != nil {
			return fmt.Errorf("failed to ensure admin for login %s: %w", seed.Login, err)
		}
	}
	return nil
}

func setupE2E(t *testing.T) *TestApp {
	once.Do(func() {
		ctx := context.Background()

		os.Setenv("DB_HOST", "test_host")
		os.Setenv("DB_USER", "test_user")
		os.Setenv("DB_PASSWORD", "test_pass")
		os.Setenv("DB_NAME", "test_db")

		os.Setenv("AUTH_SECRET", "super-secret-key-for-tests-32-chars!")
		os.Setenv("AUTH_TTL", "24h")
		os.Setenv("ADMIN_SEEDS", "test:test123")

		os.Setenv("YOOKASSA_SHOP_ID", "test_shop")
		os.Setenv("YOOKASSA_API_KEY", "test_api_key")
		os.Setenv("QR_CODE_BASE_URL", "http://localhost:8080")

		// 1. Загружаем базовый конфиг
		cfg, err := config.Load()
		require.NoError(t, err)

		// 2. Поднимаем контейнер, передавая данные из конфига (обновленный хелпер)
		pg, err := testhelpers.StartPostgresContainer(ctx)
		require.NoError(t, err)

		// Накатываем миграции в первый раз
		err = pg.MigrateUp(migrationVersion)
		require.NoError(t, err)

		// 3. ВАЖНО: Подменяем хост и порт в конфиге на динамические из Testcontainers
		cfg.DbUser = pg.Config.User
		cfg.DbPassword = pg.Config.Password
		cfg.DBName = pg.Config.Name
		cfg.DbHost = pg.Config.Host
		cfg.DbPort = pg.Config.Port

		logger := newLogger(cfg.LogLevel)

		// 4. Инициализируем клиент БД для приложения
		pgClient, err := newPostgresClient(ctx, cfg)
		require.NoError(t, err)

		// Транзакционный менеджер
		trManager := manager.Must(trmpgx.NewDefaultFactory(pgClient.Pool))

		// Репозитории
		adminRepo := adapterpg.NewAdminRepository(pgClient, trmpgx.DefaultCtxGetter)
		locationRepo := adapterpg.NewLocationRepository(pgClient, trmpgx.DefaultCtxGetter)
		itemRepo := adapterpg.NewItemRepository(pgClient, trmpgx.DefaultCtxGetter)
		locationItemRepo := adapterpg.NewLocationItemRepository(pgClient, trmpgx.DefaultCtxGetter)
		orderRepo := adapterpg.NewOrderRepository(pgClient, trmpgx.DefaultCtxGetter)
		orderItemRepo := adapterpg.NewOrderItemRepository(pgClient, trmpgx.DefaultCtxGetter)
		transactionRepo := adapterpg.NewTransactionRepository(pgClient, trmpgx.DefaultCtxGetter)

		// Инфраструктура
		tokenGen := infrajwt.NewGenerator(cfg.AuthSecret, cfg.AuthTTL)
		passHasher := infrapass.NewHasher(cfg.PasswordCost)
		qrCodeGen := infraqrcode.NewGenerator(cfg.QRCodeBaseURL, cfg.QRCodeSize)
		paymentGateway := adapteryookassa.NewPaymentGateway(cfg.YookassaShopID, cfg.YookassaAPIKey, cfg.YookassaTimeout)

		// Первичное заполнение админов
		err = seedAdmins(ctx, cfg, adminRepo, passHasher)
		require.NoError(t, err)

		// UseCases
		adminAuthUC := usecase.NewAdminAuthUC(adminRepo, passHasher, tokenGen)
		createLocationUC := usecase.NewCreateLocationUC(trManager, locationRepo, itemRepo, locationItemRepo)
		updateLocationUC := usecase.NewUpdateLocationUC(locationRepo)
		deleteLocationUC := usecase.NewDeleteLocationUC(trManager, locationRepo, locationItemRepo)
		listLocationsUC := usecase.NewListLocationsUC(locationRepo)
		getQRCodeUC := usecase.NewGetQRCodeUC(locationRepo, qrCodeGen)
		getCatalogUC := usecase.NewGetCatalogUC(locationRepo, itemRepo, locationItemRepo)
		createItemUC := usecase.NewCreateItemUC(trManager, locationRepo, itemRepo, locationItemRepo)
		updateItemUC := usecase.NewUpdateItemUC(itemRepo)
		deleteItemUC := usecase.NewDeleteItemUC(trManager, itemRepo, locationItemRepo)
		listItemsUC := usecase.NewListItemsUC(itemRepo)
		createOrderUC := usecase.NewCreateOrderUC(trManager, locationRepo, locationItemRepo, orderRepo, orderItemRepo, transactionRepo, paymentGateway)
		getInventoryUC := usecase.NewGetInventoryUC(locationRepo, locationItemRepo)
		updateInventoryUC := usecase.NewUpdateInventoryUC(trManager, locationRepo, locationItemRepo)

		// Handlers
		systemHandler := adapterhttp.NewSystemHandler(cfg.Environment, apiVersion)
		authHandler := adapterhttp.NewAuthHandler(logger, adminAuthUC)
		clientHandler := adapterhttp.NewClientHandler(logger, getCatalogUC, createOrderUC)
		locationsHandler := adapterhttp.NewLocationHandler(logger, createLocationUC, updateLocationUC, deleteLocationUC, listLocationsUC, getQRCodeUC)
		itemHandler := adapterhttp.NewItemHandler(logger, createItemUC, updateItemUC, deleteItemUC, listItemsUC)
		inventoryHandler := adapterhttp.NewInventoryHandler(logger, getInventoryUC, updateInventoryUC)

		// Роутер
		router := adapterhttp.NewRouter(
			tokenGen, systemHandler, authHandler, clientHandler,
			locationsHandler, itemHandler, inventoryHandler,
		).InitRoutes()

		// Запускаем тестовый сервер
		ts := httptest.NewServer(router)

		appInstance = &TestApp{
			Server:   ts,
			Client:   ts.Client(),
			Pg:       pg,
			DBClient: pgClient,
			Cfg:      cfg,
		}
	})

	// Очищаем и подготавливаем базу ПЕРЕД каждым тестом
	appInstance.cleanData(t, context.Background())

	return appInstance
}

// cleanData использует инструменты миграций для сброса БД в чистый вид
// и восстанавливает необходимые seed-данные.
func (a *TestApp) cleanData(t *testing.T, ctx context.Context) {
	// 1. Сносим все таблицы
	err := a.Pg.MigrateDown()
	require.NoError(t, err, "failed to migrate down")

	// 2. Накатываем структуру заново
	err = a.Pg.MigrateUp(migrationVersion)
	require.NoError(t, err, "failed to migrate up")

	// 3. Заново создаем админов, так как MigrateDown удалил таблицу admins
	adminRepo := adapterpg.NewAdminRepository(a.DBClient, trmpgx.DefaultCtxGetter)
	passHasher := infrapass.NewHasher(a.Cfg.PasswordCost)
	err = seedAdmins(ctx, a.Cfg, adminRepo, passHasher)
	require.NoError(t, err, "failed to re-seed admins")
}

func (a *TestApp) DoRequest(method, path string, body interface{}) (*http.Response, error) {
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, a.Server.URL+path, &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	return a.Client.Do(req)
}

// TearDownE2E можно вызвать при завершении пакета тестов, если нужно
// принудительно освободить ресурсы (обычно go test сам все убивает при выходе)
func TearDownE2E() {
	if appInstance != nil {
		appInstance.Server.Close()
		appInstance.DBClient.Close()
		_ = appInstance.Pg.Close(context.Background())
	}
}
