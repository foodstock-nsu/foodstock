package mapper

import (
	"backend/internal/adapter/out/postgres/sqlc"
	"backend/internal/domain/model"
	pkgpostgres "backend/pkg/postgres"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func MapOrderItemToSQLCCreate(orderID uuid.UUID, orderItem *model.OrderItem) sqlc.CreateOrderItemParams {
	price := pkgpostgres.Int64ToNumeric(orderItem.PriceAtPurchase(), -2)

	return sqlc.CreateOrderItemParams{
		ID: pgtype.UUID{
			Bytes: orderItem.ID(),
			Valid: true,
		},
		OrderID: pgtype.UUID{
			Bytes: orderID,
			Valid: true,
		},
		ItemID: pgtype.UUID{
			Bytes: orderItem.ItemID(),
			Valid: true,
		},
		ItemAmount:      int32(orderItem.ItemAmount()),
		PriceAtPurchase: price,
	}
}

func MapOrderItemsToSQLCCreateBatch(orderID uuid.UUID, orderItems []*model.OrderItem) []sqlc.CreateOrderItemsBatchParams {
	res := make([]sqlc.CreateOrderItemsBatchParams, len(orderItems))
	for i, orderItem := range orderItems {
		price := pkgpostgres.Int64ToNumeric(orderItem.PriceAtPurchase(), -2)

		res[i] = sqlc.CreateOrderItemsBatchParams{
			ID: pgtype.UUID{
				Bytes: orderItem.ID(),
				Valid: true,
			},
			OrderID: pgtype.UUID{
				Bytes: orderID,
				Valid: true,
			},
			ItemID: pgtype.UUID{
				Bytes: orderItem.ItemID(),
				Valid: true,
			},
			ItemAmount:      int32(orderItem.ItemAmount()),
			PriceAtPurchase: price,
		}
	}
	return res
}

func MapSQLCToOrderItem(raw sqlc.OrderItem) *model.OrderItem {
	price, _ := pkgpostgres.NumericToInt64(raw.PriceAtPurchase, -2)
	return model.RestoreOrderItem(
		raw.ID.Bytes,
		raw.ItemID.Bytes,
		int(raw.ItemAmount),
		price,
	)
}

func MapSQLCToOrderItems(raw []sqlc.OrderItem) []*model.OrderItem {
	res := make([]*model.OrderItem, len(raw))
	for i := range res {
		res[i] = MapSQLCToOrderItem(raw[i])
	}
	return res
}
