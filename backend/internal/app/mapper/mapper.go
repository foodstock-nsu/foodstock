package mapper

import (
	"backend/internal/app/dto"
	"backend/internal/domain/model"
)

func MapDomainToLocationDTO(location *model.Location) dto.LocationDTO {
	return dto.LocationDTO{
		ID:        location.ID(),
		Slug:      location.Slug(),
		Name:      location.Name(),
		Address:   location.Address(),
		IsActive:  location.IsActive(),
		CreatedAt: location.CreatedAt(),
	}
}

func MapDomainToLocationListDTO(locations []*model.Location) []dto.LocationDTO {
	res := make([]dto.LocationDTO, len(locations))
	for i := range res {
		res[i] = MapDomainToLocationDTO(locations[i])
	}
	return res
}

func mapDomainToItemNutritionDTO(nutrition model.Nutrition) dto.NutritionDTO {
	return dto.NutritionDTO{
		Calories: nutrition.Calories(),
		Proteins: nutrition.Proteins(),
		Fats:     nutrition.Fats(),
		Carbs:    nutrition.Carbs(),
	}
}

func MapDomainToItemDTO(item *model.Item) dto.ItemDTO {
	var nutrition dto.NutritionDTO
	if item.Nutrition() != nil {
		nutrition = mapDomainToItemNutritionDTO(*item.Nutrition())
	}

	return dto.ItemDTO{
		ID:          item.ID(),
		Name:        item.Name(),
		Description: item.Description(),
		Category:    item.Category().String(),
		PhotoURL:    item.PhotoURL(),
		Nutrition:   &nutrition,
		CreatedAt:   item.CreatedAt(),
	}
}

func MapDomainToItemListDTO(items []*model.Item) []dto.ItemDTO {
	res := make([]dto.ItemDTO, len(items))
	for i := range res {
		res[i] = MapDomainToItemDTO(items[i])
	}
	return res
}

func MapDomainToCatalogItemDTO(locItem *model.LocationItem, item *model.Item) dto.CatalogItemDTO {
	var nutrition dto.NutritionDTO
	if item.Nutrition() != nil {
		nutrition = mapDomainToItemNutritionDTO(*item.Nutrition())
	}

	return dto.CatalogItemDTO{
		ID:          locItem.ID(),
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

func mapDomainToInventoryItemDTO(item *model.LocationItem) dto.InventoryItemDTO {
	return dto.InventoryItemDTO{
		ItemID:      item.ItemID(),
		Price:       item.Price(),
		IsAvailable: item.IsAvailable(),
		StockAmount: item.StockAmount(),
	}
}

func MapDomainToInventoryItemListDTO(items []*model.LocationItem) []dto.InventoryItemDTO {
	res := make([]dto.InventoryItemDTO, len(items))
	for i := range res {
		res[i] = mapDomainToInventoryItemDTO(items[i])
	}
	return res
}
