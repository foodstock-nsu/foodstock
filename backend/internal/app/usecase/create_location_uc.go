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

type CreateLocationUC struct {
	trManager    trm.Manager
	location     port.LocationRepository
	item         port.ItemRepository
	locationItem port.LocationItemRepository
}

func NewCreateLocationUC(
	trManager trm.Manager,
	location port.LocationRepository,
	item port.ItemRepository,
	locationItem port.LocationItemRepository,
) *CreateLocationUC {
	return &CreateLocationUC{
		trManager:    trManager,
		location:     location,
		item:         item,
		locationItem: locationItem,
	}
}

func (uc *CreateLocationUC) Execute(ctx context.Context, in dto.CreateLocationInput) (dto.CreateLocationOutput, error) {
	// Rich model with validation
	location, err := model.NewLocation(in.Slug, in.Name, in.Address)
	if err != nil {
		return dto.CreateLocationOutput{}, ucerrs.Wrap(
			ucerrs.ErrInvalidInput, err,
		)
	}

	err = uc.trManager.Do(ctx, func(txCtx context.Context) error {
		// Save the location into database
		createErr := uc.location.Create(txCtx, location)
		if createErr != nil {
			if errors.Is(createErr, pkgerrs.ErrObjectAlreadyExists) {
				return ucerrs.ErrLocationAlreadyExists
			}
			return ucerrs.Wrap(ucerrs.ErrCreateLocationDB, createErr)
		}

		// Get all items and create catalog for this location
		items, listErr := uc.item.ListAll(txCtx)
		if listErr != nil {
			return ucerrs.Wrap(ucerrs.ErrListAllItemsDB, listErr)
		}

		for _, item := range items {
			locationItem, locItemErr := model.NewLocationItem(
				item.ID(),
				location.ID(),
				0,
				0,
			)
			if locItemErr != nil {
				return ucerrs.Wrap(
					ucerrs.ErrInvalidInput, locItemErr,
				)
			}

			createErr = uc.locationItem.Create(txCtx, locationItem)
			if createErr != nil {
				return ucerrs.Wrap(ucerrs.ErrCreateLocationItemDB, err)
			}
		}

		return nil
	})

	if err != nil {
		return dto.CreateLocationOutput{}, err
	}

	return dto.CreateLocationOutput{
		Location: mapper.MapDomainToLocationDTO(location),
	}, nil
}
