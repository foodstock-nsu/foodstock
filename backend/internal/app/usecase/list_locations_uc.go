package usecase

import (
	"backend/internal/app/dto"
	ucerrs "backend/internal/app/errs"
	"backend/internal/app/mapper"
	"backend/internal/domain/port"
	"context"
)

type ListLocationsUC struct {
	location port.LocationRepository
}

func NewListLocationsUC(location port.LocationRepository) *ListLocationsUC {
	return &ListLocationsUC{location: location}
}

func (uc *ListLocationsUC) Execute(ctx context.Context) (dto.ListLocationsOutput, error) {
	locations, err := uc.location.List(ctx)
	if err != nil {
		return dto.ListLocationsOutput{}, ucerrs.Wrap(
			ucerrs.ErrListLocationsDB, err,
		)
	}
	return dto.ListLocationsOutput{
		Locations: mapper.MapDomainToLocationListDTO(locations),
	}, nil
}
