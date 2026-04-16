package mapper

import (
	"backend/internal/adapter/out/postgres/sqlc"
	"backend/internal/domain/model"
	pkgpostgres "backend/pkg/postgres"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func MapOrderToSQLCCreate(order *model.Order) sqlc.CreateOrderParams {
	var paidAt pgtype.Timestamptz
	if order.PaidAt() != nil {
		paidAt = pgtype.Timestamptz{
			Time:             *order.PaidAt(),
			InfinityModifier: 0,
			Valid:            true,
		}
	}

	totalPrice := pkgpostgres.Int64ToNumeric(order.TotalPrice(), int32(-2))

	return sqlc.CreateOrderParams{
		ID: pgtype.UUID{
			Bytes: order.ID(),
			Valid: true,
		},
		LocationID: pgtype.UUID{
			Bytes: order.LocationID(),
			Valid: true,
		},
		Status:     sqlc.OrderStatus(order.Status()),
		TotalPrice: totalPrice,
		CreatedAt: pgtype.Timestamptz{
			Time:             order.CreatedAt(),
			InfinityModifier: 0,
			Valid:            true,
		},
		PaidAt: paidAt,
	}
}

func MapSQLCToOrder(raw sqlc.Order) *model.Order {
	var paidAt *time.Time
	if raw.PaidAt.Valid {
		paidAt = &raw.PaidAt.Time
	}

	totalPrice, _ := pkgpostgres.NumericToInt64(raw.TotalPrice, -2)

	return model.RestoreOrder(
		raw.ID.Bytes,
		raw.LocationID.Bytes,
		nil,
		model.OrderStatus(raw.Status),
		totalPrice,
		raw.CreatedAt.Time,
		paidAt,
	)
}

func MapOrderToSQLCUpdate(order *model.Order) sqlc.UpdateOrderParams {
	var paidAt pgtype.Timestamptz
	if order.PaidAt() != nil {
		paidAt = pgtype.Timestamptz{
			Time:             *order.PaidAt(),
			InfinityModifier: 0,
			Valid:            true,
		}
	}

	totalPrice := pkgpostgres.Int64ToNumeric(order.TotalPrice(), int32(-2))

	return sqlc.UpdateOrderParams{
		ID: pgtype.UUID{
			Bytes: order.ID(),
			Valid: true,
		},
		Status:     sqlc.OrderStatus(order.Status()),
		TotalPrice: totalPrice,
		PaidAt:     paidAt,
	}
}

func MapSQLCToOrders(raw []sqlc.Order) []*model.Order {
	items := make([]*model.Order, len(raw))
	for i := range items {
		items[i] = MapSQLCToOrder(raw[i])
	}
	return items
}
