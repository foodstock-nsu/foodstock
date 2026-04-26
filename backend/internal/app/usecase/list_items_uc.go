package usecase

import (
	"backend/internal/app/dto"
	ucerrs "backend/internal/app/errs"
	"backend/internal/app/mapper"
	"backend/internal/domain/port"
	"context"
)

type ListItemsUC struct {
	item port.ItemRepository
}

func NewListItemsUC(item port.ItemRepository) *ListItemsUC {
	return &ListItemsUC{item: item}
}

func (uc *ListItemsUC) Execute(ctx context.Context) (dto.ListItemsOutput, error) {
	items, err := uc.item.ListAll(ctx)
	if err != nil {
		return dto.ListItemsOutput{}, ucerrs.Wrap(
			ucerrs.ErrListAllItemsDB, err,
		)
	}
	return dto.ListItemsOutput{
		Items: mapper.MapDomainToItemListDTO(items),
	}, nil
}
