//go:build integration

package postgres_test

import (
	adapterpostgres "backend/internal/adapter/out/postgres"
	"backend/internal/domain/model"
	"backend/internal/testhelpers"
	pkgerrs "backend/pkg/errs"
	pkgpostgres "backend/pkg/postgres"
	"backend/pkg/utils"
	"context"
	"testing"
	"time"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type OrderItemRepoSuite struct {
	suite.Suite
	pgContainer *testhelpers.PostgresContainer
	dbClient    *pkgpostgres.Client
	repo        *adapterpostgres.OrderItemRepository
	ctx         context.Context
	migrate     *migrate.Migrate
	testOrderID uuid.UUID
	testItemID  uuid.UUID
}

func TestOrderItemRepoSuite(t *testing.T) {
	suite.Run(t, new(OrderItemRepoSuite))
}

func (s *OrderItemRepoSuite) SetupSuite() {
	const targetVersion = 5

	ctx := context.Background()

	// Init postgres container
	pgContainer, err := testhelpers.StartPostgresContainer(ctx)
	s.Require().NoError(err)

	// Init postgres client
	client, err := pkgpostgres.NewClient(ctx, pgContainer.Config)
	s.Require().NoError(err)

	// Apply migrations
	err = pgContainer.MigrateUp(targetVersion)
	s.Require().NoError(err)

	s.pgContainer = pgContainer
	s.dbClient = client
	s.repo = adapterpostgres.NewOrderItemRepository(
		s.dbClient,
		trmpgx.DefaultCtxGetter,
	)
	s.ctx = ctx

	locID := uuid.New()
	locationRepo := adapterpostgres.NewLocationRepository(s.dbClient, trmpgx.DefaultCtxGetter)
	_ = locationRepo.Create(s.ctx, model.RestoreLocation(
		locID, "test_loc", "Shop",
		"Addr", true, time.Now().UTC(), nil,
	))

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
	_ = s.pgContainer.MigrateDown()
	s.dbClient.Close()
	_ = s.pgContainer.Close(s.ctx)
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
