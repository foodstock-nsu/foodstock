package port

import (
	"backend/internal/domain/model"
	"context"

	"github.com/google/uuid"
)

type AdminRepository interface {
	Create(ctx context.Context, admin *model.Admin) error
	Upsert(ctx context.Context, admin *model.Admin) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Admin, error)
	GetByLogin(ctx context.Context, login string) (*model.Admin, error)
}
