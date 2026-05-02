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

type OrderRepository struct {
	q      *sqlc.Queries
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func NewOrderRepository(
	pgClient *pkgpostgres.Client,
	getter *trmpgx.CtxGetter,
) *OrderRepository {
	return &OrderRepository{
		q:      sqlc.New(),
		pool:   pgClient.Pool,
		getter: getter,
	}
}

func (r *OrderRepository) Create(ctx context.Context, order *model.Order) error {
	db := r.getter.DefaultTrOrDB(ctx, r.pool)
	params := mapper.MapOrderToSQLCCreate(order)

	if err := r.q.CreateOrder(ctx, db, params); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return pkgerrs.NewObjectAlreadyExistsErrorWithReason(
					"order", pgErr,
				)
			}
		}
		return err
	}

	return nil
}

func (r *OrderRepository) Get(ctx context.Context, id uuid.UUID) (*model.Order, error) {
	db := r.getter.DefaultTrOrDB(ctx, r.pool)

	rawOrder, err := r.q.GetOrder(
		ctx,
		db,
		pgtype.UUID{
			Bytes: id,
			Valid: true,
		},
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pkgerrs.NewObjectNotFoundError("order", id)
		}
		return nil, err
	}

	return mapper.MapSQLCToOrder(rawOrder), nil
}

func (r *OrderRepository) Update(ctx context.Context, order *model.Order) error {
	db := r.getter.DefaultTrOrDB(ctx, r.pool)
	params := mapper.MapOrderToSQLCUpdate(order)
	return r.q.UpdateOrder(ctx, db, params)
}

func (r *OrderRepository) ListByLocationID(ctx context.Context, locationID uuid.UUID, limit, offset int) ([]*model.Order, error) {
	db := r.getter.DefaultTrOrDB(ctx, r.pool)

	rawOrders, err := r.q.ListOrdersByLocationID(
		ctx,
		db,
		sqlc.ListOrdersByLocationIDParams{
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

	return mapper.MapSQLCToOrders(rawOrders), nil
}

func (r *OrderRepository) ListByStatus(ctx context.Context, status model.OrderStatus, limit, offset int) ([]*model.Order, error) {
	db := r.getter.DefaultTrOrDB(ctx, r.pool)

	rawOrders, err := r.q.ListOrdersByStatus(ctx, db, sqlc.OrderStatus(status))
	if err != nil {
		return nil, err
	}

	return mapper.MapSQLCToOrders(rawOrders), nil
}
