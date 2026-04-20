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
	err := uc.trManager.Do(ctx, func(ctx context.Context) error {
		// Delete location
		if err := uc.location.Delete(ctx, in.ID); err != nil {
			if errors.Is(err, pkgerrs.ErrObjectNotFound) {
				return ucerrs.ErrLocationNotFound
			}
			return ucerrs.Wrap(ucerrs.ErrDeleteLocationDB, err)
		}

		// Delete all inventory ot this location
		if err := uc.locationItem.DeleteByLocationID(ctx, in.ID); err != nil {
			return ucerrs.Wrap(
				ucerrs.ErrDeleteLocationItemByLocationIDDB, err,
			)
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
