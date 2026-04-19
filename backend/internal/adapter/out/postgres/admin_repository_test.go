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

type AdminRepoSuite struct {
	suite.Suite
	dbClient  *pkgpostgres.Client
	repo      *adapterpostgres.AdminRepository
	ctx       context.Context
	migrate   *migrate.Migrate
	testAdmin *model.Admin
}

func TestAdminRepoSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}
	suite.Run(t, new(AdminRepoSuite))
}

func (s *AdminRepoSuite) setupDatabase() {
	const targetVersion = 7

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

func (s *AdminRepoSuite) SetupSuite() {
	s.ctx = context.Background()
	s.setupDatabase()
	s.repo = adapterpostgres.NewAdminRepository(
		s.dbClient,
		trmpgx.DefaultCtxGetter,
	)
	s.testAdmin = model.RestoreAdmin(
		uuid.New(),
		"admin",
		"hash",
		time.Now().UTC(),
	)
}

func (s *AdminRepoSuite) TearDownSuite() {
	if s.migrate != nil {
		_ = s.migrate.Down()
	}
	s.dbClient.Close()
}

func (s *AdminRepoSuite) SetupTest() {
	_, err := s.dbClient.Pool.Exec(s.ctx, "TRUNCATE TABLE admins CASCADE")
	s.Require().NoError(err)
}

func (s *AdminRepoSuite) TestCreateGet() {
	err := s.repo.Create(s.ctx, s.testAdmin)
	s.Require().NoError(err)

	// Firstly get the admin by id
	admin, err := s.repo.GetByID(s.ctx, s.testAdmin.ID())
	s.Require().NoError(err)
	s.Require().NotNil(admin)
	s.Require().Equal(s.testAdmin.Login(), admin.Login())
	s.Require().Equal(s.testAdmin.PasswordHash(), admin.PasswordHash())

	// Secondly get the admin by login
	admin, err = s.repo.GetByLogin(s.ctx, s.testAdmin.Login())
	s.Require().NoError(err)
	s.Require().NotNil(admin)
	s.Require().Equal(s.testAdmin.ID(), admin.ID())
	s.Require().Equal(s.testAdmin.PasswordHash(), admin.PasswordHash())
}

func (s *AdminRepoSuite) TestCreate_AlreadyExists() {
	err := s.repo.Create(s.ctx, s.testAdmin)
	s.Require().NoError(err)

	err = s.repo.Create(s.ctx, s.testAdmin)
	s.Require().Error(err)
	s.Require().ErrorIs(err, pkgerrs.ErrObjectAlreadyExists)
}

func (s *AdminRepoSuite) TestGetByID_NotFound() {
	admin, err := s.repo.GetByID(s.ctx, uuid.New())
	s.Require().Error(err)
	s.Require().ErrorIs(err, pkgerrs.ErrObjectNotFound)
	s.Require().Nil(admin)
}

func (s *AdminRepoSuite) TestGetByLogin_NotFound() {
	admin, err := s.repo.GetByLogin(s.ctx, "unknown")
	s.Require().Error(err)
	s.Require().ErrorIs(err, pkgerrs.ErrObjectNotFound)
	s.Require().Nil(admin)
}

func (s *AdminRepoSuite) TestUpsert() {
	// Firstly create a new admin using upsert method
	err := s.repo.Upsert(s.ctx, s.testAdmin)
	s.Require().NoError(err)

	// Check if the admin was created successfully
	admin, err := s.repo.GetByID(s.ctx, s.testAdmin.ID())
	s.Require().NoError(err)
	s.Require().NotNil(admin)
	s.Require().Equal(s.testAdmin.Login(), admin.Login())
	s.Require().Equal(s.testAdmin.PasswordHash(), admin.PasswordHash())

	// Secondly, try to create admin with the same parameters except password
	var newPasswordHash = "new-hash"
	admin = model.RestoreAdmin(
		s.testAdmin.ID(),
		s.testAdmin.Login(),
		newPasswordHash,
		s.testAdmin.CreatedAt(),
	)

	err = s.repo.Upsert(s.ctx, admin)
	s.Require().NoError(err)

	// Check if password hash was updated successfully
	admin, err = s.repo.GetByID(s.ctx, s.testAdmin.ID())
	s.Require().NoError(err)
	s.Require().NotNil(admin)
	s.Require().Equal(s.testAdmin.Login(), admin.Login())
	s.Require().Equal(newPasswordHash, admin.PasswordHash())
}
