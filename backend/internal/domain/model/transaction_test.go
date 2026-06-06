package model_test

import (
	"backend/internal/domain/model"
	pkgerrs "backend/pkg/errs"
	"backend/pkg/utils"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTransactionAndIsPending(t *testing.T) {
	var (
		testOrderID = uuid.New()
		testSBP     = "aerfqw12348789asdfasd"
		testAmount  = int64(10000)
	)

	type testCase struct {
		testName         string
		orderID          uuid.UUID
		sbpTransactionID string
		amount           int64
		expect           error
	}

	var testCases = []testCase{
		{
			testName:         "Success",
			orderID:          testOrderID,
			sbpTransactionID: testSBP,
			amount:           testAmount,
			expect:           nil,
		},
		{
			testName: "Failure - invalid order id",
			orderID:  uuid.Nil,
			expect:   pkgerrs.ErrValueIsInvalid,
		},
		{
			testName:         "Failure - empty sbp transaction id",
			orderID:          testOrderID,
			sbpTransactionID: "",
			expect:           pkgerrs.ErrValueIsRequired,
		},
		{
			testName:         "Failure - invalid items",
			orderID:          testOrderID,
			sbpTransactionID: testSBP,
			amount:           int64(-10000), // too short
			expect:           pkgerrs.ErrValueIsInvalid,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			tx, err := model.NewTransaction(
				tt.orderID, tt.sbpTransactionID, tt.amount,
			)
			if tt.expect != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.expect)
				assert.Nil(t, tx)
			} else {
				require.NoError(t, err)
				require.NotNil(t, tx)

				assert.NotEmpty(t, tx.ID())
				assert.Equal(t, tt.orderID, tx.OrderID())
				assert.Equal(t, tt.sbpTransactionID, tx.SBPTransactionID())
				assert.Equal(t, tt.amount, tx.Amount())
				assert.Equal(t, model.TransactionPending, tx.Status())
				assert.True(t, tx.IsPending())
				assert.Nil(t, tx.PaidAt())
				assert.Nil(t, tx.RefundedAt())
				assert.False(t, tx.CreatedAt().After(time.Now().UTC()))
			}
		})
	}
}

func TestTransaction_ConfirmAndIsConfirmed(t *testing.T) {
	tx := model.RestoreTransaction(
		uuid.New(),
		uuid.New(),
		"aerfqw12348789asdfasd",
		int64(100),
		model.TransactionPending,
		nil,
		nil,
		time.Now().UTC(),
	)

	assert.False(t, tx.IsConfirmed())

	err := tx.Confirm()
	assert.NoError(t, err)
	assert.True(t, tx.IsConfirmed())

	err = tx.Confirm()
	assert.Error(t, err)
}

func TestTransaction_DenyAndIsDenied(t *testing.T) {
	tx := model.RestoreTransaction(
		uuid.New(),
		uuid.New(),
		"aerfqw12348789asdfasd",
		int64(100),
		model.TransactionPending,
		nil,
		nil,
		time.Now().UTC(),
	)

	assert.False(t, tx.IsDenied())

	err := tx.Deny()
	assert.NoError(t, err)
	assert.True(t, tx.IsDenied())

	err = tx.Deny()
	assert.Error(t, err)
}

func TestTransaction_RefundAndIsRefunded(t *testing.T) {
	tx := model.RestoreTransaction(
		uuid.New(),
		uuid.New(),
		"aerfqw12348789asdfasd",
		int64(100),
		model.TransactionSuccess,
		utils.VPtr(time.Now().Add(-1*time.Minute)),
		nil,
		time.Now().Add(-2*time.Minute),
	)

	assert.False(t, tx.IsRefunded())

	err := tx.Refund()
	assert.NoError(t, err)
	assert.True(t, tx.IsRefunded())

	err = tx.Refund()
	assert.Error(t, err)
}
