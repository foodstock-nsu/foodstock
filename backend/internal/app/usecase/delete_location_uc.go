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

type DeleteLocationUC struct {
	trManager    trm.Manager
	location     port.LocationRepository
	locationItem port.LocationItemRepository
}

func NewDeleteLocationUC(
	trManager trm.Manager,
	location port.LocationRepository,
	locationItem port.LocationItemRepository,
) *DeleteLocationUC {
	return &DeleteLocationUC{
		trManager:    trManager,
		location:     location,
		locationItem: locationItem,
	}
}

func (uc *DeleteLocationUC) Execute(ctx context.Context, in dto.DeleteLocationInput) error {
	err := uc.trManager.Do(ctx, func(txCtx context.Context) error {
		// Delete all inventory ot this location
		if err := uc.locationItem.DeleteByLocationID(txCtx, in.ID); err != nil {
			return ucerrs.Wrap(
				ucerrs.ErrDeleteLocationItemByLocationIDDB, err,
			)
		}

		// Get location
		location, getErr := uc.location.GetByID(txCtx, in.ID)
		if getErr != nil {
			if errors.Is(getErr, pkgerrs.ErrObjectNotFound) {
				return ucerrs.ErrLocationNotFound
			}
			return ucerrs.Wrap(ucerrs.ErrGetLocationByIDDB, getErr)
		}

		if delErr := location.Delete(); delErr != nil {
			return ucerrs.ErrLocationAlreadyDeleted
		}

		if delErr := uc.location.SoftDelete(txCtx, location); delErr != nil {
			return ucerrs.Wrap(ucerrs.ErrSoftDeleteLocationDB, delErr)
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
