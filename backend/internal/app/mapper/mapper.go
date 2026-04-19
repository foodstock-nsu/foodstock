package mapper

import (
	"backend/internal/app/dto"
	"backend/internal/domain/model"
)

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

func MapDomainToCatalogItemDTO(locItem *model.LocationItem, item *model.Item) dto.CatalogItem {
	return dto.CatalogItem{
		ID:          locItem.ID(),
		Name:        item.Name(),
		Description: item.Description(),
		Category:    dto.ItemCategory(item.Category()),
		PhotoURL:    item.PhotoURL(),
		Nutrition:   mapDomainToItemNutritionDTO(item.Nutrition()),
		Price:       locItem.Price(),
		IsAvailable: locItem.IsAvailable(),
		StockAmount: locItem.StockAmount(),
	}
}
