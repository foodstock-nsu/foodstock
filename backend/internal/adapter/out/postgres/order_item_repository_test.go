package postgres_test

import (
	adapterpostgres "backend/internal/adapter/out/postgres"
	"backend/internal/domain/model"
	"backend/migrations"
	pkgerrs "backend/pkg/errs"
	pkgpostgres "backend/pkg/postgres"
	"backend/pkg/utils"
	"context"
	"errors"
	"testing"
	"time"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type OrderItemRepoSuite struct {
	suite.Suite
	dbClient    *pkgpostgres.Client
	repo        *adapterpostgres.OrderItemRepository
	ctx         context.Context
	migrate     *migrate.Migrate
	testOrderID uuid.UUID
	testItemID  uuid.UUID
}

func TestOrderItemRepoSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}
	suite.Run(t, new(OrderItemRepoSuite))
}

func (s *OrderItemRepoSuite) setupDatabase() {
	const targetVersion = 5

	dbConfig := pkgpostgres.NewConfig(
		"localhost", 5433, "test-user",
		"test-pass", "test-db", "disable",
		5, 5,
		10*time.Second, 10*time.Second,
	)
	dsn := "postgres://test-user:test-pass@localhost:5433/test-db?sslmode=disable"

	dbClient, err := pkgpostgres.NewClient(context.Background(), dbConfig)
	s.Require().NoError(err)
	s.dbClient = dbClient

	sourceDriver, err := iofs.New(migrations.FS, ".")
	s.Require().NoError(err)

	m, err := migrate.NewWithSourceInstance("iofs", sourceDriver, dsn)
	s.Require().NoError(err)

	s.migrate = m
	err = m.Migrate(targetVersion)

	if err == nil || errors.Is(err, migrate.ErrNoChange) {
		return
	}

	var dirtyErr migrate.ErrDirty
	if errors.As(err, &dirtyErr) {
		_ = m.Force(dirtyErr.Version)
		_ = m.Down()
		err = m.Migrate(targetVersion)
		s.Require().NoError(err)
	}
}

func (s *OrderItemRepoSuite) SetupSuite() {
	s.ctx = context.Background()
	s.setupDatabase()
	s.repo = adapterpostgres.NewOrderItemRepository(
		s.dbClient,
		trmpgx.DefaultCtxGetter,
	)

	locID := uuid.New()
	locationRepo := adapterpostgres.NewLocationRepository(s.dbClient, trmpgx.DefaultCtxGetter)
	_ = locationRepo.Create(s.ctx, model.RestoreLocation(locID, "test_loc", "Shop", "Addr", true, time.Now().UTC()))

	s.testItemID = uuid.New()
	itemRepo := adapterpostgres.NewItemRepository(s.dbClient, trmpgx.DefaultCtxGetter)
	_ = itemRepo.Create(
		s.ctx,
		model.RestoreItem(
			s.testItemID,
			"Coke",
			utils.VPtr("Drink"),
			model.ItemDrinks,
			"https://example.com/photo.jpg",
			model.RestoreNutrition(
				utils.VPtr(1),
				utils.VPtr(float64(1)),
				utils.VPtr(float64(1)),
				utils.VPtr(float64(1))),
			time.Now().UTC()))

	s.testOrderID = uuid.New()
	orderRepo := adapterpostgres.NewOrderRepository(s.dbClient, trmpgx.DefaultCtxGetter)
	_ = orderRepo.Create(s.ctx, model.RestoreOrder(s.testOrderID, locID, nil, model.OrderPending, 100, time.Now().UTC(), nil))
}

func (s *OrderItemRepoSuite) TearDownSuite() {
	if s.migrate != nil {
		_ = s.migrate.Down()
	}
	s.dbClient.Close()
}

func (s *OrderItemRepoSuite) SetupTest() {
	_, err := s.dbClient.Pool.Exec(s.ctx, "TRUNCATE TABLE order_items CASCADE")
	s.Require().NoError(err)
}

func (s *OrderItemRepoSuite) TestCreateList() {
	item, _ := model.NewOrderItem(s.testItemID, 2, 5000) // 50.00

	err := s.repo.Create(s.ctx, s.testOrderID, item)
	s.Require().NoError(err)

	items, err := s.repo.List(s.ctx, s.testOrderID)
	s.Require().NoError(err)
	s.Require().Len(items, 1)
	s.Require().Equal(item.ID(), items[0].ID())
	s.Require().Equal(item.ItemID(), items[0].ItemID())
	s.Require().Equal(item.Amount(), items[0].Amount())
	s.Require().Equal(item.PriceAtPurchase(), items[0].PriceAtPurchase())
}

func (s *OrderItemRepoSuite) TestCreate_AlreadyExists() {
	item, _ := model.NewOrderItem(s.testItemID, 1, 100)

	err := s.repo.Create(s.ctx, s.testOrderID, item)
	s.Require().NoError(err)

	err = s.repo.Create(s.ctx, s.testOrderID, item)
	s.Require().Error(err)
	s.Require().ErrorIs(err, pkgerrs.ErrObjectAlreadyExists)
}

func (s *OrderItemRepoSuite) TestCreateMany() {
	item1, _ := model.NewOrderItem(s.testItemID, 1, 100)
	item2, _ := model.NewOrderItem(s.testItemID, 5, 500)
	batch := []*model.OrderItem{item1, item2}

	err := s.repo.CreateMany(s.ctx, s.testOrderID, batch)
	s.Require().NoError(err)

	items, err := s.repo.List(s.ctx, s.testOrderID)
	s.Require().NoError(err)
	s.Require().Len(items, 2)
}

func (s *OrderItemRepoSuite) TestList_Empty() {
	items, err := s.repo.List(s.ctx, uuid.New())
	s.Require().NoError(err)
	s.Require().Empty(items)
}
