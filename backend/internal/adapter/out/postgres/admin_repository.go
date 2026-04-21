package postgres

import (
	"backend/internal/adapter/out/postgres/mapper"
	"backend/internal/adapter/out/postgres/sqlc"
	"backend/internal/domain/model"
	pkgerrs "backend/pkg/errs"
	pkgpostgres "backend/pkg/postgres"
	"context"
	"errors"
	"time"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AdminRepository struct {
	q      *sqlc.Queries
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func NewAdminRepository(
	pgClient *pkgpostgres.Client,
	getter *trmpgx.CtxGetter,
) *AdminRepository {
	return &AdminRepository{
		q:      sqlc.New(),
		pool:   pgClient.Pool,
		getter: getter,
	}
}

func (r *AdminRepository) Create(ctx context.Context, admin *model.Admin) error {
	db := r.getter.DefaultTrOrDB(ctx, r.pool)
	params := mapper.MapAdminToSQLCCreate(admin)

	if err := r.q.CreateAdmin(ctx, db, params); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return pkgerrs.NewObjectAlreadyExistsErrorWithReason(
					"admin", pgErr,
				)
			}
		}
		return err
	}

	return nil
}

func (r *AdminRepository) Upsert(ctx context.Context, admin *model.Admin) error {
	db := r.getter.DefaultTrOrDB(ctx, r.pool)
	params := mapper.MapAdminToSQLCUpsert(admin)
	return r.q.UpsertAdmin(ctx, db, params)
}

func (r *AdminRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Admin, error) {
	db := r.getter.DefaultTrOrDB(ctx, r.pool)

	rawAdmin, err := r.q.GetAdminByID(
		ctx,
		db,
		pgtype.UUID{
			Bytes: id,
			Valid: true,
		},
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pkgerrs.NewObjectNotFoundError("admin", id)
		}
		return nil, err
	}

	return mapper.MapSQLCToAdmin(rawAdmin), nil
}

func (r *AdminRepository) GetByLogin(ctx context.Context, login string) (*model.Admin, error) {
	db := r.getter.DefaultTrOrDB(ctx, r.pool)

	rawAdmin, err := r.q.GetAdminByLogin(ctx, db, login)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pkgerrs.NewObjectNotFoundError("admin", login)
		}
		return nil, err
	}

	return mapper.MapSQLCToAdmin(rawAdmin), nil
}

func (r *AdminRepository) EnsureAdmin(ctx context.Context, login string, passwordHash string) error {
	db := r.getter.DefaultTrOrDB(ctx, r.pool)

	adminNamespace := uuid.MustParse("00000000-0000-0000-0000-000000000001")

	id := uuid.NewSHA1(adminNamespace, []byte(login))

	params := sqlc.UpsertAdminParams{
		ID:           pgtype.UUID{Bytes: id, Valid: true},
		Login:        login,
		PasswordHash: passwordHash,
		CreatedAt:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	return r.q.UpsertAdmin(ctx, db, params)
}
