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

type ItemRepository struct {
	q      *sqlc.Queries
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func NewItemRepository(
	pgClient *pkgpostgres.Client,
	getter *trmpgx.CtxGetter,
) *ItemRepository {
	return &ItemRepository{
		q:      sqlc.New(),
		pool:   pgClient.Pool,
		getter: getter,
	}
}

func (r *ItemRepository) Create(ctx context.Context, item *model.Item) error {
	db := r.getter.DefaultTrOrDB(ctx, r.pool)
	params := mapper.MapItemToSQLCCreate(item)

	if err := r.q.CreateItem(ctx, db, params); err != nil {
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

func (r *ItemRepository) Get(ctx context.Context, id uuid.UUID) (*model.Item, error) {
	db := r.getter.DefaultTrOrDB(ctx, r.pool)

	rawItem, err := r.q.GetItem(
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

	return mapper.MapSQLCToItem(rawItem), nil
}

func (r *ItemRepository) Update(ctx context.Context, item *model.Item) error {
	db := r.getter.DefaultTrOrDB(ctx, r.pool)
	params := mapper.MapItemToSQLCUpdate(item)
	return r.q.UpdateItem(ctx, db, params)
}

func (r *ItemRepository) Delete(ctx context.Context, id uuid.UUID) error {
	db := r.getter.DefaultTrOrDB(ctx, r.pool)
	return r.q.DeleteItem(
		ctx,
		db,
		pgtype.UUID{
			Bytes: id,
			Valid: true,
		},
	)
}

func (r *ItemRepository) List(ctx context.Context, limit, offset int32) ([]*model.Item, error) {
	db := r.getter.DefaultTrOrDB(ctx, r.pool)

	rawItems, err := r.q.ListItems(
		ctx,
		db,
		sqlc.ListItemsParams{
			Limit:  limit,
			Offset: offset,
		},
	)
	if err != nil {
		return nil, err
	}

	return mapper.MapSQLCToItems(rawItems), nil
}
