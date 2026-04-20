package postgres

import (
	"backend/internal/adapter/out/postgres/mapper"
	"backend/internal/adapter/out/postgres/sqlc"
	"backend/internal/domain/model"
	pkgerrs "backend/pkg/errs"
	pkgpostgres "backend/pkg/postgres"
	"context"
	"errors"
	"fmt"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LocationItemRepository struct {
	q      *sqlc.Queries
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func NewLocationItemRepository(
	pgClient *pkgpostgres.Client,
	getter *trmpgx.CtxGetter,
) *LocationItemRepository {
	return &LocationItemRepository{
		q:      sqlc.New(),
		pool:   pgClient.Pool,
		getter: getter,
	}
}

func (r *LocationItemRepository) Create(ctx context.Context, locItem *model.LocationItem) error {
	db := r.getter.DefaultTrOrDB(ctx, r.pool)
	params := mapper.MapLocationItemToSQLCCreate(locItem)

	if err := r.q.CreateLocationItem(ctx, db, params); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return pkgerrs.NewObjectAlreadyExistsErrorWithReason(
					"location_item", pgErr,
				)
			}
		}
		return err
	}

	return nil
}

func (r *LocationItemRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.LocationItem, error) {
	db := r.getter.DefaultTrOrDB(ctx, r.pool)

	rawLocItem, err := r.q.GetLocationItemByID(
		ctx,
		db,
		pgtype.UUID{
			Bytes: id,
			Valid: true,
		},
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pkgerrs.NewObjectNotFoundError("location_item", id)
		}
		return nil, err
	}

	return mapper.MapSQLCToLocationItem(rawLocItem), nil
}

func (r *LocationItemRepository) GetByLocationAndItem(ctx context.Context, locationID, itemID uuid.UUID) (*model.LocationItem, error) {
	db := r.getter.DefaultTrOrDB(ctx, r.pool)

	rawLocItem, err := r.q.GetLocationItemByLocationAndItem(
		ctx,
		db,
		sqlc.GetLocationItemByLocationAndItemParams{
			ItemID: pgtype.UUID{
				Bytes: itemID,
				Valid: true,
			},
			LocationID: pgtype.UUID{
				Bytes: locationID,
				Valid: true,
			},
		},
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pkgerrs.NewObjectNotFoundError(
				"location_id | item_id",
				fmt.Sprintf("%s | %s", locationID, itemID),
			)
		}
		return nil, err
	}

	return mapper.MapSQLCToLocationItem(rawLocItem), nil
}

func (r *LocationItemRepository) Update(ctx context.Context, locItem *model.LocationItem) error {
	db := r.getter.DefaultTrOrDB(ctx, r.pool)
	params := mapper.MapLocationItemToSQLCUpdate(locItem)
	return r.q.UpdateLocationItem(ctx, db, params)
}

func (r *LocationItemRepository) DeleteByItemID(ctx context.Context, itemID uuid.UUID) error {
	db := r.getter.DefaultTrOrDB(ctx, r.pool)
	return r.q.DeleteLocationItemByItemID(
		ctx,
		db,
		pgtype.UUID{
			Bytes: itemID,
			Valid: true,
		},
	)
}

func (r *LocationItemRepository) DeleteByLocationID(ctx context.Context, locationID uuid.UUID) error {
	db := r.getter.DefaultTrOrDB(ctx, r.pool)
	return r.q.DeleteLocationItemsByLocationID(
		ctx,
		db,
		pgtype.UUID{
			Bytes: locationID,
			Valid: true,
		},
	)
}

func (r *LocationItemRepository) List(ctx context.Context, locationID uuid.UUID, limit, offset int) ([]*model.LocationItem, error) {
	db := r.getter.DefaultTrOrDB(ctx, r.pool)

	rawLocItems, err := r.q.ListLocationItems(
		ctx,
		db,
		sqlc.ListLocationItemsParams{
			LocationID: pgtype.UUID{
				Bytes: locationID,
				Valid: true,
			},
			Limit:  int32(limit),
			Offset: int32(offset),
		},
	)
	if err != nil {
		return nil, err
	}

	return mapper.MapSQLCToLocationItems(rawLocItems), nil
}
