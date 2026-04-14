package model_test

import (
	"backend/domain/model"
	pkgerrs "backend/pkg/errs"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTransaction(t *testing.T) {
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
			testName:         "Failure - invalid amount",
			orderID:          testOrderID,
			sbpTransactionID: testSBP,
			amount:           int64(-10000), // too short
			expect:           pkgerrs.ErrValueIsInvalid,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			tr, err := model.NewTransaction(
				tt.orderID, tt.sbpTransactionID, tt.amount,
			)
			if tt.expect != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.expect)
				assert.Nil(t, tr)
			} else {
				require.NoError(t, err)
				require.NotNil(t, tr)

				assert.NotEmpty(t, tr.ID())
				assert.Equal(t, tt.orderID, tr.OrderID())
				assert.Equal(t, tt.sbpTransactionID, tr.SBPTransactionID())
				assert.Equal(t, tt.amount, tr.Amount())
				assert.Equal(t, model.TransactionPending, tr.Status())
				assert.Nil(t, tr.WebhookReceivedAt())
				assert.False(t, tr.CreatedAt().After(time.Now().UTC()))
			}
		})
	}
}

func TestTransaction_Confirm(t *testing.T) {
	tr := model.RestoreTransaction(
		uuid.New(),
		uuid.New(),
		"aerfqw12348789asdfasd",
		int64(100),
		model.TransactionPending,
		nil,
		time.Now().UTC(),
	)

	err := tr.Confirm()
	assert.NoError(t, err)
	assert.Equal(t, model.TransactionSuccess, tr.Status())

	err = tr.Confirm()
	assert.Error(t, err)
}

func TestTransaction_Deny(t *testing.T) {
	tr := model.RestoreTransaction(
		uuid.New(),
		uuid.New(),
		"aerfqw12348789asdfasd",
		int64(100),
		model.TransactionPending,
		nil,
		time.Now().UTC(),
	)

	err := tr.Deny()
	assert.NoError(t, err)
	assert.Equal(t, model.TransactionFailed, tr.Status())

	err = tr.Deny()
	assert.Error(t, err)
}
