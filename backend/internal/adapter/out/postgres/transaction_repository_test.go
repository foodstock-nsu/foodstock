//go:build integration

package postgres_test

import (
	adapterpostgres "backend/internal/adapter/out/postgres"
	"backend/internal/domain/model"
	"backend/internal/testhelpers"
	pkgerrs "backend/pkg/errs"
	pkgpostgres "backend/pkg/postgres"
	"context"
	"testing"
	"time"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type TransactionRepoSuite struct {
	suite.Suite
	pgContainer *testhelpers.PostgresContainer
	dbClient    *pkgpostgres.Client
	repo        *adapterpostgres.TransactionRepository
	ctx         context.Context
	migrate     *migrate.Migrate
	testOrderID uuid.UUID
	testTx      *model.Transaction
}

func TestTransactionRepoSuite(t *testing.T) {
	suite.Run(t, new(TransactionRepoSuite))
}

func (s *TransactionRepoSuite) SetupSuite() {
	const targetVersion = 6

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
	s.repo = adapterpostgres.NewTransactionRepository(
		s.dbClient,
		trmpgx.DefaultCtxGetter,
	)
	s.ctx = ctx

	locID := uuid.New()
	locationRepo := adapterpostgres.NewLocationRepository(s.dbClient, trmpgx.DefaultCtxGetter)
	_ = locationRepo.Create(s.ctx, model.RestoreLocation(
		locID, "tx_test_loc", "Shop",
		"Addr", true, time.Now(), nil,
	))

	s.testOrderID = uuid.New()
	orderRepo := adapterpostgres.NewOrderRepository(s.dbClient, trmpgx.DefaultCtxGetter)
	_ = orderRepo.Create(s.ctx, model.RestoreOrder(s.testOrderID, locID, nil, model.OrderPending, 1000, time.Now(), nil))

	s.testTx = model.RestoreTransaction(
		uuid.New(),
		s.testOrderID,
		"sbp-unique-id-123",
		1000,
		model.TransactionPending,
		nil,
		time.Now().UTC().Truncate(time.Microsecond),
	)
}

func (s *TransactionRepoSuite) TearDownSuite() {
	_ = s.pgContainer.MigrateDown()
	s.dbClient.Close()
	_ = s.pgContainer.Close(s.ctx)
}

func (s *TransactionRepoSuite) SetupTest() {
	_, err := s.dbClient.Pool.Exec(s.ctx, "TRUNCATE TABLE transactions CASCADE")
	s.Require().NoError(err)
}

func (s *TransactionRepoSuite) TestCreateGetByID() {
	err := s.repo.Create(s.ctx, s.testTx)
	s.Require().NoError(err)

	tx, err := s.repo.GetByID(s.ctx, s.testTx.ID())
	s.Require().NoError(err)
	s.Require().NotNil(tx)
	s.Require().Equal(s.testTx.ID(), tx.ID())
	s.Require().Equal(s.testTx.SBPTransactionID(), tx.SBPTransactionID())
	s.Require().Equal(s.testTx.Status(), tx.Status())
}

func (s *TransactionRepoSuite) TestCreate_AlreadyExists() {
	err := s.repo.Create(s.ctx, s.testTx)
	s.Require().NoError(err)

	err = s.repo.Create(s.ctx, s.testTx)
	s.Require().Error(err)
	s.Require().ErrorIs(err, pkgerrs.ErrObjectAlreadyExists)
}

func (s *TransactionRepoSuite) TestGetByID_NotFound() {
	tx, err := s.repo.GetByID(s.ctx, uuid.New())
	s.Require().Error(err)
	s.Require().ErrorIs(err, pkgerrs.ErrObjectNotFound)
	s.Require().Nil(tx)
}

func (s *TransactionRepoSuite) TestGetLatestByOrderID() {
	err := s.repo.Create(s.ctx, s.testTx)
	s.Require().NoError(err)

	tx, err := s.repo.GetLatestByOrderID(s.ctx, s.testTx.OrderID())
	s.Require().NoError(err)
	s.Require().NotNil(tx)
	s.Require().Equal(s.testTx.ID(), tx.ID())
}

func (s *TransactionRepoSuite) TestGetLatestByOrderID_NotFound() {
	tx, err := s.repo.GetLatestByOrderID(s.ctx, uuid.New())
	s.Require().Error(err)
	s.Require().ErrorIs(err, pkgerrs.ErrObjectNotFound)
	s.Require().Nil(tx)
}

func (s *TransactionRepoSuite) TestGetBySbpID() {
	err := s.repo.Create(s.ctx, s.testTx)
	s.Require().NoError(err)

	tx, err := s.repo.GetBySbpID(s.ctx, s.testTx.SBPTransactionID())
	s.Require().NoError(err)
	s.Require().NotNil(tx)
	s.Require().Equal(s.testTx.ID(), tx.ID())
}

func (s *TransactionRepoSuite) TestGetBySbpID_NotFound() {
	tx, err := s.repo.GetBySbpID(s.ctx, "non-existent")
	s.Require().Error(err)
	s.Require().ErrorIs(err, pkgerrs.ErrObjectNotFound)
	s.Require().Nil(tx)
}

func (s *TransactionRepoSuite) TestUpdate() {
	err := s.repo.Create(s.ctx, s.testTx)
	s.Require().NoError(err)

	paidAt := time.Now().UTC().Truncate(time.Microsecond)
	updatedTx := model.RestoreTransaction(
		s.testTx.ID(),
		s.testTx.OrderID(),
		s.testTx.SBPTransactionID(),
		s.testTx.Amount(),
		model.TransactionSuccess,
		&paidAt,
		s.testTx.CreatedAt(),
	)

	err = s.repo.Update(s.ctx, updatedTx)
	s.Require().NoError(err)

	res, err := s.repo.GetByID(s.ctx, s.testTx.ID())
	s.Require().NoError(err)
	s.Require().Equal(model.TransactionSuccess, res.Status())
	s.Require().NotNil(res.PaidAt())
	s.Require().True(paidAt.Equal(*res.PaidAt()))
}

func (s *TransactionRepoSuite) TestList() {
	err := s.repo.Create(s.ctx, s.testTx)
	s.Require().NoError(err)

	tx2 := model.RestoreTransaction(
		uuid.New(),
		s.testOrderID,
		"sbp-id-2",
		2000,
		model.TransactionFailed,
		nil,
		time.Now().UTC(),
	)
	err = s.repo.Create(s.ctx, tx2)
	s.Require().NoError(err)

	txs, err := s.repo.List(s.ctx, s.testOrderID)
	s.Require().NoError(err)
	s.Require().Len(txs, 2)
}
