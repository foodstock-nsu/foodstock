package usecase

import (
	"backend/internal/app/dto"
	ucerrs "backend/internal/app/errs"
	"backend/internal/app/mapper"
	"backend/internal/domain/model"
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
	/*
		Get each inventory item
		Update it
		Save fixes in database
	*/
	updatedItems := make([]*model.LocationItem, 0, len(in.Inventory))
	err := uc.trManager.Do(ctx, func(txCtx context.Context) error {
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

			updErr := locationItem.Update(
				inputItem.Price,
				inputItem.IsAvailable,
				inputItem.StockAmount,
			)
			if updErr != nil {
				return ucerrs.Wrap(ucerrs.ErrInvalidInput, updErr)
			}

			if updErr = uc.locationItem.Update(txCtx, locationItem); updErr != nil {
				if errors.Is(updErr, pkgerrs.ErrObjectNotFound) {
					return ucerrs.ErrLocationItemNotFound
				}
				return ucerrs.Wrap(
					ucerrs.ErrUpdateLocationItemDB, updErr,
				)
			}

			updatedItems = append(updatedItems, locationItem)
		}
		return nil
	})

	if err != nil {
		return dto.UpdateInventoryOutput{}, err
	}

	return dto.UpdateInventoryOutput{
		Inventory: mapper.MapDomainToInventoryItemListDTO(updatedItems),
	}, nil
}
