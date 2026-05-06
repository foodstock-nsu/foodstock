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

type AdminRepoSuite struct {
	suite.Suite
	pgContainer *testhelpers.PostgresContainer
	dbClient    *pkgpostgres.Client
	repo        *adapterpostgres.AdminRepository
	ctx         context.Context
	migrate     *migrate.Migrate
	testAdmin   *model.Admin
}

func TestAdminRepoSuite(t *testing.T) {
	suite.Run(t, new(AdminRepoSuite))
}

func (s *AdminRepoSuite) SetupSuite() {
	const targetVersion = 7

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
	s.repo = adapterpostgres.NewAdminRepository(
		s.dbClient,
		trmpgx.DefaultCtxGetter,
	)
	s.ctx = ctx
	s.testAdmin = model.RestoreAdmin(
		uuid.New(),
		"admin",
		"hash",
		time.Now().UTC(),
	)
}

func (s *AdminRepoSuite) TearDownSuite() {
	_ = s.pgContainer.MigrateDown()
	s.dbClient.Close()
	_ = s.pgContainer.Close(s.ctx)
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
