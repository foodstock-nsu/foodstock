package mapper

import (
	"backend/internal/app/dto"
	"backend/internal/domain/model"
)

func MapDomainToLocationDTO(location *model.Location) dto.LocationResponse {
	return dto.LocationResponse{
		ID:        location.ID(),
		Slug:      location.Slug(),
		Name:      location.Name(),
		Address:   location.Address(),
		IsActive:  location.IsActive(),
		CreatedAt: location.CreatedAt(),
	}
}

func MapDomainToLocationListDTO(locations []*model.Location) []dto.LocationResponse {
	res := make([]dto.LocationResponse, len(locations))
	for i := range res {
		res[i] = MapDomainToLocationDTO(locations[i])
	}
	return res
}

func mapDomainToItemNutritionDTO(nutrition model.Nutrition) dto.NutritionResponse {
	return dto.NutritionResponse{
		Calories: nutrition.Calories(),
		Proteins: nutrition.Proteins(),
		Fats:     nutrition.Fats(),
		Carbs:    nutrition.Carbs(),
	}
}

func MapDomainToItemDTO(item *model.Item) dto.ItemResponse {
	var nutrition dto.NutritionResponse
	if item.Nutrition() != nil {
		nutrition = mapDomainToItemNutritionDTO(*item.Nutrition())
	}

	return dto.ItemResponse{
		ID:          item.ID(),
		Name:        item.Name(),
		Description: item.Description(),
		Category:    item.Category().String(),
		PhotoURL:    item.PhotoURL(),
		Nutrition:   &nutrition,
		CreatedAt:   item.CreatedAt(),
	}
}

func MapDomainToItemListDTO(items []*model.Item) []dto.ItemResponse {
	res := make([]dto.ItemResponse, len(items))
	for i := range res {
		res[i] = MapDomainToItemDTO(items[i])
	}
	return res
}

func MapDomainToCatalogItemDTO(locItem *model.LocationItem, item *model.Item) dto.CatalogItemResponse {
	var nutrition dto.NutritionResponse
	if item.Nutrition() != nil {
		nutrition = mapDomainToItemNutritionDTO(*item.Nutrition())
	}

	return dto.CatalogItemResponse{
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

func mapDomainToInventoryItemDTO(item *model.LocationItem) dto.InventoryItemResponse {
	return dto.InventoryItemResponse{
		ItemID:      item.ItemID(),
		Price:       item.Price(),
		IsAvailable: item.IsAvailable(),
		StockAmount: item.StockAmount(),
	}
}

func MapDomainToInventoryItemListDTO(items []*model.LocationItem) []dto.InventoryItemResponse {
	res := make([]dto.InventoryItemResponse, len(items))
	for i := range res {
		res[i] = mapDomainToInventoryItemDTO(items[i])
	}
	return res
}
