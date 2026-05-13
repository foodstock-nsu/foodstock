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

type DeleteItemUC struct {
	trManager    trm.Manager
	item         port.ItemRepository
	locationItem port.LocationItemRepository
}

func NewDeleteItemUC(
	trManager trm.Manager,
	item port.ItemRepository,
	locationItem port.LocationItemRepository,
) *DeleteItemUC {
	return &DeleteItemUC{
		trManager:    trManager,
		item:         item,
		locationItem: locationItem,
	}
}

func (uc *DeleteItemUC) Execute(ctx context.Context, in dto.DeleteItemInput) error {
	err := uc.trManager.Do(ctx, func(ctx context.Context) error {
		// Get item by id
		item, getErr := uc.item.Get(ctx, in.ID)
		if getErr != nil {
			if errors.Is(getErr, pkgerrs.ErrObjectNotFound) {
				return ucerrs.ErrItemNotFound
			}
			return ucerrs.Wrap(ucerrs.ErrGetItemDB, getErr)
		}

		// Delete it in items list
		delErr := item.Delete()
		if delErr != nil {
			return ucerrs.ErrItemAlreadyDeleted
		}

		if delErr = uc.item.SoftDelete(ctx, item); delErr != nil {
			return ucerrs.Wrap(ucerrs.ErrSoftDeleteItemDB, delErr)
		}

		// Delete item in each catalog
		if delErr = uc.locationItem.DeleteByItemID(ctx, item.ID()); delErr != nil {
			return ucerrs.Wrap(ucerrs.ErrDeleteLocationItemsByItemIDDB, delErr)
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
