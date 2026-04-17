package mapper_test

import (
	"backend/internal/adapter/out/postgres/mapper"
	"backend/internal/adapter/out/postgres/sqlc"
	"backend/internal/domain/model"
	pkgpostgres "backend/pkg/postgres"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMapTransactionToSQLCCreate(t *testing.T) {
	id := uuid.New()
	orderID := uuid.New()
	sbpID := "sbp-12345"
	amount := int64(150050) // 1500.50
	status := model.TransactionSuccess
	createdAt := time.Now().UTC()
	paidAt := createdAt.Add(time.Minute)

	tx := model.RestoreTransaction(
		id,
		orderID,
		sbpID,
		amount,
		status,
		&paidAt,
		createdAt,
	)

	mapped := mapper.MapTransactionToSQLCCreate(tx)

	assert.Equal(t, [16]byte(id), mapped.ID.Bytes)
	assert.Equal(t, [16]byte(orderID), mapped.OrderID.Bytes)
	assert.Equal(t, sbpID, mapped.SbpTransactionID)
	assert.Equal(t, pkgpostgres.Int64ToNumeric(amount, -2), mapped.Amount)
	assert.Equal(t, sqlc.TransactionStatus(status), mapped.Status)
	assert.True(t, mapped.PaidAt.Valid)
	assert.True(t, paidAt.Equal(mapped.PaidAt.Time))
	assert.True(t, createdAt.Equal(mapped.CreatedAt.Time))
}

func TestMapTransactionToSQLCCreate_NilPaidAt(t *testing.T) {
	tx := model.RestoreTransaction(
		uuid.New(),
		uuid.New(),
		"sbp-id",
		1000,
		model.TransactionPending,
		nil,
		time.Now().UTC(),
	)

	mapped := mapper.MapTransactionToSQLCCreate(tx)

	assert.False(t, mapped.PaidAt.Valid)
}

func TestMapSQLCToTransaction(t *testing.T) {
	id := uuid.New()
	orderID := uuid.New()
	now := time.Now().UTC()

	raw := sqlc.Transaction{
		ID:               pgtype.UUID{Bytes: id, Valid: true},
		OrderID:          pgtype.UUID{Bytes: orderID, Valid: true},
		SbpTransactionID: "sbp-999",
		Amount:           pkgpostgres.Int64ToNumeric(250000, -2),
		Status:           sqlc.TransactionStatusFAILED,
		PaidAt:           pgtype.Timestamptz{Time: now.Add(time.Second), Valid: true},
		CreatedAt:        pgtype.Timestamptz{Time: now, Valid: true},
	}

	mapped := mapper.MapSQLCToTransaction(raw)

	assert.Equal(t, id, mapped.ID())
	assert.Equal(t, orderID, mapped.OrderID())
	assert.Equal(t, "sbp-999", mapped.SBPTransactionID())
	assert.Equal(t, int64(250000), mapped.Amount())
	assert.Equal(t, model.TransactionFailed, mapped.Status())
	assert.NotNil(t, mapped.PaidAt())
	assert.True(t, now.Add(time.Second).Equal(*mapped.PaidAt()))
	assert.True(t, now.Equal(mapped.CreatedAt()))
}

func TestMapSQLCToTransactions(t *testing.T) {
	raw := []sqlc.Transaction{
		{ID: pgtype.UUID{Bytes: uuid.New(), Valid: true}},
		{ID: pgtype.UUID{Bytes: uuid.New(), Valid: true}},
	}

	mapped := mapper.MapSQLCToTransactions(raw)

	require.Len(t, mapped, 2)
	assert.Equal(t, raw[0].ID.Bytes, [16]byte(mapped[0].ID()))
	assert.Equal(t, raw[1].ID.Bytes, [16]byte(mapped[1].ID()))
}
