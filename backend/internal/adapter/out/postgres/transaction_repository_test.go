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

type TransactionRepoSuite struct {
	suite.Suite
	dbClient    *pkgpostgres.Client
	repo        *adapterpostgres.TransactionRepository
	ctx         context.Context
	migrate     *migrate.Migrate
	testOrderID uuid.UUID
	testTx      *model.Transaction
}

func TestTransactionRepoSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}
	suite.Run(t, new(TransactionRepoSuite))
}

func (s *TransactionRepoSuite) setupDatabase() {
	const targetVersion = 6

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

func (s *TransactionRepoSuite) SetupSuite() {
	s.ctx = context.Background()
	s.setupDatabase()
	s.repo = adapterpostgres.NewTransactionRepository(
		s.dbClient,
		trmpgx.DefaultCtxGetter,
	)

	locID := uuid.New()
	locationRepo := adapterpostgres.NewLocationRepository(s.dbClient, trmpgx.DefaultCtxGetter)
	_ = locationRepo.Create(s.ctx, model.RestoreLocation(locID, "tx_test_loc", "Shop", "Addr", true, time.Now()))

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
	if s.migrate != nil {
		_ = s.migrate.Down()
	}
	s.dbClient.Close()
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
