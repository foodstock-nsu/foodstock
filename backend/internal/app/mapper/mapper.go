package mapper

import (
	"backend/internal/app/dto"
	"backend/internal/domain/model"
)

func MapDomainToLocationDTO(location *model.Location) dto.LocationOutput {
	return dto.LocationOutput{
		ID:        location.ID(),
		Slug:      location.Slug(),
		Name:      location.Name(),
		Address:   location.Address(),
		IsActive:  location.IsActive(),
		CreatedAt: location.CreatedAt(),
		DeletedAt: location.DeletedAt(),
	}
}

func MapDomainToLocationListDTO(locations []*model.Location) []dto.LocationOutput {
	res := make([]dto.LocationOutput, len(locations))
	for i := range res {
		res[i] = MapDomainToLocationDTO(locations[i])
	}
	return res
}

func mapDomainToItemNutritionDTO(nutrition model.Nutrition) dto.NutritionOutput {
	return dto.NutritionOutput{
		Calories: nutrition.Calories(),
		Proteins: nutrition.Proteins(),
		Fats:     nutrition.Fats(),
		Carbs:    nutrition.Carbs(),
	}
}

func MapDomainToItemDTO(item *model.Item) dto.ItemOutput {
	var nutrition dto.NutritionOutput
	if item.Nutrition() != nil {
		nutrition = mapDomainToItemNutritionDTO(*item.Nutrition())
	}

	return dto.ItemOutput{
		ID:          item.ID(),
		Name:        item.Name(),
		Description: item.Description(),
		Category:    item.Category().String(),
		PhotoURL:    item.PhotoURL(),
		Nutrition:   &nutrition,
		CreatedAt:   item.CreatedAt(),
		DeletedAt:   item.DeletedAt(),
	}
}

func MapDomainToItemListDTO(items []*model.Item) []dto.ItemOutput {
	res := make([]dto.ItemOutput, len(items))
	for i := range res {
		res[i] = MapDomainToItemDTO(items[i])
	}
	return res
}

func MapDomainToCatalogItemDTO(locItem *model.LocationItem, item *model.Item) dto.CatalogItemOutput {
	var nutrition dto.NutritionOutput
	if item.Nutrition() != nil {
		nutrition = mapDomainToItemNutritionDTO(*item.Nutrition())
	}

	return dto.CatalogItemOutput{
		ItemID:      item.ID(),
		Name:        item.Name(),
		Description: item.Description(),
		Category:    item.Category().String(),
		PhotoURL:    item.PhotoURL(),
		Nutrition:   &nutrition,
		Price:       locItem.Price(),
		IsAvailable: locItem.IsAvailable(),
		StockAmount: locItem.StockAmount(),
	}
}

func mapDomainToInventoryItemDTO(item *model.LocationItem) dto.InventoryItemOutput {
	return dto.InventoryItemOutput{
		ItemID:      item.ItemID(),
		Price:       item.Price(),
		IsAvailable: item.IsAvailable(),
		StockAmount: item.StockAmount(),
	}
}

func MapDomainToInventoryItemListDTO(items []*model.LocationItem) []dto.InventoryItemOutput {
	res := make([]dto.InventoryItemOutput, len(items))
	for i := range res {
		res[i] = mapDomainToInventoryItemDTO(items[i])
	}
	return res
}
