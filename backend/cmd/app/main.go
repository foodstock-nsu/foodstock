package main

import (
	"backend/cmd/app/config"
	adapterhttp "backend/internal/adapter/in/http"
	adapterpg "backend/internal/adapter/out/postgres"
	"backend/internal/app/service"
	"backend/internal/app/usecase"
	infrajwt "backend/internal/infrastructure/jwt"
	infrapass "backend/internal/infrastructure/password"
	infraqrcode "backend/internal/infrastructure/qrcode"
	pkgpostgres "backend/pkg/postgres"
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
)

const (
	apiVersion      = "v1"
	shutdownTimeout = 10 * time.Second
	cleanupDelay    = 10 * time.Minute
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

	pgClient, err := pkgpostgres.NewClient(ctx, pgConfig)
	if err != nil {
		return nil, err
	}

	return pgClient, nil
}

func closePostgresClient(
	ctx context.Context,
	logger *slog.Logger,
	pgClient *pkgpostgres.Client,
) {
	logger.InfoContext(ctx, "closing postgres connection...")
	pgClient.Close()
}

// Fill database with seed data
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

		err = adminRepo.EnsureAdmin(ctx, seed.Login, hash)
		if err != nil {
			return fmt.Errorf("failed to ensure admin for login %s: %w", seed.Login, err)
		}
	}
	return nil
}

func runServer(ctx context.Context, cfg *config.Config, logger *slog.Logger) error {
	// Postgres client
	pgClient, err := newPostgresClient(ctx, cfg)
	if err != nil {
		return fmt.Errorf("failed to init postgres client: %w", err)
	}

	// Close Postgres
	defer closePostgresClient(ctx, logger, pgClient)

	// Transaction manager
	trManager := manager.Must(trmpgx.NewDefaultFactory(pgClient.Pool))

	// Repositories
	adminRepo := adapterpg.NewAdminRepository(pgClient, trmpgx.DefaultCtxGetter)
	locationRepo := adapterpg.NewLocationRepository(pgClient, trmpgx.DefaultCtxGetter)
	itemRepo := adapterpg.NewItemRepository(pgClient, trmpgx.DefaultCtxGetter)
	locationItemRepo := adapterpg.NewLocationItemRepository(pgClient, trmpgx.DefaultCtxGetter)
	orderRepo := adapterpg.NewOrderRepository(pgClient, trmpgx.DefaultCtxGetter)
	orderItemRepo := adapterpg.NewOrderItemRepository(pgClient, trmpgx.DefaultCtxGetter)

	// Infrastructure
	tokenGen := infrajwt.NewGenerator(cfg.AuthSecret, cfg.AuthTTL)
	passHasher := infrapass.NewHasher(cfg.PasswordCost)
	qrCodeGen := infraqrcode.NewGenerator(cfg.QRCodeBaseURL, cfg.QRCodeSize)

	// Fill database with seed data
	err = seedAdmins(ctx, cfg, adminRepo, passHasher)
	if err != nil {
		return fmt.Errorf("failed to add seed data: %w", err)
	}

	// UseСases
	adminAuthUC := usecase.NewAdminAuthUC(adminRepo, passHasher, tokenGen)
	createLocationUC := usecase.NewCreateLocationUC(
		trManager, locationRepo, itemRepo, locationItemRepo,
	)
	updateLocationUC := usecase.NewUpdateLocationUC(locationRepo)
	deleteLocationUC := usecase.NewDeleteLocationUC(
		trManager, locationRepo, locationItemRepo,
	)
	listLocationsUC := usecase.NewListLocationsUC(locationRepo)
	getQRCodeUC := usecase.NewGetQRCodeUC(locationRepo, qrCodeGen)
	getCatalogUC := usecase.NewGetCatalogUC(itemRepo, locationItemRepo)
	createItemUC := usecase.NewCreateItemUC(
		trManager, locationRepo, itemRepo, locationItemRepo,
	)
	updateItemUC := usecase.NewUpdateItemUC(itemRepo)
	deleteItemUC := usecase.NewDeleteItemUC(
		trManager, itemRepo, locationItemRepo,
	)
	listItemsUC := usecase.NewListItemsUC(itemRepo)

	// Services
	orderCleaner := service.NewOrderCleaner(
		trManager, locationItemRepo, orderRepo, orderItemRepo)

	// Handlers
	systemHandler := adapterhttp.NewSystemHandler(cfg.Environment, apiVersion)
	authHandler := adapterhttp.NewAuthHandler(logger, adminAuthUC)
	clientHandler := adapterhttp.NewClientHandler(logger, getCatalogUC)
	locationsHandler := adapterhttp.NewLocationHandler(
		logger, createLocationUC, updateLocationUC,
		deleteLocationUC, listLocationsUC, getQRCodeUC,
	)
	itemHandler := adapterhttp.NewItemHandler(
		logger, createItemUC,
		updateItemUC, deleteItemUC, listItemsUC,
	)

	// Router
	router := adapterhttp.NewRouter(
		tokenGen,
		systemHandler,
		authHandler,
		clientHandler,
		locationsHandler,
		itemHandler,
	).InitRoutes()

	// Launch background job - cleaning expired orders
	go cleanupExpiredOrders(ctx, logger, orderCleaner)

	// Launch server with graceful shutdown
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.HttpPort),
		Handler: router,
	}

	errCh := make(chan error, 1)

	go func() {
		logger.Info("starting server", slog.String("address", fmt.Sprintf(":%d", cfg.HttpPort)))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
			return
		}
		errCh <- nil
	}()

	select {
	case <-ctx.Done():
		logger.Info("shutdown signal received")
	case err := <-errCh:
		if err != nil {
			logger.Error("server failed", slog.Any("err", err))
			return err
		}
		logger.Info("server stopped")
		return nil
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("graceful shutdown failed", slog.Any("err", err))
		_ = srv.Close() // fallback
		return err
	}

	logger.Info("server exited properly")
	return nil
}

func cleanupExpiredOrders(
	ctx context.Context,
	logger *slog.Logger,
	cleaner *service.OrderCleaner,
) {
	time.Sleep(20 * time.Second) // Wait until the server wakes up

	ticker := time.NewTicker(cleanupDelay)
	defer ticker.Stop()

	logger.InfoContext(ctx, "background worker started", slog.Duration("delay", cleanupDelay))

	for {
		select {
		case <-ctx.Done():
			logger.Info("stopping background cleanup worker...")
			return
		case <-ticker.C:
			logger.InfoContext(ctx, "Background job: cleaning expired orders...")
			if err := cleaner.Cleanup(ctx); err != nil {
				logger.ErrorContext(ctx, "Background job errors", slog.Any("err", err))
			}
		}
	}
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	logger := newLogger(cfg.LogLevel)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := runServer(ctx, cfg, logger); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
