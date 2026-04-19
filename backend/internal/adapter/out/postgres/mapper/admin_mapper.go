package mapper

import (
	"backend/internal/adapter/out/postgres/sqlc"
	"backend/internal/domain/model"

	"github.com/jackc/pgx/v5/pgtype"
)

func MapAdminToSQLCCreate(admin *model.Admin) sqlc.CreateAdminParams {
	return sqlc.CreateAdminParams{
		ID: pgtype.UUID{
			Bytes: admin.ID(),
			Valid: true,
		},
		Login:        admin.Login(),
		PasswordHash: admin.PasswordHash(),
		CreatedAt: pgtype.Timestamptz{
			Time:             admin.CreatedAt(),
			InfinityModifier: 0,
			Valid:            true,
		},
	}
}

func MapSQLCToAdmin(raw sqlc.Admin) *model.Admin {
	return model.RestoreAdmin(
		raw.ID.Bytes,
		raw.Login,
		raw.PasswordHash,
		raw.CreatedAt.Time,
	)
}
