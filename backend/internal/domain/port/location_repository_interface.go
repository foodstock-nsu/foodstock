package port

import (
	"backend/internal/domain/model"
	"context"

	"github.com/google/uuid"
)

type LocationRepository interface {
	Create(ctx context.Context, location *model.Location) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Location, error)
	GetBySlug(ctx context.Context, slug string) (*model.Location, error)
	Update(ctx context.Context, location *model.Location) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context) ([]*model.Location, error)
}
