package port

import (
	"backend/internal/domain/model"
	"context"

	"github.com/google/uuid"
)

type ItemRepository interface {
	Create(ctx context.Context, item *model.Item) error
	Get(ctx context.Context, id uuid.UUID) (*model.Item, error)
	Update(ctx context.Context, item *model.Item) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListAll(ctx context.Context, limit, offset int) ([]*model.Item, error)
	ListByIDs(ctx context.Context, ids []uuid.UUID) ([]*model.Item, error)
}
