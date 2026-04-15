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

func TestMapLocationToSQLCCreate(t *testing.T) {
	location, _ := model.NewLocation(
		"nsu_1",
		"Novosibirsk State University | Store №1",
		"Novosibirsk, some st., 6300019",
	)

	mapped := mapper.MapLocationToSQLCCreate(location)

	require.True(t, mapped.ID.Valid)
	require.True(t, mapped.CreatedAt.Valid)

	assert.Equal(t, [16]byte(location.ID()), mapped.ID.Bytes)
	assert.Equal(t, location.Slug(), mapped.Slug)
	assert.Equal(t, location.Name(), mapped.Name)
	assert.Equal(t, location.Address(), mapped.Address)
	assert.Equal(t, location.IsActive(), mapped.IsActive)
	assert.Equal(t, location.CreatedAt(), mapped.CreatedAt.Time)
}

func TestMapSQLCToLocation(t *testing.T) {
	raw := sqlc.Location{
		ID: pgtype.UUID{
			Bytes: uuid.New(),
			Valid: true,
		},
		Slug:     "nsu_1",
		Name:     "Novosibirsk State University | Store №1",
		Address:  "Novosibirsk, some st., 6300019",
		IsActive: true,
		CreatedAt: pgtype.Timestamptz{
			Time:             time.Now().UTC(),
			InfinityModifier: 0,
			Valid:            true,
		},
	}

	mapped := mapper.MapSQLCToLocation(raw)

	assert.Equal(t, raw.ID.Bytes, [16]byte(mapped.ID()))
	assert.Equal(t, raw.Slug, mapped.Slug())
	assert.Equal(t, raw.Name, mapped.Name())
	assert.Equal(t, raw.Address, mapped.Address())
	assert.Equal(t, raw.IsActive, mapped.IsActive())
	assert.Equal(t, raw.CreatedAt.Time, mapped.CreatedAt())
}

func TestMapLocationToSQLCUpdate(t *testing.T) {
	location, _ := model.NewLocation(
		"nsu_1",
		"Novosibirsk State University | Store №1",
		"Novosibirsk, some st., 6300019",
	)

	mapped := mapper.MapLocationToSQLCUpdate(location)

	require.True(t, mapped.ID.Valid)
	assert.Equal(t, [16]byte(location.ID()), mapped.ID.Bytes)
	assert.Equal(t, location.Slug(), mapped.Slug)
	assert.Equal(t, location.Name(), mapped.Name)
	assert.Equal(t, location.Address(), mapped.Address)
	assert.Equal(t, location.IsActive(), mapped.IsActive)
}

func TestMapSQLCToLocations(t *testing.T) {
	raw := []sqlc.Location{
		{
			ID: pgtype.UUID{
				Bytes: uuid.New(),
				Valid: true,
			},
			Slug:     "nsu_1",
			Name:     "Novosibirsk State University | Store №1",
			Address:  "Novosibirsk, some st., 6300019",
			IsActive: true,
			CreatedAt: pgtype.Timestamptz{
				Time:             time.Now().UTC(),
				InfinityModifier: 0,
				Valid:            true,
			},
		},
		{
			ID: pgtype.UUID{
				Bytes: uuid.New(),
				Valid: true,
			},
			Slug:     "nsu_2",
			Name:     "Novosibirsk State University | Store №2",
			Address:  "Novosibirsk, some st., 6300019",
			IsActive: true,
			CreatedAt: pgtype.Timestamptz{
				Time:             time.Now().UTC(),
				InfinityModifier: 0,
				Valid:            true,
			},
		},
	}

	mapped := mapper.MapSQLCToLocations(raw)
	require.Len(t, mapped, len(raw))
	assert.Equal(t, raw[0].ID.Bytes, [16]byte(mapped[0].ID()))
	assert.Equal(t, raw[1].ID.Bytes, [16]byte(mapped[1].ID()))
	assert.Equal(t, raw[0].Slug, mapped[0].Slug())
	assert.Equal(t, raw[1].Address, mapped[1].Address())
}
