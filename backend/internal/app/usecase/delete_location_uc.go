package usecase

import (
	"backend/internal/app/dto"
	ucerrs "backend/internal/app/errs"
	"backend/internal/domain/port"
	pkgerrs "backend/pkg/errs"
	"context"
	"errors"
)

type DeleteLocationUC struct {
	location     port.LocationRepository
	locationItem port.LocationItemRepository
}

func NewDeleteLocationUC(
	location port.LocationRepository,
	locationItem port.LocationItemRepository,
) *DeleteLocationUC {
	return &DeleteLocationUC{
		location:     location,
		locationItem: locationItem,
	}
}

func (uc *DeleteLocationUC) Execute(ctx context.Context, in dto.DeleteLocationInput) error {
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
}
