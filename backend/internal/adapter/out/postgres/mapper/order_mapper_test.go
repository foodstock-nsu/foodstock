package mapper_test

import (
	"backend/internal/adapter/out/postgres/mapper"
	"backend/internal/adapter/out/postgres/sqlc"
	"backend/internal/domain/model"
	pkgpostgres "backend/pkg/postgres"
	"backend/pkg/utils"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMapOrderToSQLCCreate(t *testing.T) {
	now := time.Now().UTC()
	paidAt := now.Add(time.Minute)
	order := model.RestoreOrder(
		uuid.New(),
		uuid.New(),
		nil,
		model.OrderPaid,
		1500,
		now,
		&paidAt,
	)

	mapped := mapper.MapOrderToSQLCCreate(order)

	require.True(t, mapped.ID.Valid)
	require.True(t, mapped.LocationID.Valid)
	require.True(t, mapped.CreatedAt.Valid)
	require.True(t, mapped.PaidAt.Valid)

	assert.Equal(t, [16]byte(order.ID()), mapped.ID.Bytes)
	assert.Equal(t, [16]byte(order.LocationID()), mapped.LocationID.Bytes)
	assert.Equal(t, order.Status().String(), string(mapped.Status))
	assert.Equal(t, order.CreatedAt(), mapped.CreatedAt.Time)
	assert.Equal(t, *order.PaidAt(), mapped.PaidAt.Time)

	expectedPrice := pkgpostgres.Int64ToNumeric(1500, -2)
	assert.Equal(t, expectedPrice, mapped.TotalPrice)
}

func TestMapSQLCToOrder(t *testing.T) {
	now := time.Now().UTC()
	raw := sqlc.Order{
		ID:         pgtype.UUID{Bytes: uuid.New(), Valid: true},
		LocationID: pgtype.UUID{Bytes: uuid.New(), Valid: true},
		Status:     sqlc.OrderStatusPENDING,
		TotalPrice: pkgpostgres.Int64ToNumeric(2500, -2),
		CreatedAt:  pgtype.Timestamptz{Time: now, Valid: true},
		PaidAt:     pgtype.Timestamptz{Time: now.Add(time.Second), Valid: true},
	}

	mapped := mapper.MapSQLCToOrder(raw)

	assert.Equal(t, raw.ID.Bytes, [16]byte(mapped.ID()))
	assert.Equal(t, raw.LocationID.Bytes, [16]byte(mapped.LocationID()))
	assert.Equal(t, string(raw.Status), mapped.Status().String())
	assert.Equal(t, int64(2500), mapped.TotalPrice())
	assert.Equal(t, raw.CreatedAt.Time, mapped.CreatedAt())
	assert.True(t, raw.PaidAt.Valid)
	assert.Equal(t, raw.PaidAt.Time, *mapped.PaidAt())
}

func TestMapOrderToSQLCUpdate(t *testing.T) {
	order := model.RestoreOrder(
		uuid.New(),
		uuid.New(),
		nil,
		model.OrderCancelled,
		0,
		time.Now().UTC(),
		utils.VPtr(time.Now().UTC().Add(time.Minute)),
	)

	mapped := mapper.MapOrderToSQLCUpdate(order)

	require.True(t, mapped.ID.Valid)
	assert.Equal(t, [16]byte(order.ID()), mapped.ID.Bytes)
	assert.Equal(t, order.Status().String(), string(mapped.Status))
	assert.True(t, mapped.PaidAt.Valid)

	expectedPrice := pkgpostgres.Int64ToNumeric(0, -2)
	assert.Equal(t, expectedPrice, mapped.TotalPrice)
}

func TestMapSQLCToOrders(t *testing.T) {
	raw := []sqlc.Order{
		{
			ID:         pgtype.UUID{Bytes: uuid.New(), Valid: true},
			LocationID: pgtype.UUID{Bytes: uuid.New(), Valid: true},
			Status:     sqlc.OrderStatusPENDING,
			CreatedAt:  pgtype.Timestamptz{Time: time.Now().UTC(), Valid: true},
		},
		{
			ID:         pgtype.UUID{Bytes: uuid.New(), Valid: true},
			LocationID: pgtype.UUID{Bytes: uuid.New(), Valid: true},
			Status:     sqlc.OrderStatusPAID,
			CreatedAt:  pgtype.Timestamptz{Time: time.Now().UTC(), Valid: true},
		},
	}

	mapped := mapper.MapSQLCToOrders(raw)

	require.Len(t, mapped, len(raw))
	assert.Equal(t, raw[0].ID.Bytes, [16]byte(mapped[0].ID()))
	assert.Equal(t, raw[1].ID.Bytes, [16]byte(mapped[1].ID()))
	assert.Equal(t, string(raw[0].Status), mapped[0].Status().String())
	assert.Equal(t, string(raw[1].Status), mapped[1].Status().String())
}
