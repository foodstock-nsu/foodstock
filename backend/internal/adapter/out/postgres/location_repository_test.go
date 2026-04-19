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

type LocationRepoSuite struct {
	suite.Suite
	dbClient     *pkgpostgres.Client
	repo         *adapterpostgres.LocationRepository
	ctx          context.Context
	migrate      *migrate.Migrate
	testLocation *model.Location
}

func TestLocationRepoSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}
	suite.Run(t, new(LocationRepoSuite))
}

func (s *LocationRepoSuite) setupDatabase() {
	// Version of the lowest migration to apply
	const targetVersion = 1

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
	s.Require().NoError(err, "failed to create iofs driver")

	m, err := migrate.NewWithSourceInstance(
		"iofs",
		sourceDriver,
		dsn,
	)
	s.Require().NoError(err, "failed to create migration instance")

	s.migrate = m

	err = m.Migrate(targetVersion)

	// If migration is correct - setup has done
	if err == nil || errors.Is(err, migrate.ErrNoChange) {
		return
	}

	// Except dirty db as a normal scenario
	var dirtyErr migrate.ErrDirty
	if !errors.As(err, &dirtyErr) {
		s.FailNowf("failed to migrate up", "unexpected error: %v", err)
	}

	// ================ Restore dirty database ================
	_ = m.Force(dirtyErr.Version)

	err = m.Down()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		s.Require().NoError(err, "failed to migrate down during recovery")
	}

	err = m.Migrate(targetVersion)
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		s.Require().NoError(err, "failed to migrate up after recovery")
	}
}

func (s *LocationRepoSuite) SetupSuite() {
	s.ctx = context.Background()
	s.setupDatabase()
	s.repo = adapterpostgres.NewLocationRepository(
		s.dbClient,
		trmpgx.DefaultCtxGetter,
	)

	s.testLocation, _ = model.NewLocation(
		"test-slug",
		"Test Location Name",
		"Test Address 123456789",
	)
}

func (s *LocationRepoSuite) TearDownSuite() {
	if s.migrate != nil {
		if err := s.migrate.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			s.Require().NoError(err, "failed to migrate down")
		}
	}
	s.dbClient.Close()
}

func (s *LocationRepoSuite) SetupTest() {
	_, err := s.dbClient.Pool.Exec(s.ctx, "TRUNCATE TABLE locations CASCADE")
	s.Require().NoError(err)
}

func (s *LocationRepoSuite) TestCreateGetByID() {
	// Create a test location at first
	err := s.repo.Create(s.ctx, s.testLocation)
	s.Require().NoError(err)

	// Then get it by id
	loc, err := s.repo.GetByID(s.ctx, s.testLocation.ID())
	s.Require().NoError(err)
	s.Require().NotNil(loc)
	s.Require().Equal(s.testLocation.Slug(), loc.Slug())
	s.Require().Equal(s.testLocation.Name(), loc.Name())
}

func (s *LocationRepoSuite) TestCreate_AlreadyExists() {
	_ = s.repo.Create(s.ctx, s.testLocation)

	err := s.repo.Create(s.ctx, s.testLocation)

	s.Require().Error(err)
	s.Require().ErrorIs(err, pkgerrs.ErrObjectAlreadyExists)

	s.Assert().Contains(err.Error(), "location")
}

func (s *LocationRepoSuite) TestGetByID_NotFound() {
	// Try to get a non-existing location by id
	var unexistingID = uuid.New()
	loc, err := s.repo.GetByID(s.ctx, unexistingID)

	s.Require().Error(err)
	s.Require().ErrorIs(err, pkgerrs.ErrObjectNotFound)
	s.Require().Nil(loc)
}

func (s *LocationRepoSuite) TestGetBySlug() {
	// Create the location in advance
	_ = s.repo.Create(s.ctx, s.testLocation)

	// Get it by slug
	loc, err := s.repo.GetBySlug(s.ctx, s.testLocation.Slug())
	s.Require().NoError(err)
	s.Require().NotNil(loc)
	s.Require().Equal(s.testLocation.ID(), loc.ID())
}

func (s *LocationRepoSuite) TestGetBySlug_NotFound() {
	// Try to get a non-existing location by slug
	loc, err := s.repo.GetBySlug(s.ctx, "non-existing-slug")

	s.Require().Error(err)
	s.Require().ErrorIs(err, pkgerrs.ErrObjectNotFound)
	s.Require().Nil(loc)
}

func (s *LocationRepoSuite) TestUpdate() {
	// Create the location in advance
	_ = s.repo.Create(s.ctx, s.testLocation)

	// Modify the model
	_ = s.testLocation.Update(
		utils.VPtr("updated-slug"),
		utils.VPtr("Updated Name"),
		utils.VPtr("Updated Address123456789"),
	)

	// Update it in repository
	err := s.repo.Update(s.ctx, s.testLocation)
	s.Require().NoError(err)

	// Check the result
	loc, _ := s.repo.GetByID(s.ctx, s.testLocation.ID())
	s.Require().Equal("updated-slug", loc.Slug())
	s.Require().Equal("Updated Name", loc.Name())
}

func (s *LocationRepoSuite) TestDelete() {
	// Create the location in advance
	_ = s.repo.Create(s.ctx, s.testLocation)

	// DeleteByID it
	err := s.repo.Delete(s.ctx, s.testLocation.ID())
	s.Require().NoError(err)

	// Try to get it again to ensure it's deleted
	loc, err := s.repo.GetByID(s.ctx, s.testLocation.ID())
	s.Require().Error(err)
	s.Require().Nil(loc)
}

func (s *LocationRepoSuite) TestDelete_NotFound() {
	// Delete an unexisting location
	err := s.repo.Delete(s.ctx, uuid.New())
	s.Require().Error(err)
}

func (s *LocationRepoSuite) TestList() {
	// Create in advance
	_ = s.repo.Create(s.ctx, s.testLocation)

	// Expect 1 item in result
	locs, err := s.repo.List(s.ctx)
	s.Require().NoError(err)
	s.Require().Len(locs, 1)
	s.Require().Equal(s.testLocation.ID(), locs[0].ID())
}
