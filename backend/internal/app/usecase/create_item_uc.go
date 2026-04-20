package usecase

import (
	"backend/internal/app/dto"
	ucerrs "backend/internal/app/errs"
	"backend/internal/app/mapper"
	"backend/internal/domain/model"
	"backend/internal/domain/port"
	"context"
)

type CreateItemUC struct {
	item port.ItemRepository
}

func NewCreateItemUC(item port.ItemRepository) *CreateItemUC {
	return &CreateItemUC{item: item}
}

func (uc *CreateItemUC) Execute(ctx context.Context, in dto.CreateItemInput) (dto.CreateItemOutput, error) {
	// Create rich model with validation
	nutrition, err := model.NewNutrition(
		in.Nutrition.Calories,
		in.Nutrition.Proteins,
		in.Nutrition.Fats,
		in.Nutrition.Carbs,
	)
	if err != nil {
		return dto.CreateItemOutput{}, ucerrs.ErrInvalidInput
	}

	item, err := model.NewItem(
		in.Name,
		in.Description,
		in.Category,
		in.PhotoURL,
		nutrition,
	)
	if err != nil {
		return dto.CreateItemOutput{}, ucerrs.ErrInvalidInput
	}

	// Save into database
	err = uc.item.Create(ctx, item)
	if err != nil {
		return dto.CreateItemOutput{}, ucerrs.Wrap(
			ucerrs.ErrCreateItemDB, err,
		)
	}

	return dto.CreateItemOutput{
		Item: mapper.MapDomainToItemDTO(item),
	}, nil
}
