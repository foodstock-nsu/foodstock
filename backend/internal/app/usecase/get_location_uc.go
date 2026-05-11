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

type GetLocationUC struct {
	location port.LocationRepository
}

func NewGetLocationUC(location port.LocationRepository) *GetLocationUC {
	return &GetLocationUC{location: location}
}

func (uc *GetLocationUC) Execute(ctx context.Context, in dto.GetLocationInput) (dto.GetLocationOutput, error) {
	// Get location by slug and validate it
	location, err := uc.location.GetBySlug(ctx, in.Slug)
	if err != nil {
		if errors.Is(err, pkgerrs.ErrObjectNotFound) {
			return dto.GetLocationOutput{}, ucerrs.ErrLocationNotFound
		}
		return dto.GetLocationOutput{}, ucerrs.Wrap(
			ucerrs.ErrGetLocationBySlugDB, err,
		)
	}

	if location.IsDeleted() {
		return dto.GetLocationOutput{}, ucerrs.ErrLocationAlreadyDeleted
	}

	return dto.GetLocationOutput{
		Location: mapper.MapDomainToLocationDTO(location),
	}, nil
}
