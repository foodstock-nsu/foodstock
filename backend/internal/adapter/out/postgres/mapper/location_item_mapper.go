package mapper

import (
	"backend/internal/adapter/out/postgres/sqlc"
	"backend/internal/domain/model"
	pkgpostgres "backend/pkg/postgres"

	"github.com/jackc/pgx/v5/pgtype"
)

func MapLocationItemToSQLCCreate(item *model.LocationItem) sqlc.CreateLocationItemParams {
	price := pkgpostgres.Int64ToNumeric(item.Price(), int32(-2))

	return sqlc.CreateLocationItemParams{
		ID: pgtype.UUID{
			Bytes: item.ID(),
			Valid: true,
		},
		ItemID: pgtype.UUID{
			Bytes: item.ItemID(),
			Valid: true,
		},
		LocationID: pgtype.UUID{
			Bytes: item.LocationID(),
			Valid: true,
		},
		Price:       price,
		IsAvailable: item.IsAvailable(),
		StockAmount: int32(item.StockAmount()),
	}
}

func MapSQLCToLocationItem(raw sqlc.LocationItem) *model.LocationItem {
	price, _ := pkgpostgres.NumericToInt64(raw.Price, int32(-2))

	return model.RestoreLocationItem(
		raw.ID.Bytes,
		raw.ItemID.Bytes,
		raw.LocationID.Bytes,
		price,
		raw.IsAvailable,
		int(raw.StockAmount),
	)
}

func MapLocationItemToSQLCUpdate(item *model.LocationItem) sqlc.UpdateLocationItemParams {
	price := pkgpostgres.Int64ToNumeric(item.Price(), int32(-2))

	return sqlc.UpdateLocationItemParams{
		ID: pgtype.UUID{
			Bytes: item.ID(),
			Valid: true,
		},
		Price:       price,
		StockAmount: int32(item.StockAmount()),
	}
}

func MapSQLCToLocationItems(raw []sqlc.LocationItem) []*model.LocationItem {
	items := make([]*model.LocationItem, len(raw))
	for i := range items {
		items[i] = MapSQLCToLocationItem(raw[i])
	}
	return items
}
