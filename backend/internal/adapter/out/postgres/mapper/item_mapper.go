package mapper

import (
	"backend/internal/adapter/out/postgres/sqlc"
	"backend/internal/domain/model"
	pkgpostgres "backend/pkg/postgres"
	"backend/pkg/utils"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func MapItemToSQLCCreate(item *model.Item) sqlc.CreateItemParams {
	var (
		desc                  pgtype.Text
		calories              pgtype.Int4
		proteins, fats, carbs pgtype.Numeric
	)

	if item.Description() != nil {
		desc = pgtype.Text{
			String: *item.Description(),
			Valid:  true,
		}
	}

	if item.Nutrition() != nil {
		if item.Nutrition().Calories() != nil {
			calories = pgtype.Int4{
				Int32: int32(*item.Nutrition().Calories()),
				Valid: true,
			}
		}

		if item.Nutrition().Proteins() != nil {
			proteins, _ = pkgpostgres.Float64ToNumeric(*item.Nutrition().Proteins(), 1)
		}

		if item.Nutrition().Fats() != nil {
			fats, _ = pkgpostgres.Float64ToNumeric(*item.Nutrition().Fats(), 1)
		}

		if item.Nutrition().Carbs() != nil {
			carbs, _ = pkgpostgres.Float64ToNumeric(*item.Nutrition().Carbs(), 1)
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
		PhotoUrl:    item.PhotoURL(),
		Calories:    calories,
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
		desc      *string
		deletedAt *time.Time
	)

	if raw.Description.Valid {
		desc = &raw.Description.String
	}
	if raw.DeletedAt.Valid {
		deletedAt = &raw.DeletedAt.Time
	}

	proteins, _ := pkgpostgres.NumericToFloat64(raw.Proteins)
	fats, _ := pkgpostgres.NumericToFloat64(raw.Fats)
	carbs, _ := pkgpostgres.NumericToFloat64(raw.Carbs)

	return model.RestoreItem(
		raw.ID.Bytes,
		raw.Name,
		desc,
		model.ItemCategory(raw.Category),
		raw.PhotoUrl,
		model.RestoreNutrition(
			utils.VPtr(int(raw.Calories.Int32)),
			utils.VPtr(proteins),
			utils.VPtr(fats),
			utils.VPtr(carbs),
		),
		raw.CreatedAt.Time,
		deletedAt,
	)
}

func MapItemToSQLCUpdate(item *model.Item) sqlc.UpdateItemParams {
	var (
		desc                  pgtype.Text
		calories              pgtype.Int4
		proteins, fats, carbs pgtype.Numeric
	)

	if item.Description() != nil {
		desc = pgtype.Text{
			String: *item.Description(),
			Valid:  true,
		}
	}

	if item.Nutrition() != nil {
		if item.Nutrition().Calories() != nil {
			calories = pgtype.Int4{
				Int32: int32(*item.Nutrition().Calories()),
				Valid: true,
			}
		}

		if item.Nutrition().Proteins() != nil {
			proteins, _ = pkgpostgres.Float64ToNumeric(*item.Nutrition().Proteins(), 1)
		}

		if item.Nutrition().Fats() != nil {
			fats, _ = pkgpostgres.Float64ToNumeric(*item.Nutrition().Fats(), 1)
		}

		if item.Nutrition().Carbs() != nil {
			carbs, _ = pkgpostgres.Float64ToNumeric(*item.Nutrition().Carbs(), 1)
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
		PhotoUrl:    item.PhotoURL(),
		Calories:    calories,
		Proteins:    proteins,
		Fats:        fats,
		Carbs:       carbs,
	}
}

func MapItemToSQLCSoftDelete(item *model.Item) sqlc.DeleteItemSoftParams {
	var deletedAt pgtype.Timestamptz
	if item.DeletedAt() != nil {
		deletedAt = pgtype.Timestamptz{
			Time:  *item.DeletedAt(),
			Valid: true,
		}
	}
	return sqlc.DeleteItemSoftParams{
		ID: pgtype.UUID{
			Bytes: item.ID(),
			Valid: true,
		},
		DeletedAt: deletedAt,
	}
}

func MapSQLCToItems(raw []sqlc.Item) []*model.Item {
	items := make([]*model.Item, len(raw))
	for i := range items {
		items[i] = MapSQLCToItem(raw[i])
	}
	return items
}
