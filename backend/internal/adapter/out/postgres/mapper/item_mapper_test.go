package mapper_test

import (
	"backend/internal/adapter/out/postgres/mapper"
	"backend/internal/adapter/out/postgres/sqlc"
	"backend/internal/domain/model"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMapItemToSQLCCreate(t *testing.T) {
	desc := "Test Description"
	photo := "https://example.com/photo.jpg"
	nutrition := model.RestoreNutrition(100, 10.5, 5.5, 20.5)

	item := model.RestoreItem(
		uuid.New(),
		"Test Item",
		&desc,
		model.ItemLunch,
		&photo,
		nutrition,
		time.Now().UTC(),
	)

	mapped := mapper.MapItemToSQLCCreate(item)

	require.True(t, mapped.ID.Valid)
	require.True(t, mapped.CreatedAt.Valid)
	require.True(t, mapped.Description.Valid)
	require.True(t, mapped.PhotoUrl.Valid)

	assert.Equal(t, [16]byte(item.ID()), mapped.ID.Bytes)
	assert.Equal(t, item.Name(), mapped.Name)
	assert.Equal(t, *item.Description(), mapped.Description.String)
	assert.Equal(t, string(item.Category()), string(mapped.Category))
	assert.Equal(t, *item.PhotoURL(), mapped.PhotoUrl.String)
	assert.Equal(t, int32(item.Nutrition().Calories()), mapped.Calories)
	assert.Equal(t, item.CreatedAt(), mapped.CreatedAt.Time)
}

func TestMapSQLCToItem(t *testing.T) {
	raw := sqlc.Item{
		ID: pgtype.UUID{
			Bytes: uuid.New(),
			Valid: true,
		},
		Name: "Test Item",
		Description: pgtype.Text{
			String: "Test Description",
			Valid:  true,
		},
		Category: sqlc.ItemCategory("breakfast"),
		PhotoUrl: pgtype.Text{
			String: "https://example.com/photo.jpg",
			Valid:  true,
		},
		Calories: 200,
		Proteins: pgtype.Numeric{Valid: true},
		Fats:     pgtype.Numeric{Valid: true},
		Carbs:    pgtype.Numeric{Valid: true},
		CreatedAt: pgtype.Timestamptz{
			Time:  time.Now().UTC(),
			Valid: true,
		},
	}

	mapped := mapper.MapSQLCToItem(raw)

	assert.Equal(t, raw.ID.Bytes, [16]byte(mapped.ID()))
	assert.Equal(t, raw.Name, mapped.Name())
	assert.Equal(t, raw.Description.String, *mapped.Description())
	assert.Equal(t, string(raw.Category), string(mapped.Category()))
	assert.Equal(t, raw.PhotoUrl.String, *mapped.PhotoURL())
	assert.Equal(t, int(raw.Calories), mapped.Nutrition().Calories())
	assert.Equal(t, raw.CreatedAt.Time, mapped.CreatedAt())
}

func TestMapItemToSQLCUpdate(t *testing.T) {
	desc := "Update Description"
	nutrition := model.RestoreNutrition(150, 12.0, 6.0, 25.0)

	item := model.RestoreItem(
		uuid.New(),
		"Update Item",
		&desc,
		model.ItemDrinks,
		nil,
		nutrition,
		time.Now().UTC(),
	)

	mapped := mapper.MapItemToSQLCUpdate(item)

	require.True(t, mapped.ID.Valid)
	assert.Equal(t, [16]byte(item.ID()), mapped.ID.Bytes)
	assert.Equal(t, item.Name(), mapped.Name)
	assert.Equal(t, *item.Description(), mapped.Description.String)
	assert.False(t, mapped.PhotoUrl.Valid)
	assert.Equal(t, int32(item.Nutrition().Calories()), mapped.Calories)
}

func TestMapSQLCToItems(t *testing.T) {
	raw := []sqlc.Item{
		{
			ID:        pgtype.UUID{Bytes: uuid.New(), Valid: true},
			Name:      "Item 1",
			Category:  sqlc.ItemCategory("lunch"),
			CreatedAt: pgtype.Timestamptz{Time: time.Now().UTC(), Valid: true},
		},
		{
			ID:        pgtype.UUID{Bytes: uuid.New(), Valid: true},
			Name:      "Item 2",
			Category:  sqlc.ItemCategory("breakfast"),
			CreatedAt: pgtype.Timestamptz{Time: time.Now().UTC(), Valid: true},
		},
	}

	mapped := mapper.MapSQLCToItems(raw)
	require.Len(t, mapped, len(raw))
	assert.Equal(t, raw[0].ID.Bytes, [16]byte(mapped[0].ID()))
	assert.Equal(t, raw[1].ID.Bytes, [16]byte(mapped[1].ID()))
	assert.Equal(t, raw[0].Name, mapped[0].Name())
	assert.Equal(t, string(raw[1].Category), string(mapped[1].Category()))
}
