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

func TestMapLocationItemToSQLCCreate(t *testing.T) {
	id := uuid.New()
	itemID := uuid.New()
	locationID := uuid.New()
	price := int64(10050)
	stock := 10

	item := model.RestoreLocationItem(
		id,
		itemID,
		locationID,
		price,
		true,
		stock,
	)

	mapped := mapper.MapLocationItemToSQLCCreate(item)

	require.True(t, mapped.ID.Valid)
	require.True(t, mapped.ItemID.Valid)
	require.True(t, mapped.LocationID.Valid)

	assert.Equal(t, [16]byte(item.ID()), mapped.ID.Bytes)
	assert.Equal(t, [16]byte(item.ItemID()), mapped.ItemID.Bytes)
	assert.Equal(t, [16]byte(item.LocationID()), mapped.LocationID.Bytes)
	assert.Equal(t, pkgpostgres.Int64ToNumeric(item.Price(), -2), mapped.Price)
	assert.Equal(t, item.IsAvailable(), mapped.IsAvailable)
	assert.Equal(t, int32(item.StockAmount()), mapped.StockAmount)
}

func TestMapSQLCToLocationItem(t *testing.T) {
	raw := sqlc.LocationItem{
		ID: pgtype.UUID{
			Bytes: uuid.New(),
			Valid: true,
		},
		ItemID: pgtype.UUID{
			Bytes: uuid.New(),
			Valid: true,
		},
		LocationID: pgtype.UUID{
			Bytes: uuid.New(),
			Valid: true,
		},
		Price:       pkgpostgres.Int64ToNumeric(15000, -2),
		IsAvailable: true,
		StockAmount: 5,
	}

	mapped := mapper.MapSQLCToLocationItem(raw)

	expectedPrice, _ := pkgpostgres.NumericToInt64(raw.Price, -2)

	assert.Equal(t, raw.ID.Bytes, [16]byte(mapped.ID()))
	assert.Equal(t, raw.ItemID.Bytes, [16]byte(mapped.ItemID()))
	assert.Equal(t, raw.LocationID.Bytes, [16]byte(mapped.LocationID()))
	assert.Equal(t, expectedPrice, mapped.Price())
	assert.Equal(t, raw.IsAvailable, mapped.IsAvailable())
	assert.Equal(t, int(raw.StockAmount), mapped.StockAmount())
}

func TestMapLocationItemToSQLCUpdate(t *testing.T) {
	id := uuid.New()
	price := int64(20000)
	stock := 25

	item := model.RestoreLocationItem(
		id,
		uuid.New(),
		uuid.New(),
		price,
		false,
		stock,
	)

	mapped := mapper.MapLocationItemToSQLCUpdate(item)

	require.True(t, mapped.ID.Valid)
	assert.Equal(t, [16]byte(item.ID()), mapped.ID.Bytes)
	assert.Equal(t, pkgpostgres.Int64ToNumeric(item.Price(), -2), mapped.Price)
	assert.Equal(t, int32(item.StockAmount()), mapped.StockAmount)
}

func TestMapSQLCToLocationItems(t *testing.T) {
	raw := []sqlc.LocationItem{
		{
			ID:          pgtype.UUID{Bytes: uuid.New(), Valid: true},
			Price:       pkgpostgres.Int64ToNumeric(1000, -2),
			StockAmount: 10,
		},
		{
			ID:          pgtype.UUID{Bytes: uuid.New(), Valid: true},
			Price:       pkgpostgres.Int64ToNumeric(2000, -2),
			StockAmount: 20,
		},
	}

	mapped := mapper.MapSQLCToLocationItems(raw)
	require.Len(t, mapped, len(raw))

	assert.Equal(t, raw[0].ID.Bytes, [16]byte(mapped[0].ID()))
	assert.Equal(t, raw[1].ID.Bytes, [16]byte(mapped[1].ID()))
	assert.Equal(t, int(raw[0].StockAmount), mapped[0].StockAmount())
	assert.Equal(t, int(raw[1].StockAmount), mapped[1].StockAmount())
}
