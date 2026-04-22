package port

import (
	"backend/internal/domain/model"
	"context"

	"github.com/google/uuid"
)

type OrderRepository interface {
	Create(ctx context.Context, order *model.Order) error
	Get(ctx context.Context, id uuid.UUID) (*model.Order, error)
	Update(ctx context.Context, order *model.Order) error
	ListByLocationID(ctx context.Context, locationID uuid.UUID, limit, offset int) ([]*model.Order, error)
	ListByStatus(ctx context.Context, status model.OrderStatus, limit, offset int) ([]*model.Order, error)
}
