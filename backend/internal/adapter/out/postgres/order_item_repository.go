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
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderItemRepository struct {
	q      *sqlc.Queries
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func NewOrderItemRepository(
	pgClient *pkgpostgres.Client,
	getter *trmpgx.CtxGetter,
) *OrderItemRepository {
	return &OrderItemRepository{
		q:      sqlc.New(),
		pool:   pgClient.Pool,
		getter: getter,
	}
}

func (r *OrderItemRepository) Create(ctx context.Context, orderID uuid.UUID, orderItem *model.OrderItem) error {
	db := r.getter.DefaultTrOrDB(ctx, r.pool)
	params := mapper.MapOrderItemToSQLCCreate(orderID, orderItem)

	if err := r.q.CreateOrderItem(ctx, db, params); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return pkgerrs.NewObjectAlreadyExistsErrorWithReason(
					"order_item", pgErr,
				)
			}
		}
		return err
	}

	return nil
}

func (r *OrderItemRepository) CreateMany(ctx context.Context, orderID uuid.UUID, orderItems []*model.OrderItem) error {
	db := r.getter.DefaultTrOrDB(ctx, r.pool)
	params := mapper.MapOrderItemsToSQLCCreateBatch(orderID, orderItems)

	if _, err := r.q.CreateOrderItemsBatch(ctx, db, params); err != nil {
		return err
	}

	return nil
}

func (r *OrderItemRepository) List(ctx context.Context, orderID uuid.UUID) ([]*model.OrderItem, error) {
	db := r.getter.DefaultTrOrDB(ctx, r.pool)

	rawOrderItems, err := r.q.ListOrderItems(
		ctx,
		db,
		pgtype.UUID{
			Bytes: orderID,
			Valid: true,
		},
	)
	if err != nil {
		return nil, err
	}

	return mapper.MapSQLCToOrderItems(rawOrderItems), nil
}
