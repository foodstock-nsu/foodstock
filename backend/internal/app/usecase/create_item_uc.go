package usecase

import (
	"backend/internal/app/dto"
	ucerrs "backend/internal/app/errs"
	"backend/internal/app/mapper"
	"backend/internal/domain/model"
	"backend/internal/domain/port"
	"context"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
)

type CreateItemUC struct {
	trManager    trm.Manager
	location     port.LocationRepository
	item         port.ItemRepository
	locationItem port.LocationItemRepository
}

func NewCreateItemUC(
	trManager trm.Manager,
	location port.LocationRepository,
	item port.ItemRepository,
	locationItem port.LocationItemRepository,
) *CreateItemUC {
	return &CreateItemUC{
		trManager:    trManager,
		location:     location,
		item:         item,
		locationItem: locationItem,
	}
}

func (uc *CreateItemUC) Execute(ctx context.Context, in dto.CreateItemInput) (dto.CreateItemOutput, error) {
	// Create rich model for item with validation
	var (
		nutrition *model.Nutrition
		err       error
	)
	if in.Nutrition != nil {
		nutrition, err = model.NewNutrition(
			in.Nutrition.Calories,
			in.Nutrition.Proteins,
			in.Nutrition.Fats,
			in.Nutrition.Carbs,
		)
		if err != nil {
			return dto.CreateItemOutput{}, ucerrs.Wrap(
				ucerrs.ErrInvalidInput, err,
			)
		}
	} else {
		nutrition = nil
	}

	item, err := model.NewItem(
		in.Name,
		in.Description,
		in.Category,
		in.PhotoURL,
		nutrition,
	)
	if err != nil {
		return dto.CreateItemOutput{}, ucerrs.Wrap(
			ucerrs.ErrInvalidInput, err,
		)
	}

	err = uc.trManager.Do(ctx, func(ctx context.Context) error {
		// Save the item into database
		err = uc.item.Create(ctx, item)
		if err != nil {
			return ucerrs.Wrap(ucerrs.ErrCreateItemDB, err)
		}

		// Get all locations
		locations, listLocsErr := uc.location.List(ctx)
		if listLocsErr != nil {
			return ucerrs.Wrap(ucerrs.ErrListLocationsDB, err)
		}

		// Create location item for each location and save it into database
		for _, location := range locations {
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

			createErr := uc.locationItem.Create(ctx, locationItem)
			if createErr != nil {
				return ucerrs.Wrap(ucerrs.ErrCreateLocationItemDB, err)
			}
		}

		return nil
	})

	if err != nil {
		return dto.CreateItemOutput{}, err
	}

	return dto.CreateItemOutput{
		Item: mapper.MapDomainToItemDTO(item),
	}, nil
}
