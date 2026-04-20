package port

import (
	"backend/internal/domain/model"
	"context"

	"github.com/google/uuid"
)

type LocationItemRepository interface {
	Create(ctx context.Context, locationItem *model.LocationItem) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.LocationItem, error)
	GetByLocationAndItem(ctx context.Context, locationID, itemID uuid.UUID) (*model.LocationItem, error)
	Update(ctx context.Context, locationItem *model.LocationItem) error
	DeleteByItemID(ctx context.Context, itemID uuid.UUID) error
	DeleteByLocationID(ctx context.Context, locationID uuid.UUID) error
	List(ctx context.Context, locationID uuid.UUID, limit, offset int) ([]*model.LocationItem, error)
}
