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

type ItemRepoSuite struct {
	suite.Suite
	pgContainer *testhelpers.PostgresContainer
	dbClient    *pkgpostgres.Client
	repo        *adapterpostgres.ItemRepository
	ctx         context.Context
	migrate     *migrate.Migrate
	testItem    *model.Item
}

func TestItemRepoSuite(t *testing.T) {
	suite.Run(t, new(ItemRepoSuite))
}

func (s *ItemRepoSuite) SetupSuite() {
	const targetVersion = 2

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
	s.repo = adapterpostgres.NewItemRepository(
		s.dbClient,
		trmpgx.DefaultCtxGetter,
	)
	s.ctx = ctx
	s.testItem = model.RestoreItem(
		uuid.New(),
		"Test Item",
		utils.VPtr("Delicious test item description long enough"),
		model.ItemLunch,
		"https://example.com/photo.jpg",
		model.RestoreNutrition(
			utils.VPtr(100),
			utils.VPtr(float64(10)),
			utils.VPtr(float64(5)),
			utils.VPtr(float64(15))),
		time.Now().UTC().Truncate(time.Microsecond),
	)
}

func (s *ItemRepoSuite) TearDownSuite() {
	_ = s.pgContainer.MigrateDown()
	s.dbClient.Close()
	_ = s.pgContainer.Close(s.ctx)
}

func (s *ItemRepoSuite) SetupTest() {
	_, err := s.dbClient.Pool.Exec(s.ctx, "TRUNCATE TABLE items CASCADE")
	s.Require().NoError(err)
}

func (s *ItemRepoSuite) TestCreateGet() {
	err := s.repo.Create(s.ctx, s.testItem)
	s.Require().NoError(err)

	item, err := s.repo.Get(s.ctx, s.testItem.ID())
	s.Require().NoError(err)
	s.Require().NotNil(item)
	s.Require().Equal(s.testItem.ID(), item.ID())
	s.Require().Equal(s.testItem.Name(), item.Name())
	s.Require().Equal(s.testItem.Category(), item.Category())
	s.Require().Equal(s.testItem.Nutrition().Calories(), item.Nutrition().Calories())
}

func (s *ItemRepoSuite) TestCreate_AlreadyExists() {
	err := s.repo.Create(s.ctx, s.testItem)
	s.Require().NoError(err)

	err = s.repo.Create(s.ctx, s.testItem)
	s.Require().Error(err)
	s.Require().ErrorIs(err, pkgerrs.ErrObjectAlreadyExists)
}

func (s *ItemRepoSuite) TestGet_NotFound() {
	item, err := s.repo.Get(s.ctx, uuid.New())
	s.Require().Error(err)
	s.Require().ErrorIs(err, pkgerrs.ErrObjectNotFound)
	s.Require().Nil(item)
}

func (s *ItemRepoSuite) TestUpdate() {
	err := s.repo.Create(s.ctx, s.testItem)
	s.Require().NoError(err)

	newName := "Updated Item Name"
	newDesc := "Updated description that is also long enough"

	// Using RestoreItem to simulate updated state since model might be rich
	updatedItem := model.RestoreItem(
		s.testItem.ID(),
		newName,
		&newDesc,
		model.ItemBreakfast,
		s.testItem.PhotoURL(),
		model.RestoreNutrition(
			utils.VPtr(100),
			utils.VPtr(float64(20)),
			utils.VPtr(float64(10)),
			utils.VPtr(float64(30))),
		s.testItem.CreatedAt(),
	)

	err = s.repo.Update(s.ctx, updatedItem)
	s.Require().NoError(err)

	res, err := s.repo.Get(s.ctx, s.testItem.ID())
	s.Require().NoError(err)
	s.Require().Equal(newName, res.Name())
	s.Require().Equal(model.ItemCategory("breakfast"), res.Category())
	s.Require().Equal(100, *res.Nutrition().Calories())
}

func (s *ItemRepoSuite) TestDelete() {
	err := s.repo.Create(s.ctx, s.testItem)
	s.Require().NoError(err)

	err = s.repo.Delete(s.ctx, s.testItem.ID())
	s.Require().NoError(err)

	res, err := s.repo.Get(s.ctx, s.testItem.ID())
	s.Require().Error(err)
	s.Require().ErrorIs(err, pkgerrs.ErrObjectNotFound)
	s.Require().Nil(res)
}

func (s *ItemRepoSuite) TestDelete_NotFound() {
	err := s.repo.Delete(s.ctx, uuid.New())
	s.Require().Error(err)
	s.Require().ErrorIs(err, pkgerrs.ErrObjectNotFound)
}

func (s *ItemRepoSuite) TestListAll() {
	err := s.repo.Create(s.ctx, s.testItem)
	s.Require().NoError(err)

	// Second item
	item2 := model.RestoreItem(
		uuid.New(),
		"Item 2",
		utils.VPtr("Another description"),
		model.ItemDrinks,
		"https://example.com/photo.jpg",
		model.RestoreNutrition(
			utils.VPtr(100),
			utils.VPtr(float64(20)),
			utils.VPtr(float64(10)),
			utils.VPtr(float64(30))),
		time.Now().UTC(),
	)
	err = s.repo.Create(s.ctx, item2)
	s.Require().NoError(err)

	items, err := s.repo.ListAll(s.ctx)
	s.Require().NoError(err)
	s.Require().Len(items, 2)
}

func (s *ItemRepoSuite) TestListByIDs() {
	err := s.repo.Create(s.ctx, s.testItem)
	s.Require().NoError(err)

	// Second item
	item2 := model.RestoreItem(
		uuid.New(),
		"Item 2",
		utils.VPtr("Another description"),
		model.ItemDrinks,
		"https://example.com/photo.jpg",
		model.RestoreNutrition(
			utils.VPtr(100),
			utils.VPtr(float64(20)),
			utils.VPtr(float64(10)),
			utils.VPtr(float64(30))),
		time.Now().UTC(),
	)
	err = s.repo.Create(s.ctx, item2)
	s.Require().NoError(err)

	items, err := s.repo.ListByIDs(
		s.ctx, []uuid.UUID{s.testItem.ID(), item2.ID()},
	)
	s.Require().NoError(err)
	s.Require().Len(items, 2)

	// Test if only one specified
	items, err = s.repo.ListByIDs(s.ctx, []uuid.UUID{item2.ID()})
	s.Require().NoError(err)
	s.Require().Len(items, 1)

	// Test if no single id is specified
	items, err = s.repo.ListByIDs(s.ctx, nil)
	s.Require().NoError(err)
	s.Require().Nil(items)

}
