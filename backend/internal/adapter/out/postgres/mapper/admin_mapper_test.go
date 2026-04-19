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

func TestMapAdminToSQLCCreate(t *testing.T) {
	admin := model.RestoreAdmin(
		uuid.New(),
		"admin",
		"hash",
		time.Now().UTC(),
	)

	mapped := mapper.MapAdminToSQLCCreate(admin)

	require.True(t, mapped.ID.Valid)
	require.True(t, mapped.CreatedAt.Valid)

	assert.Equal(t, [16]byte(admin.ID()), mapped.ID.Bytes)
	assert.Equal(t, admin.Login(), mapped.Login)
	assert.Equal(t, admin.PasswordHash(), mapped.PasswordHash)
	assert.Equal(t, admin.CreatedAt(), mapped.CreatedAt.Time)
}

func TestMapAdminToSQLCUpsert(t *testing.T) {
	admin := model.RestoreAdmin(
		uuid.New(),
		"admin",
		"hash",
		time.Now().UTC(),
	)

	mapped := mapper.MapAdminToSQLCUpsert(admin)

	require.True(t, mapped.ID.Valid)
	require.True(t, mapped.CreatedAt.Valid)

	assert.Equal(t, [16]byte(admin.ID()), mapped.ID.Bytes)
	assert.Equal(t, admin.Login(), mapped.Login)
	assert.Equal(t, admin.PasswordHash(), mapped.PasswordHash)
	assert.Equal(t, admin.CreatedAt(), mapped.CreatedAt.Time)
}

func TestMapSQLCToAdmin(t *testing.T) {
	raw := sqlc.Admin{
		ID: pgtype.UUID{
			Bytes: uuid.New(),
			Valid: true,
		},
		Login:        "admin",
		PasswordHash: "hash",
		CreatedAt: pgtype.Timestamptz{
			Time:             time.Now().UTC(),
			InfinityModifier: 0,
			Valid:            true,
		},
	}

	mapped := mapper.MapSQLCToAdmin(raw)

	assert.Equal(t, raw.ID.Bytes, [16]byte(mapped.ID()))
	assert.Equal(t, raw.Login, mapped.Login())
	assert.Equal(t, raw.PasswordHash, mapped.PasswordHash())
	assert.Equal(t, raw.CreatedAt.Time, mapped.CreatedAt())
}
