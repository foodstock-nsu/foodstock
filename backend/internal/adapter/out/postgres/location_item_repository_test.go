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

type LocationItemRepoSuite struct {
	suite.Suite
	pgContainer *testhelpers.PostgresContainer
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
	suite.Run(t, new(LocationItemRepoSuite))
}

func (s *LocationItemRepoSuite) SetupSuite() {
	const targetVersion = 3

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
	s.repo = adapterpostgres.NewLocationItemRepository(
		s.dbClient,
		trmpgx.DefaultCtxGetter,
	)
	s.ctx = ctx

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
			nil,
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
	_ = s.pgContainer.MigrateDown()
	s.dbClient.Close()
	_ = s.pgContainer.Close(s.ctx)
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

	_ = s.testLocItem.Update(
		utils.VPtr(int64(2500)), nil, utils.VPtr(50),
	)

	err = s.repo.Update(s.ctx, s.testLocItem)
	s.Require().NoError(err)

	res, err := s.repo.GetByID(s.ctx, s.testLocItem.ID())
	s.Require().NoError(err)
	s.Require().Equal(int64(2500), res.Price())
	s.Require().Equal(50, res.StockAmount())
	s.Require().True(res.IsAvailable())
}

func (s *LocationItemRepoSuite) TestDeleteByItemID() {
	err := s.repo.Create(s.ctx, s.testLocItem)
	s.Require().NoError(err)

	err = s.repo.DeleteByItemID(s.ctx, s.testLocItem.ItemID())
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

	items, err := s.repo.List(s.ctx, s.testLocID)
	s.Require().NoError(err)
	s.Require().Len(items, 2)
}
