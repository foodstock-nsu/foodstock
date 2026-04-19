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

type LocationItemRepoSuite struct {
	suite.Suite
	dbClient    *pkgpostgres.Client
	repo        *adapterpostgres.LocationItemRepository
	ctx         context.Context
	migrate     *migrate.Migrate
	testLocItem *model.LocationItem
	testLocID   uuid.UUID
	testItemID1 uuid.UUID
	testItemID2 uuid.UUID
}

func TestLocationItemRepoSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}
	suite.Run(t, new(LocationItemRepoSuite))
}

func (s *LocationItemRepoSuite) setupDatabase() {
	const targetVersion = 3

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

func (s *LocationItemRepoSuite) SetupSuite() {
	s.ctx = context.Background()
	s.setupDatabase()
	s.repo = adapterpostgres.NewLocationItemRepository(
		s.dbClient,
		trmpgx.DefaultCtxGetter,
	)

	s.testLocID = uuid.New()
	s.testItemID1 = uuid.New()
	s.testItemID2 = uuid.New()

	locationRepo := adapterpostgres.NewLocationRepository(
		s.dbClient,
		trmpgx.DefaultCtxGetter,
	)
	_ = locationRepo.Create(
		s.ctx,
		model.RestoreLocation(
			s.testLocID,
			"nsu_1",
			"Novosibirsk State University | Store №1",
			"Novosibirsk, some st., 6300019",
			true,
			time.Now().UTC(),
		),
	)

	itemRepo := adapterpostgres.NewItemRepository(
		s.dbClient,
		trmpgx.DefaultCtxGetter,
	)
	_ = itemRepo.Create(
		s.ctx,
		model.RestoreItem(
			s.testItemID1,
			"Chicken Sandwich",
			utils.VPtr("Chicken sandwich with fresh vegetables"),
			model.ItemLunch,
			"https://example.com/photo.jpg",
			model.RestoreNutrition(
				utils.VPtr(100),
				utils.VPtr(float64(20)),
				utils.VPtr(float64(10)),
				utils.VPtr(float64(30))),
			time.Now().UTC(),
		),
	)
	_ = itemRepo.Create(
		s.ctx,
		model.RestoreItem(
			s.testItemID2,
			"Beef Sandwich",
			utils.VPtr("Beef sandwich with fresh vegetables"),
			model.ItemLunch,
			"https://example.com/photo.jpg",
			model.RestoreNutrition(
				utils.VPtr(100),
				utils.VPtr(float64(20)),
				utils.VPtr(float64(10)),
				utils.VPtr(float64(30))),
			time.Now().UTC(),
		),
	)

	s.testLocItem = model.RestoreLocationItem(
		uuid.New(),
		s.testItemID1,
		s.testLocID,
		1500,
		true,
		100,
	)
}

func (s *LocationItemRepoSuite) TearDownSuite() {
	if s.migrate != nil {
		_ = s.migrate.Down()
	}
	s.dbClient.Close()
}

func (s *LocationItemRepoSuite) SetupTest() {
	_, err := s.dbClient.Pool.Exec(s.ctx, "TRUNCATE TABLE location_items CASCADE")
	s.Require().NoError(err)
}

func (s *LocationItemRepoSuite) TestCreateGetByID() {
	err := s.repo.Create(s.ctx, s.testLocItem)
	s.Require().NoError(err)

	item, err := s.repo.GetByID(s.ctx, s.testLocItem.ID())
	s.Require().NoError(err)
	s.Require().NotNil(item)
	s.Require().Equal(s.testLocItem.ID(), item.ID())
	s.Require().Equal(s.testLocItem.Price(), item.Price())
	s.Require().Equal(s.testLocItem.StockAmount(), item.StockAmount())
	s.Require().Equal(s.testLocItem.IsAvailable(), item.IsAvailable())
}

func (s *LocationItemRepoSuite) TestCreate_AlreadyExists() {
	err := s.repo.Create(s.ctx, s.testLocItem)
	s.Require().NoError(err)

	err = s.repo.Create(s.ctx, s.testLocItem)
	s.Require().Error(err)
	s.Require().ErrorIs(err, pkgerrs.ErrObjectAlreadyExists)
}

func (s *LocationItemRepoSuite) TestGetByID_NotFound() {
	item, err := s.repo.GetByID(s.ctx, uuid.New())
	s.Require().Error(err)
	s.Require().ErrorIs(err, pkgerrs.ErrObjectNotFound)
	s.Require().Nil(item)
}

func (s *LocationItemRepoSuite) TestGetByLocationAndItem() {
	err := s.repo.Create(s.ctx, s.testLocItem)
	s.Require().NoError(err)

	item, err := s.repo.GetByLocationAndItem(s.ctx, s.testLocID, s.testItemID1)
	s.Require().NoError(err)
	s.Require().NotNil(item)
	s.Require().Equal(s.testLocItem.ID(), item.ID())
}

func (s *LocationItemRepoSuite) TestUpdate() {
	err := s.repo.Create(s.ctx, s.testLocItem)
	s.Require().NoError(err)

	_ = s.testLocItem.Update(utils.VPtr(int64(2500)), utils.VPtr(50))

	err = s.repo.Update(s.ctx, s.testLocItem)
	s.Require().NoError(err)

	res, err := s.repo.GetByID(s.ctx, s.testLocItem.ID())
	s.Require().NoError(err)
	s.Require().Equal(int64(2500), res.Price())
	s.Require().Equal(50, res.StockAmount())
	s.Require().True(res.IsAvailable())
}

func (s *LocationItemRepoSuite) TestDeleteByID() {
	err := s.repo.Create(s.ctx, s.testLocItem)
	s.Require().NoError(err)

	err = s.repo.DeleteByID(s.ctx, s.testLocItem.ID())
	s.Require().NoError(err)

	res, err := s.repo.GetByID(s.ctx, s.testLocItem.ID())
	s.Require().Error(err)
	s.Require().ErrorIs(err, pkgerrs.ErrObjectNotFound)
	s.Require().Nil(res)
}

func (s *LocationItemRepoSuite) TestDeleteByLocationID() {
	err := s.repo.Create(s.ctx, s.testLocItem)
	s.Require().NoError(err)

	err = s.repo.DeleteByLocationID(s.ctx, s.testLocItem.LocationID())
	s.Require().NoError(err)

	res, err := s.repo.GetByID(s.ctx, s.testLocItem.ID())
	s.Require().Error(err)
	s.Require().ErrorIs(err, pkgerrs.ErrObjectNotFound)
	s.Require().Nil(res)
}

func (s *LocationItemRepoSuite) TestList() {
	err := s.repo.Create(s.ctx, s.testLocItem)
	s.Require().NoError(err)

	// Second item for the same location
	item2 := model.RestoreLocationItem(
		uuid.New(),
		s.testItemID2,
		s.testLocID,
		1000,
		true,
		10,
	)
	err = s.repo.Create(s.ctx, item2)
	s.Require().NoError(err)

	items, err := s.repo.List(s.ctx, s.testLocID, 10, 0)
	s.Require().NoError(err)
	s.Require().Len(items, 2)

	// Test pagination
	itemsPaged, err := s.repo.List(s.ctx, s.testLocID, 1, 1)
	s.Require().NoError(err)
	s.Require().Len(itemsPaged, 1)
}
