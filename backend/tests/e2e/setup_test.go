///go:build e2e

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
	"strings"
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

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/stretchr/testify/require"
)

const (
	migrationVersion = 7 // Убедись, что тут актуальная версия
	apiVersion       = "v1"
)

type testApp struct {
	server     *httptest.Server
	client     *http.Client
	pg         *testhelpers.PostgresContainer
	dbClient   *pkgpostgres.Client
	cfg        *config.Config
	adminToken *string
}

var (
	appInstance *testApp
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

func setupE2E(t *testing.T) *testApp {
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
		getLocationUC := usecase.NewGetLocationUC(locationRepo)
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
		locationsHandler := adapterhttp.NewLocationHandler(logger, createLocationUC, getLocationUC, updateLocationUC, deleteLocationUC, listLocationsUC, getQRCodeUC)
		itemHandler := adapterhttp.NewItemHandler(logger, createItemUC, updateItemUC, deleteItemUC, listItemsUC)
		inventoryHandler := adapterhttp.NewInventoryHandler(logger, getInventoryUC, updateInventoryUC)

		// Роутер
		router := adapterhttp.NewRouter(
			tokenGen, systemHandler, authHandler, clientHandler,
			locationsHandler, itemHandler, inventoryHandler,
		).InitRoutes()

		// Запускаем тестовый сервер
		ts := httptest.NewServer(router)

		appInstance = &testApp{
			server:   ts,
			client:   ts.Client(),
			pg:       pg,
			dbClient: pgClient,
			cfg:      cfg,
		}
	})

	// Очищаем и подготавливаем базу ПЕРЕД каждым тестом
	appInstance.cleanData(t, context.Background())

	return appInstance
}

// cleanData использует инструменты миграций для сброса БД в чистый вид
// и восстанавливает необходимые seed-данные.
func (a *testApp) cleanData(t *testing.T, ctx context.Context) {
	// Список таблиц, которые нужно чистить (в порядке, учитывающем FK, либо используем CASCADE)
	// Добавь сюда все свои таблицы
	tables := []string{
		"order_items",
		"orders",
		"transactions",
		"location_items",
		"items",
		"locations",
		"admins",
	}

	// Формируем запрос: TRUNCATE table1, table2 RESTART IDENTITY CASCADE;
	// RESTART IDENTITY сбрасывает счетчики SERIAL/BIGSERIAL в 0
	// CASCADE удаляет зависимости, если они есть
	query := fmt.Sprintf("TRUNCATE %s RESTART IDENTITY CASCADE", strings.Join(tables, ", "))

	_, err := a.dbClient.Pool.Exec(ctx, query)
	require.NoError(t, err, "failed to truncate tables")

	// После очистки данных кэш планов Postgres (Prepared Statements)
	// НЕ ломается, так как таблицы не пересоздавались.
	// Но на всякий случай можно оставить DISCARD PLANS, он очень легкий.
	_, _ = a.dbClient.Pool.Exec(ctx, "DISCARD PLANS")

	// Заново создаем админов, так как TRUNCATE их стер
	adminRepo := adapterpg.NewAdminRepository(a.dbClient, trmpgx.DefaultCtxGetter)
	passHasher := infrapass.NewHasher(a.cfg.PasswordCost)
	err = seedAdmins(ctx, a.cfg, adminRepo, passHasher)
	require.NoError(t, err, "failed to re-seed admins")
}

func (a *testApp) doRequest(method, path string, body interface{}) (*http.Response, error) {
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, a.server.URL+path, &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	return a.client.Do(req)
}

func (a *testApp) doRequestAuth(method, path string, body interface{}, token string) (*http.Response, error) {
	var buf bytes.Buffer
	if body != nil {
		_ = json.NewEncoder(&buf).Encode(body)
	}

	req, _ := http.NewRequest(method, a.server.URL+path, &buf)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	return a.client.Do(req)
}

func (a *testApp) getAdminToken(t *testing.T) string {
	if a.adminToken != nil {
		return *a.adminToken
	}

	resp, err := a.doRequest(
		"POST",
		"/api/v1/admin/auth",
		map[string]interface{}{
			"login":    "test",
			"password": "test123",
		},
	)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var body map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&body)

	require.NoError(t, err)
	require.NotEmpty(t, body["token"])

	token, ok := body["token"].(string)
	require.True(t, ok)

	a.adminToken = &token

	return token
}

// tearDownE2E можно вызвать при завершении пакета тестов, если нужно
// принудительно освободить ресурсы (обычно go test сам все убивает при выходе)
func tearDownE2E() {
	if appInstance != nil {
		appInstance.server.Close()
		appInstance.dbClient.Close()
		_ = appInstance.pg.Close(context.Background())
	}
}
