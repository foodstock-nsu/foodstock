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
)

type CreateLocationUC struct {
	location port.LocationRepository
}

func (uc *CreateLocationUC) Execute(ctx context.Context, in dto.CreateLocationInput) (dto.CreateLocationOutput, error) {
	// Rich model with validation
	location, err := model.NewLocation(in.Slug, in.Name, in.Address)
	if err != nil {
		return dto.CreateLocationOutput{}, ucerrs.Wrap(
			ucerrs.ErrInvalidInput, err,
		)
	}

	// Save it into database
	err = uc.location.Create(ctx, location)
	if err != nil {
		if errors.Is(err, pkgerrs.ErrObjectAlreadyExists) {
			return dto.CreateLocationOutput{}, ucerrs.ErrLocationAlreadyExists
		}
		return dto.CreateLocationOutput{}, ucerrs.Wrap(
			ucerrs.ErrCreateLocationDB, err,
		)
	}

	return dto.CreateLocationOutput{
		Location: mapper.MapDomainToLocationDTO(location),
	}, nil
}
