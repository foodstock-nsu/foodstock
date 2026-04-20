package mapper

import (
	"backend/internal/app/dto"
	"backend/internal/domain/model"
)

func MapDomainToLocationDTO(location *model.Location) dto.Location {
	return dto.Location{
		ID:        location.ID(),
		Slug:      location.Slug(),
		Name:      location.Name(),
		Address:   location.Address(),
		IsActive:  location.IsActive(),
		CreatedAt: location.CreatedAt(),
	}
}

func MapDomainToLocationListDTO(locations []*model.Location) []dto.Location {
	res := make([]dto.Location, len(locations))
	for i := range res {
		res[i] = MapDomainToLocationDTO(locations[i])
	}
	return res
}

func mapDomainToItemNutritionDTO(nutrition *model.Nutrition) *dto.ItemNutrition {
	if nutrition == nil {
		return nil
	}

	return &dto.ItemNutrition{
		Calories: nutrition.Calories(),
		Proteins: nutrition.Proteins(),
		Fats:     nutrition.Fats(),
		Carbs:    nutrition.Carbs(),
	}
}

func MapDomainToItemDTO(item *model.Item) dto.Item {
	return dto.Item{
		ID:          item.ID(),
		Name:        item.Name(),
		Description: item.Description(),
		Category:    item.Category().String(),
		PhotoURL:    item.PhotoURL(),
		Nutrition:   mapDomainToItemNutritionDTO(item.Nutrition()),
	}
}

func MapDomainToCatalogItemDTO(locItem *model.LocationItem, item *model.Item) dto.CatalogItem {
	return dto.CatalogItem{
		ID:          locItem.ID(),
		Name:        item.Name(),
		Description: item.Description(),
		Category:    item.Category().String(),
		PhotoURL:    item.PhotoURL(),
		Nutrition:   mapDomainToItemNutritionDTO(item.Nutrition()),
		Price:       locItem.Price(),
		IsAvailable: locItem.IsAvailable(),
		StockAmount: locItem.StockAmount(),
	}
}
