package mapper_test

import (
	"backend/internal/adapter/out/postgres/mapper"
	"backend/internal/adapter/out/postgres/sqlc"
	"backend/internal/domain/model"
	pkgpostgres "backend/pkg/postgres"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMapOrderItemToSQLCCreate(t *testing.T) {
	orderID := uuid.New()
	itemID := uuid.New()
	amount := 3
	price := int64(45000)
	item := model.RestoreOrderItem(uuid.New(), itemID, amount, price)

	mapped := mapper.MapOrderItemToSQLCCreate(orderID, item)

	assert.Equal(t, [16]byte(item.ID()), mapped.ID.Bytes)
	assert.Equal(t, [16]byte(orderID), mapped.OrderID.Bytes)
	assert.Equal(t, [16]byte(itemID), mapped.ItemID.Bytes)
	assert.Equal(t, int32(amount), mapped.ItemAmount)
	assert.Equal(t, pkgpostgres.Int64ToNumeric(price, -2), mapped.PriceAtPurchase)
}

func TestMapOrderItemsToSQLCCreateBatch(t *testing.T) {
	orderID := uuid.New()
	itemID1 := uuid.New()
	itemID2 := uuid.New()
	item1 := model.RestoreOrderItem(uuid.New(), itemID1, 1, 100)
	item2 := model.RestoreOrderItem(uuid.New(), itemID2, 2, 200)
	items := []*model.OrderItem{item1, item2}

	mapped := mapper.MapOrderItemsToSQLCCreateBatch(orderID, items)

	require.Len(t, mapped, 2)
	assert.Equal(t, [16]byte(orderID), mapped[0].OrderID.Bytes)
	assert.Equal(t, [16]byte(itemID1), mapped[0].ItemID.Bytes)
	assert.Equal(t, [16]byte(itemID2), mapped[1].ItemID.Bytes)
}

func TestMapSQLCToOrderItem(t *testing.T) {
	id := uuid.New()
	itemID := uuid.New()
	raw := sqlc.OrderItem{
		ID:              pgtype.UUID{Bytes: id, Valid: true},
		ItemID:          pgtype.UUID{Bytes: itemID, Valid: true},
		ItemAmount:      5,
		PriceAtPurchase: pkgpostgres.Int64ToNumeric(125050, -2),
	}

	mapped := mapper.MapSQLCToOrderItem(raw)

	assert.Equal(t, id, mapped.ID())
	assert.Equal(t, itemID, mapped.ItemID())
	assert.Equal(t, 5, mapped.Amount())
	assert.Equal(t, int64(125050), mapped.PriceAtPurchase())
}

func TestMapSQLCToOrderItems(t *testing.T) {
	rawItems := []sqlc.OrderItem{
		{ID: pgtype.UUID{Bytes: uuid.New(), Valid: true}},
		{ID: pgtype.UUID{Bytes: uuid.New(), Valid: true}},
	}

	mapped := mapper.MapSQLCToOrderItems(rawItems)

	require.Len(t, mapped, 2)
	assert.Equal(t, rawItems[0].ID.Bytes, [16]byte(mapped[0].ID()))
	assert.Equal(t, rawItems[1].ID.Bytes, [16]byte(mapped[1].ID()))
}
