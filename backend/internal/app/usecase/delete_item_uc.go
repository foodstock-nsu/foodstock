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
		// Delete item in each catalog
		deleteErr := uc.locationItem.DeleteByItemID(ctx, in.ID)
		if deleteErr != nil {
			return ucerrs.Wrap(
				ucerrs.ErrDeleteLocationItemsByItemIDDB, deleteErr,
			)
		}

		// Delete item in items list
		deleteErr = uc.item.Delete(ctx, in.ID)
		if deleteErr != nil {
			if errors.Is(deleteErr, pkgerrs.ErrObjectNotFound) {
				return ucerrs.ErrItemNotFound
			}
			return ucerrs.Wrap(ucerrs.ErrDeleteItemDB, deleteErr)
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
