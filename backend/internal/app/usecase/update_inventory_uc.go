package usecase

import (
	"backend/internal/app/dto"
	ucerrs "backend/internal/app/errs"
	"backend/internal/domain/port"
	pkgerrs "backend/pkg/errs"
	"context"
	"errors"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
)

type UpdateInventoryUC struct {
	trManager    trm.Manager
	locationItem port.LocationItemRepository
}

func NewUpdateInventoryUC(
	trManager trm.Manager,
	locationItem port.LocationItemRepository,
) *UpdateInventoryUC {
	return &UpdateInventoryUC{
		trManager:    trManager,
		locationItem: locationItem,
	}
}

func (uc *UpdateInventoryUC) Execute(ctx context.Context, in dto.UpdateInventoryInput) (dto.UpdateInventoryOutput, error) {
	err := uc.trManager.Do(ctx, func(txCtx context.Context) error {
		/*
			Get each inventory item
			Update it
			Save fixes in database
		*/
		for _, inputItem := range in.Inventory {
			locationItem, getErr := uc.locationItem.GetByLocationAndItem(
				txCtx, in.LocationID, inputItem.ItemID,
			)
			if getErr != nil {
				if errors.Is(getErr, pkgerrs.ErrObjectNotFound) {
					return ucerrs.ErrLocationItemNotFound
				}
				return ucerrs.Wrap(
					ucerrs.ErrGetLocationItemByLocationAndItemDB, getErr,
				)
			}

			updErr := locationItem.Update()
		}
	})
}
