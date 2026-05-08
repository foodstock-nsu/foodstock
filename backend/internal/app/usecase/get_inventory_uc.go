package usecase

import (
	"backend/internal/app/dto"
	ucerrs "backend/internal/app/errs"
	"backend/internal/app/mapper"
	"backend/internal/domain/port"
	pkgerrs "backend/pkg/errs"
	"context"
	"errors"
)

type GetInventoryUC struct {
	location     port.LocationRepository
	locationItem port.LocationItemRepository
}

func NewGetInventoryUC(locationItem port.LocationItemRepository) *GetInventoryUC {
	return &GetInventoryUC{locationItem: locationItem}
}

func (uc *GetInventoryUC) Execute(ctx context.Context, in dto.GetInventoryInput) (dto.GetInventoryOutput, error) {
	// Get the location by slug and validate it
	location, err := uc.location.GetBySlug(ctx, in.Slug)
	if err != nil {
		if errors.Is(err, pkgerrs.ErrObjectNotFound) {
			return dto.GetInventoryOutput{}, ucerrs.ErrLocationNotFound
		}
		return dto.GetInventoryOutput{}, ucerrs.Wrap(
			ucerrs.ErrGetLocationBySlugDB, err,
		)
	}

	if location.IsDeleted() {
		return dto.GetInventoryOutput{}, ucerrs.ErrLocationAlreadyDeleted
	}

	// Get the inventory
	inventory, err := uc.locationItem.List(ctx, location.ID())
	if err != nil {
		return dto.GetInventoryOutput{}, ucerrs.Wrap(
			ucerrs.ErrListLocationItemsDB, err,
		)
	}
	return dto.GetInventoryOutput{
		Inventory: mapper.MapDomainToInventoryItemListDTO(inventory),
	}, nil
}
