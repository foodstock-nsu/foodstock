package mapper

import (
	"backend/internal/adapter/out/postgres/sqlc"
	"backend/internal/domain/model"
	pkgpostgres "backend/pkg/postgres"

	"github.com/jackc/pgx/v5/pgtype"
)

func MapItemToSQLCCreate(item *model.Item) sqlc.CreateItemParams {
	var (
		desc  pgtype.Text
		photo pgtype.Text
	)

	if item.Description() != nil {
		desc = pgtype.Text{
			String: *item.Description(),
			Valid:  true,
		}
	}
	if item.PhotoURL() != nil {
		photo = pgtype.Text{
			String: *item.PhotoURL(),
			Valid:  true,
		}
	}

	proteins, _ := pkgpostgres.Float64ToNumeric(item.Nutrition().Proteins(), 1)
	fats, _ := pkgpostgres.Float64ToNumeric(item.Nutrition().Fats(), 1)
	carbs, _ := pkgpostgres.Float64ToNumeric(item.Nutrition().Carbs(), 1)

	return sqlc.CreateItemParams{
		ID: pgtype.UUID{
			Bytes: item.ID(),
			Valid: true,
		},
		Name:        item.Name(),
		Description: desc,
		Category:    sqlc.ItemCategory(item.Category()),
		PhotoUrl:    photo,
		Calories:    int32(item.Nutrition().Calories()),
		Proteins:    proteins,
		Fats:        fats,
		Carbs:       carbs,
		CreatedAt: pgtype.Timestamptz{
			Time:             item.CreatedAt(),
			InfinityModifier: 0,
			Valid:            true,
		},
	}
}

func MapSQLCToItem(raw sqlc.Item) *model.Item {
	var (
		desc  *string
		photo *string
	)

	if raw.Description.Valid {
		desc = &raw.Description.String
	}
	if raw.PhotoUrl.Valid {
		photo = &raw.PhotoUrl.String
	}

	proteins, _ := pkgpostgres.NumericToFloat64(raw.Proteins)
	fats, _ := pkgpostgres.NumericToFloat64(raw.Fats)
	carbs, _ := pkgpostgres.NumericToFloat64(raw.Carbs)

	return model.RestoreItem(
		raw.ID.Bytes,
		raw.Name,
		desc,
		model.ItemCategory(raw.Category),
		photo,
		model.RestoreNutrition(
			int(raw.Calories),
			proteins,
			fats,
			carbs,
		),
		raw.CreatedAt.Time,
	)
}

func MapItemToSQLCUpdate(item *model.Item) sqlc.UpdateItemParams {
	var (
		desc  pgtype.Text
		photo pgtype.Text
	)

	if item.Description() != nil {
		desc = pgtype.Text{
			String: *item.Description(),
			Valid:  true,
		}
	}
	if item.PhotoURL() != nil {
		photo = pgtype.Text{
			String: *item.PhotoURL(),
			Valid:  true,
		}
	}

	proteins, _ := pkgpostgres.Float64ToNumeric(item.Nutrition().Proteins(), 1)
	fats, _ := pkgpostgres.Float64ToNumeric(item.Nutrition().Fats(), 1)
	carbs, _ := pkgpostgres.Float64ToNumeric(item.Nutrition().Carbs(), 1)

	return sqlc.UpdateItemParams{
		ID: pgtype.UUID{
			Bytes: item.ID(),
			Valid: true,
		},
		Name:        item.Name(),
		Description: desc,
		Category:    sqlc.ItemCategory(item.Category()),
		PhotoUrl:    photo,
		Calories:    int32(item.Nutrition().Calories()),
		Proteins:    proteins,
		Fats:        fats,
		Carbs:       carbs,
	}
}

func MapSQLCToItems(raw []sqlc.Item) []*model.Item {
	items := make([]*model.Item, len(raw))
	for i := range items {
		items[i] = MapSQLCToItem(raw[i])
	}
	return items
}
