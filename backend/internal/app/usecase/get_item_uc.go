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

type GetItemUC struct {
	item port.ItemRepository
}

func NewGetItemUC(item port.ItemRepository) *GetItemUC {
	return &GetItemUC{item: item}
}

func (uc *GetItemUC) Execute(ctx context.Context, in dto.GetItemInput) (dto.GetItemOutput, error) {
	// Get item and validate it
	item, err := uc.item.Get(ctx, in.ID)
	if err != nil {
		if errors.Is(err, pkgerrs.ErrObjectNotFound) {
			return dto.GetItemOutput{}, ucerrs.ErrItemNotFound
		}
		return dto.GetItemOutput{}, ucerrs.Wrap(
			ucerrs.ErrGetItemDB, err,
		)
	}

	if item.IsDeleted() {
		return dto.GetItemOutput{}, ucerrs.ErrItemAlreadyDeleted
	}

	return dto.GetItemOutput{Item: mapper.MapDomainToItemDTO(item)}, nil
}
