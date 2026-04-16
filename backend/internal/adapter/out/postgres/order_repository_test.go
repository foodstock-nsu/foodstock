package postgres_test

import (
	adapterpostgres "backend/internal/adapter/out/postgres"
	"backend/internal/domain/model"
	"backend/migrations"
	pkgerrs "backend/pkg/errs"
	pkgpostgres "backend/pkg/postgres"
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

type OrderRepoSuite struct {
	suite.Suite
	dbClient  *pkgpostgres.Client
	repo      *adapterpostgres.OrderRepository
	ctx       context.Context
	migrate   *migrate.Migrate
	testOrder *model.Order
	testLocID uuid.UUID
}

func TestOrderRepoSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}
	suite.Run(t, new(OrderRepoSuite))
}

func (s *OrderRepoSuite) setupDatabase() {
	const targetVersion = 4

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

func (s *OrderRepoSuite) SetupSuite() {
	s.ctx = context.Background()
	s.setupDatabase()
	s.repo = adapterpostgres.NewOrderRepository(
		s.dbClient,
		trmpgx.DefaultCtxGetter,
	)

	s.testLocID = uuid.New()
	locationRepo := adapterpostgres.NewLocationRepository(
		s.dbClient,
		trmpgx.DefaultCtxGetter,
	)

	_ = locationRepo.Create(
		s.ctx,
		model.RestoreLocation(
			s.testLocID,
			"nsu_order_test",
			"NSU Test Shop",
			"Pirogova 1",
			true,
			time.Now().UTC(),
		),
	)

	s.testOrder = model.RestoreOrder(
		uuid.New(),
		s.testLocID,
		nil,
		model.OrderPending,
		50000, // 500.00
		time.Now().UTC().Truncate(time.Microsecond),
		nil,
	)
}

func (s *OrderRepoSuite) TearDownSuite() {
	if s.migrate != nil {
		_ = s.migrate.Down()
	}
	s.dbClient.Close()
}

func (s *OrderRepoSuite) SetupTest() {
	_, err := s.dbClient.Pool.Exec(s.ctx, "TRUNCATE TABLE orders CASCADE")
	s.Require().NoError(err)
}

func (s *OrderRepoSuite) TestCreateGet() {
	err := s.repo.Create(s.ctx, s.testOrder)
	s.Require().NoError(err)

	order, err := s.repo.Get(s.ctx, s.testOrder.ID())
	s.Require().NoError(err)
	s.Require().NotNil(order)
	s.Require().Equal(s.testOrder.ID(), order.ID())
	s.Require().Equal(s.testOrder.LocationID(), order.LocationID())
	s.Require().Equal(s.testOrder.Status(), order.Status())
	s.Require().Equal(s.testOrder.TotalPrice(), order.TotalPrice())
	s.Require().Nil(order.PaidAt())
}

func (s *OrderRepoSuite) TestCreate_AlreadyExists() {
	err := s.repo.Create(s.ctx, s.testOrder)
	s.Require().NoError(err)

	err = s.repo.Create(s.ctx, s.testOrder)
	s.Require().Error(err)
	s.Require().ErrorIs(err, pkgerrs.ErrObjectAlreadyExists)
}

func (s *OrderRepoSuite) TestGet_NotFound() {
	order, err := s.repo.Get(s.ctx, uuid.New())
	s.Require().Error(err)
	s.Require().ErrorIs(err, pkgerrs.ErrObjectNotFound)
	s.Require().Nil(order)
}

func (s *OrderRepoSuite) TestUpdate() {
	err := s.repo.Create(s.ctx, s.testOrder)
	s.Require().NoError(err)

	paidAt := time.Now().UTC().Truncate(time.Microsecond)
	updatedOrder := model.RestoreOrder(
		s.testOrder.ID(),
		s.testOrder.LocationID(),
		nil,
		model.OrderPaid,
		60000,
		s.testOrder.CreatedAt(),
		&paidAt,
	)

	err = s.repo.Update(s.ctx, updatedOrder)
	s.Require().NoError(err)

	res, err := s.repo.Get(s.ctx, s.testOrder.ID())
	s.Require().NoError(err)
	s.Require().Equal(model.OrderPaid, res.Status())
	s.Require().Equal(int64(60000), res.TotalPrice())
	s.Require().NotNil(res.PaidAt())
	s.Require().True(paidAt.Equal(*res.PaidAt()))
}

func (s *OrderRepoSuite) TestListByLocationID() {
	err := s.repo.Create(s.ctx, s.testOrder)
	s.Require().NoError(err)

	// Second order for same location
	order2 := model.RestoreOrder(
		uuid.New(),
		s.testLocID,
		nil,
		model.OrderPending,
		1000,
		time.Now().UTC(),
		nil,
	)
	err = s.repo.Create(s.ctx, order2)
	s.Require().NoError(err)

	orders, err := s.repo.ListByLocationID(s.ctx, s.testLocID, 10, 0)
	s.Require().NoError(err)
	s.Require().Len(orders, 2)

	// Test pagination
	ordersPaged, err := s.repo.ListByLocationID(s.ctx, s.testLocID, 1, 1)
	s.Require().NoError(err)
	s.Require().Len(ordersPaged, 1)
}

func (s *OrderRepoSuite) TestListByStatus() {
	err := s.repo.Create(s.ctx, s.testOrder) // PENDING
	s.Require().NoError(err)

	paidAt := time.Now().UTC()
	order2 := model.RestoreOrder(
		uuid.New(),
		s.testLocID,
		nil,
		model.OrderPaid,
		2000,
		time.Now().UTC(),
		&paidAt,
	)
	err = s.repo.Create(s.ctx, order2)
	s.Require().NoError(err)

	pendingOrders, err := s.repo.ListByStatus(s.ctx, model.OrderPending, 10, 0)
	s.Require().NoError(err)
	s.Require().Len(pendingOrders, 1)
	s.Require().Equal(s.testOrder.ID(), pendingOrders[0].ID())

	paidOrders, err := s.repo.ListByStatus(s.ctx, model.OrderPaid, 10, 0)
	s.Require().NoError(err)
	s.Require().Len(paidOrders, 1)
	s.Require().Equal(order2.ID(), paidOrders[0].ID())
}
