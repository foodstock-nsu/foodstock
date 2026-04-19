package usecase

import (
	"backend/internal/app/dto"
	ucerrs "backend/internal/app/errs"
	"backend/internal/app/mapper"
	"backend/internal/domain/model"
	"backend/internal/domain/port"
	"context"

	"github.com/google/uuid"
)

type GetCatalogUC struct {
	item         port.ItemRepository
	locationItem port.LocationItemRepository
}

func NewGetCatalogUC(
	item port.ItemRepository,
	locationItem port.LocationItemRepository,
) *GetCatalogUC {
	return &GetCatalogUC{
		item:         item,
		locationItem: locationItem,
	}
}

func (uc *GetCatalogUC) Execute(ctx context.Context, in dto.GetCatalogInput) (dto.GetCatalogOutput, error) {
	// Get an inventory for specified location
	inventory, err := uc.locationItem.List(
		ctx, in.LocationID,
		in.Limit, in.Offset,
	)
	if err != nil {
		return dto.GetCatalogOutput{}, ucerrs.Wrap(
			ucerrs.ErrListLocationItemsDB, err,
		)
	}

	if len(inventory) == 0 {
		return dto.GetCatalogOutput{}, nil
	}

	// Get a list of all items
	allItems, err := uc.item.List(ctx)
	if err != nil {
		return dto.GetCatalogOutput{}, ucerrs.Wrap(
			ucerrs.ErrListItemsDB, err,
		)
	}

	allItemsMap := make(map[uuid.UUID]*model.Item, len(allItems))
	for _, item := range allItems {
		allItemsMap[item.ID()] = item
	}

	// ========== Make the catalog ==========

	categories := make([]dto.ItemCategory, 0)
	items := make([]dto.CatalogItem, len(inventory))

	for i := range inventory {
		item, ok := allItemsMap[inventory[i].ItemID()]
		if !ok {
			continue
		}

		// Add the inventory item into output
		items[i] = mapper.MapDomainToCatalogItemDTO(
			inventory[i], item,
		)

		// Add the category of the item if it's not there yet
		var found bool
		for _, category := range categories {
			if item.Category().String() == string(category) {
				found = true
				break
			}
		}
		if !found {
			categories = append(categories, dto.ItemCategory(item.Category()))
		}
	}

	return dto.GetCatalogOutput{
		Categories: categories,
		Items:      items,
	}, nil
}
