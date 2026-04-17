package mapper

import (
	"backend/internal/adapter/out/postgres/sqlc"
	"backend/internal/domain/model"

	"github.com/jackc/pgx/v5/pgtype"
)

func MapLocationToSQLCCreate(loc *model.Location) sqlc.CreateLocationParams {
	return sqlc.CreateLocationParams{
		ID: pgtype.UUID{
			Bytes: loc.ID(),
			Valid: true,
		},
		Slug:     loc.Slug(),
		Name:     loc.Name(),
		Address:  loc.Address(),
		IsActive: loc.IsActive(),
		CreatedAt: pgtype.Timestamptz{
			Time:             loc.CreatedAt(),
			InfinityModifier: 0,
			Valid:            true,
		},
	}
}

func MapSQLCToLocation(raw sqlc.Location) *model.Location {
	return model.RestoreLocation(
		raw.ID.Bytes,
		raw.Slug,
		raw.Name,
		raw.Address,
		raw.IsActive,
		raw.CreatedAt.Time,
	)
}

func MapLocationToSQLCUpdate(loc *model.Location) sqlc.UpdateLocationParams {
	return sqlc.UpdateLocationParams{
		ID: pgtype.UUID{
			Bytes: loc.ID(),
			Valid: true,
		},
		Slug:     loc.Slug(),
		Name:     loc.Name(),
		Address:  loc.Address(),
		IsActive: loc.IsActive(),
	}
}

func MapSQLCToLocations(raw []sqlc.Location) []*model.Location {
	locs := make([]*model.Location, len(raw))
	for i := range locs {
		locs[i] = MapSQLCToLocation(raw[i])
	}
	return locs
}
