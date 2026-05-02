package port

import (
	"context"

	"github.com/google/uuid"
)

type PaymentGateway interface {
	Create(ctx context.Context, amount int64, returnURL string, orderID uuid.UUID) (string, string, error)
	GetStatus(ctx context.Context, externalID string) (string, error)
}
