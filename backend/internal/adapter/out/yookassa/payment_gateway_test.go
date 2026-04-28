package yookassa_test

import (
	"backend/internal/adapter/out/yookassa"
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

	t.Run("CreatePayment Success", func(t *testing.T) {
		paymentID, payURL, err := gateway.Create(ctx, amount, returnURL, orderID)

		require.NoError(t, err)
		assert.NotEmpty(t, paymentID)
		assert.NotEmpty(t, payURL)

		t.Logf("Payment ID: %s", paymentID)
		t.Logf("Pay URL: %s", payURL)

		t.Run("GetStatus Success", func(t *testing.T) {
			status, err := gateway.GetStatus(ctx, paymentID)
			require.NoError(t, err)
			assert.Equal(t, "pending", status)
		})
	})

	t.Run("Auth Error", func(t *testing.T) {
		badGateway := yookassa.NewPaymentGateway("wrong", "wrong", 5*time.Second)
		_, _, err := badGateway.Create(ctx, amount, returnURL, orderID)

		assert.ErrorIs(t, err, yookassa.ErrAuthInvalid)
	})
}
