package port

import (
	"backend/internal/domain/model"
	"context"

	"github.com/google/uuid"
)

type OrderItemRepository interface {
	Create(ctx context.Context, orderID uuid.UUID, orderItem *model.OrderItem) error
	CreateMany(ctx context.Context, orderID uuid.UUID, orderItems []*model.OrderItem) error
	List(ctx context.Context, orderID uuid.UUID) ([]*model.OrderItem, error)
}
