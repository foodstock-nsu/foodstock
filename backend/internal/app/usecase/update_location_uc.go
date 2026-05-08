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

type UpdateLocationUC struct {
	location port.LocationRepository
}

func NewUpdateLocationUC(location port.LocationRepository) *UpdateLocationUC {
	return &UpdateLocationUC{location: location}
}

func (uc *UpdateLocationUC) Execute(ctx context.Context, in dto.UpdateLocationInput) (dto.UpdateLocationOutput, error) {
	// Get location by slug and validate it
	location, err := uc.location.GetBySlug(ctx, in.Slug)
	if err != nil {
		if errors.Is(err, pkgerrs.ErrObjectNotFound) {
			return dto.UpdateLocationOutput{}, ucerrs.ErrLocationNotFound
		}
		return dto.UpdateLocationOutput{}, ucerrs.Wrap(
			ucerrs.ErrGetLocationBySlugDB, err,
		)
	}

	if location.IsDeleted() {
		return dto.UpdateLocationOutput{}, ucerrs.ErrLocationAlreadyDeleted
	}

	// Update
	err = location.Update(in.Name, in.Address)
	if err != nil {
		return dto.UpdateLocationOutput{}, ucerrs.Wrap(
			ucerrs.ErrInvalidInput, err,
		)
	}

	if in.IsActive != nil {
		if *in.IsActive {
			location.Activate()
		} else {
			location.Deactivate()
		}
	}

	err = uc.location.Update(ctx, location)
	if err != nil {
		return dto.UpdateLocationOutput{}, ucerrs.Wrap(
			ucerrs.ErrUpdateLocationDB, err,
		)
	}

	return dto.UpdateLocationOutput{
		Location: mapper.MapDomainToLocationDTO(location),
	}, nil
}
