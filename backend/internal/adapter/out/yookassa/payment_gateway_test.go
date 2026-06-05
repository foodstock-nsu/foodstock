///go:build integration

package yookassa_test

import (
	"backend/internal/adapter/out/yookassa"
	"backend/internal/domain/model"
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPaymentGateway_Integration(t *testing.T) {
	var (
		shopID    = "1344511"
		apiKey    = "test_ivvXDKubIzQ-vo5_RMv5Z9a4zSQ9BHfhr7VybxhzabE"
		ctx       = context.Background()
		amount    = int64(10050)
		returnURL = "https://google.com"
		orderID   = uuid.New()
	)

	gateway := yookassa.NewPaymentGateway(shopID, apiKey, 15*time.Second)

	var externalID string

	t.Run("Create Payment Success", func(t *testing.T) {
		paymentID, payURL, err := gateway.Create(ctx, amount, returnURL, orderID)
		externalID = paymentID

		require.NoError(t, err)
		assert.NotEmpty(t, paymentID)
		assert.NotEmpty(t, payURL)

		t.Logf("Payment ID: %s", paymentID)
		t.Logf("Pay URL: %s", payURL)
	})

	t.Run("Get Status Success", func(t *testing.T) {
		status, err := gateway.GetStatus(ctx, externalID)
		require.NoError(t, err)
		assert.Equal(t, model.TransactionPending, status)
	})

	t.Run("Auth Error", func(t *testing.T) {
		badGateway := yookassa.NewPaymentGateway("wrong", "wrong", 5*time.Second)
		_, _, err := badGateway.Create(ctx, amount, returnURL, orderID)

		assert.ErrorIs(t, err, yookassa.ErrAuthInvalid)
	})
}
