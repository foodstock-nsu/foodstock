package mapper

import (
	"backend/internal/adapter/out/postgres/sqlc"
	"backend/internal/domain/model"

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
		Proteins:    toNumeric(item.Nutrition().Proteins(), 1),
		Fats:        toNumeric(item.Nutrition().Fats(), 1),
		Carbs:       toNumeric(item.Nutrition().Carbs(), 1),
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

	return model.RestoreItem(
		raw.ID.Bytes,
		raw.Name,
		desc,
		model.ItemCategory(raw.Category),
		photo,
		model.RestoreNutrition(
			int(raw.Calories),
			numericToFloat(raw.Proteins),
			numericToFloat(raw.Fats),
			numericToFloat(raw.Carbs),
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
		Proteins:    toNumeric(item.Nutrition().Proteins(), 1),
		Fats:        toNumeric(item.Nutrition().Fats(), 1),
		Carbs:       toNumeric(item.Nutrition().Carbs(), 1),
	}
}

func MapSQLCToItems(raw []sqlc.Item) []*model.Item {
	items := make([]*model.Item, len(raw))
	for i := range items {
		items[i] = MapSQLCToItem(raw[i])
	}
	return items
}
