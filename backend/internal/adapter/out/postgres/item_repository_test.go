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

type ItemRepoSuite struct {
	suite.Suite
	dbClient *pkgpostgres.Client
	repo     *adapterpostgres.ItemRepository
	ctx      context.Context
	migrate  *migrate.Migrate
	testItem *model.Item
}

func TestItemRepoSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}
	suite.Run(t, new(ItemRepoSuite))
}

func (s *ItemRepoSuite) setupDatabase() {
	const targetVersion = 2

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

func (s *ItemRepoSuite) SetupSuite() {
	s.ctx = context.Background()
	s.setupDatabase()
	s.repo = adapterpostgres.NewItemRepository(
		s.dbClient,
		trmpgx.DefaultCtxGetter,
	)

	desc := "Delicious test item description long enough"
	s.testItem = model.RestoreItem(
		uuid.New(),
		"Test Item",
		&desc,
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
	if s.migrate != nil {
		_ = s.migrate.Down()
	}
	s.dbClient.Close()
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

func (s *ItemRepoSuite) TestList() {
	err := s.repo.Create(s.ctx, s.testItem)
	s.Require().NoError(err)

	// Second item
	desc2 := "Another description"
	item2 := model.RestoreItem(
		uuid.New(),
		"Item 2",
		&desc2,
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

	items, err := s.repo.List(s.ctx)
	s.Require().NoError(err)
	s.Require().Len(items, 2)
}
