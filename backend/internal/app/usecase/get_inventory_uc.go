package usecase

import (
	"backend/internal/app/dto"
	ucerrs "backend/internal/app/errs"
	"backend/internal/app/mapper"
	"backend/internal/domain/port"
	"context"
)

type GetInventoryUC struct {
	locationItem port.LocationItemRepository
}

func NewGetInventoryUC(locationItem port.LocationItemRepository) *GetInventoryUC {
	return &GetInventoryUC{locationItem: locationItem}
}

func (uc *GetInventoryUC) Execute(ctx context.Context, in dto.GetInventoryInput) (dto.GetInventoryOutput, error) {
	inventory, err := uc.locationItem.List(ctx, in.LocationID)
	if err != nil {
		return dto.GetInventoryOutput{}, ucerrs.Wrap(
			ucerrs.ErrListLocationItemsDB, err,
		)
	}
	return dto.GetInventoryOutput{
		Inventory: mapper.MapDomainToInventoryItemListDTO(inventory),
	}, nil
}
