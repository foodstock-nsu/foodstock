package postgres

import (
	"backend/internal/adapter/out/postgres/mapper"
	"backend/internal/adapter/out/postgres/sqlc"
	"backend/internal/domain/model"
	pkgerrs "backend/pkg/errs"
	pkgpostgres "backend/pkg/postgres"
	"context"
	"errors"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LocationRepository struct {
	q      *sqlc.Queries
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func NewLocationRepository(
	pgClient *pkgpostgres.Client,
	getter *trmpgx.CtxGetter,
) *LocationRepository {
	return &LocationRepository{
		q:      sqlc.New(),
		pool:   pgClient.Pool,
		getter: getter,
	}
}

func (r *LocationRepository) Create(ctx context.Context, loc *model.Location) error {
	db := r.getter.DefaultTrOrDB(ctx, r.pool)
	params := mapper.MapLocationToSQLCCreate(loc)

	if err := r.q.CreateLocation(ctx, db, params); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return pkgerrs.NewObjectAlreadyExistsErrorWithReason(
					"location", pgErr,
				)
			}
		}
		return err
	}

	return nil
}

func (r *LocationRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Location, error) {
	db := r.getter.DefaultTrOrDB(ctx, r.pool)

	rawLoc, err := r.q.GetLocationByID(
		ctx,
		db,
		pgtype.UUID{
			Bytes: id,
			Valid: true,
		},
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pkgerrs.NewObjectNotFoundError("location", id)
		}
		return nil, err
	}

	return mapper.MapSQLCToLocation(rawLoc), nil
}

func (r *LocationRepository) GetBySlug(ctx context.Context, slug string) (*model.Location, error) {
	db := r.getter.DefaultTrOrDB(ctx, r.pool)

	rawLoc, err := r.q.GetLocationBySlug(ctx, db, slug)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pkgerrs.NewObjectNotFoundError("location", slug)
		}
		return nil, err
	}

	return mapper.MapSQLCToLocation(rawLoc), nil
}

func (r *LocationRepository) Update(ctx context.Context, loc *model.Location) error {
	db := r.getter.DefaultTrOrDB(ctx, r.pool)
	params := mapper.MapLocationToSQLCUpdate(loc)
	return r.q.UpdateLocation(ctx, db, params)
}

func (r *LocationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	db := r.getter.DefaultTrOrDB(ctx, r.pool)

	rowsAffected, err := r.q.DeleteLocation(
		ctx,
		db,
		pgtype.UUID{
			Bytes: id,
			Valid: true,
		},
	)
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return pkgerrs.NewObjectNotFoundError("location", id)
	}

	return nil
}

func (r *LocationRepository) List(ctx context.Context) ([]*model.Location, error) {
	db := r.getter.DefaultTrOrDB(ctx, r.pool)

	rawLocs, err := r.q.ListLocations(ctx, db)
	if err != nil {
		return nil, err
	}

	return mapper.MapSQLCToLocations(rawLocs), nil
}
