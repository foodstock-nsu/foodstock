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
	}
}

func MapDomainToLocationListDTO(locations []*model.Location) []dto.LocationOutput {
	res := make([]dto.LocationOutput, len(locations))
	for i := range res {
		res[i] = MapDomainToLocationDTO(locations[i])
	}
	return res
}

func mapDomainToItemNutritionDTO(nutrition *model.Nutrition) *dto.NutritionOutput {
	if nutrition == nil {
		return nil
	}

	return &dto.NutritionOutput{
		Calories: nutrition.Calories(),
		Proteins: nutrition.Proteins(),
		Fats:     nutrition.Fats(),
		Carbs:    nutrition.Carbs(),
	}
}

func MapDomainToItemDTO(item *model.Item) dto.ItemOutput {
	return dto.ItemOutput{
		ID:          item.ID(),
		Name:        item.Name(),
		Description: item.Description(),
		Category:    item.Category().String(),
		PhotoURL:    item.PhotoURL(),
		Nutrition:   mapDomainToItemNutritionDTO(item.Nutrition()),
		CreatedAt:   item.CreatedAt(),
	}
}

func MapDomainToCatalogItemDTO(locItem *model.LocationItem, item *model.Item) dto.CatalogItemOutput {
	return dto.CatalogItemOutput{
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
