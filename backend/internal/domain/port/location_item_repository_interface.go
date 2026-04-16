package port

import (
	"backend/internal/domain/model"
	"context"

	"github.com/google/uuid"
)

type LocationItemRepository interface {
	Create(ctx context.Context, locItem *model.LocationItem) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.LocationItem, error)
	GetByLocationAndItem(ctx context.Context, locationID, itemID uuid.UUID) (*model.LocationItem, error)
	Update(ctx context.Context, locItem *model.LocationItem) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, locationID uuid.UUID, limit, offset int) ([]*model.LocationItem, error)
}
