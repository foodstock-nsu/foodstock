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

type UpdateItemUC struct {
	item port.ItemRepository
}

func NewUpdateItemUC(item port.ItemRepository) *UpdateItemUC {
	return &UpdateItemUC{item: item}
}

func (uc *UpdateItemUC) Execute(ctx context.Context, in dto.UpdateItemInput) (dto.UpdateItemOutput, error) {
	// Get item
	item, err := uc.item.Get(ctx, in.ID)
	if err != nil {
		if errors.Is(err, pkgerrs.ErrObjectNotFound) {
			return dto.UpdateItemOutput{}, ucerrs.ErrItemNotFound
		}
		return dto.UpdateItemOutput{}, ucerrs.Wrap(
			ucerrs.ErrGetItemDB, err,
		)
	}

	// Update nutrition
	nutrition := item.Nutrition()
	err = nutrition.Update(
		in.Nutrition.Calories,
		in.Nutrition.Proteins,
		in.Nutrition.Fats,
		in.Nutrition.Carbs,
	)
	if err != nil {
		return dto.UpdateItemOutput{}, ucerrs.Wrap(
			ucerrs.ErrInvalidInput, err,
		)
	}

	// Update item
	err = item.Update(
		in.Name,
		in.Description,
		in.Category,
		in.PhotoURL,
		nutrition,
	)
	if err != nil {
		return dto.UpdateItemOutput{}, ucerrs.Wrap(
			ucerrs.ErrInvalidInput, err,
		)
	}

	err = uc.item.Update(ctx, item)
	if err != nil {
		return dto.UpdateItemOutput{}, ucerrs.Wrap(
			ucerrs.ErrUpdateItemDB, err,
		)
	}

	return dto.UpdateItemOutput{Item: mapper.MapDomainToItemDTO(item)}, nil
}
