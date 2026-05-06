//go:build e2e

package e2e

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"testing"
	"time"

	"backend/cmd/app/config"
	adapterhttp "backend/internal/adapter/in/http"
	adapterpg "backend/internal/adapter/out/postgres"
	adapteryookassa "backend/internal/adapter/out/yookassa"
	"backend/internal/app/service"
	"backend/internal/app/usecase"
	infrajwt "backend/internal/infrastructure/jwt"
	infrapass "backend/internal/infrastructure/password"
	infraqrcode "backend/internal/infrastructure/qrcode"
	pkgpostgres "backend/pkg/postgres"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// Глобальные переменные для всех E2E тестов.
// testApp содержит ВЕСЬ твой собранный API.
var (
	testApp http.Handler
	testDB  *pkgpostgres.Client
)

// TestMain запускается ОДИН РАЗ перед всеми тестами в папке e2e.
func TestMain(m *testing.M) {
	ctx := context.Background()

	// 1. ПОДНИМАЕМ TESTCONTAINERS С БАЗОЙ
	dbName := "foodstock_test"
	dbUser := "postgres"
	dbPass := "postgres"

	pgContainer, err := postgres.Run(ctx,
		"postgres:15-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPass),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second)),
	)
	if err != nil {
		log.Fatalf("failed to start postgres container: %v", err)
	}

	// Гарантируем, что контейнер будет убит после завершения тестов
	defer func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate container: %v", err)
		}
	}()

	host, _ := pgContainer.Host(ctx)
	port, _ := pgContainer.MappedPort(ctx, "5432")

	// 2. ПОДКЛЮЧАЕМСЯ К ТЕСТОВОЙ БАЗЕ
	pgConfig := pkgpostgres.NewConfig(
		host, port.Int(), dbUser, dbPass, dbName, "disable",
		10, 2, time.Hour, time.Minute,
	)

	testDB, err = pkgpostgres.NewClient(ctx, pgConfig)
	if err != nil {
		log.Fatalf("failed to connect to test db: %v", err)
	}
	defer testDB.Close()

	// Apply migrations
	applyTestMigrations(testDB)

	// Config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	// =====================================================================
	// ТОЧНАЯ КОПИЯ СБОРКИ ИЗ ТВОЕГО main.go (runServer)
	// =====================================================================

	// Transaction manager
	trManager := manager.Must(trmpgx.NewDefaultFactory(testDB.Pool))

	// Repositories
	adminRepo := adapterpg.NewAdminRepository(testDB, trmpgx.DefaultCtxGetter)
	locationRepo := adapterpg.NewLocationRepository(testDB, trmpgx.DefaultCtxGetter)
	itemRepo := adapterpg.NewItemRepository(testDB, trmpgx.DefaultCtxGetter)
	locationItemRepo := adapterpg.NewLocationItemRepository(testDB, trmpgx.DefaultCtxGetter)
	orderRepo := adapterpg.NewOrderRepository(testDB, trmpgx.DefaultCtxGetter)
	orderItemRepo := adapterpg.NewOrderItemRepository(testDB, trmpgx.DefaultCtxGetter)
	transactionRepo := adapterpg.NewTransactionRepository(testDB, trmpgx.DefaultCtxGetter)

	// Infrastructure
	tokenGen := infrajwt.NewGenerator(cfg.AuthSecret, cfg.AuthTTL)
	passHasher := infrapass.NewHasher(cfg.PasswordCost)
	qrCodeGen := infraqrcode.NewGenerator(cfg.QRCodeBaseURL, cfg.QRCodeSize)

	// Payment Gateway
	paymentGateway := adapteryookassa.NewPaymentGateway(
		cfg.YookassaShopID, cfg.YookassaAPIKey, cfg.YookassaTimeout,
	)

	// Fill database with seed data
	for _, seed := range cfg.GetAdminSeeds() {
		hash, hashErr := passHasher.Hash(seed.Password)
		if hashErr != nil {
			log.Fatalf("failed to hash password for login %s: %v", seed.Login, hashErr)
		}
		if err := adminRepo.EnsureAdmin(ctx, seed.Login, hash); err != nil {
			log.Fatalf("failed to ensure admin for login %s: %v", seed.Login, err)
		}
	}

	// UseСases
	adminAuthUC := usecase.NewAdminAuthUC(adminRepo, passHasher, tokenGen)
	createLocationUC := usecase.NewCreateLocationUC(trManager, locationRepo, itemRepo, locationItemRepo)
	updateLocationUC := usecase.NewUpdateLocationUC(locationRepo)
	deleteLocationUC := usecase.NewDeleteLocationUC(trManager, locationRepo, locationItemRepo)
	listLocationsUC := usecase.NewListLocationsUC(locationRepo)
	getQRCodeUC := usecase.NewGetQRCodeUC(locationRepo, qrCodeGen)
	getCatalogUC := usecase.NewGetCatalogUC(itemRepo, locationItemRepo)
	createItemUC := usecase.NewCreateItemUC(trManager, locationRepo, itemRepo, locationItemRepo)
	updateItemUC := usecase.NewUpdateItemUC(itemRepo)
	deleteItemUC := usecase.NewDeleteItemUC(trManager, itemRepo, locationItemRepo)
	listItemsUC := usecase.NewListItemsUC(itemRepo)
	createOrderUC := usecase.NewCreateOrderUC(trManager, locationRepo, locationItemRepo, orderRepo, orderItemRepo, transactionRepo, paymentGateway)
	getInventoryUC := usecase.NewGetInventoryUC(locationItemRepo)
	updateInventoryUC := usecase.NewUpdateInventoryUC(trManager, locationItemRepo)

	_ = service.NewExpirationService(trManager, locationItemRepo, orderRepo, orderItemRepo, transactionRepo)

	// Handlers
	systemHandler := adapterhttp.NewSystemHandler(cfg.Environment, "v1")
	authHandler := adapterhttp.NewAuthHandler(logger, adminAuthUC)
	clientHandler := adapterhttp.NewClientHandler(logger, getCatalogUC, createOrderUC)
	locationsHandler := adapterhttp.NewLocationHandler(logger, createLocationUC, updateLocationUC, deleteLocationUC, listLocationsUC, getQRCodeUC)
	itemHandler := adapterhttp.NewItemHandler(logger, createItemUC, updateItemUC, deleteItemUC, listItemsUC)
	inventoryHandler := adapterhttp.NewInventoryHandler(logger, getInventoryUC, updateInventoryUC)

	// Router
	testApp = adapterhttp.NewRouter(
		tokenGen,
		systemHandler,
		authHandler,
		clientHandler,
		locationsHandler,
		itemHandler,
		inventoryHandler,
	).InitRoutes()

	code := m.Run()

	os.Exit(code)
}

// applyTestMigrations создает таблицы в тестовой БД.
func applyTestMigrations(client *pkgpostgres.Client) {
	// СЮДА НУЖНО ВСТАВИТЬ SQL ДЛЯ СОЗДАНИЯ ТВОИХ ТАБЛИЦ.
	// Либо, если у тебя есть инструмент миграций (golang-migrate), вызвать его здесь.

	query := `
	CREATE TABLE IF NOT EXISTS admins (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		login TEXT NOT NULL UNIQUE,
		password_hash TEXT NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);
	-- ТУТ ДОЛЖНЫ БЫТЬ CREATE TABLE ДЛЯ items, locations, orders и т.д.
	`
	_, err := client.Pool.Exec(context.Background(), query)
	if err != nil {
		log.Fatalf("failed to apply test migrations: %v", err)
	}
}
