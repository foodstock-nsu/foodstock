package port

import (
	"backend/internal/domain/model"
	"context"

	"github.com/google/uuid"
)

type TransactionRepository interface {
	Create(ctx context.Context, transaction *model.Transaction) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Transaction, error)
	GetBySbpID(ctx context.Context, sbpID string) (*model.Transaction, error)
	Update(ctx context.Context, transaction *model.Transaction) error
	List(ctx context.Context, orderID uuid.UUID) ([]*model.Transaction, error)
}
