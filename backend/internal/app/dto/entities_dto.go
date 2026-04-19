package dto

import "github.com/google/uuid"

type ItemCategory string

type ItemNutrition struct {
	Calories *int
	Proteins *float64
	Fats     *float64
	Carbs    *float64
}

type Item struct {
	ID          uuid.UUID
	Name        string
	Description *string
	Category    ItemCategory
	PhotoURL    string
	Nutrition   *ItemNutrition
}

type CatalogItem struct {
	ID          uuid.UUID // here is a location item id, not an item id
	Name        string
	Description *string
	Category    ItemCategory
	PhotoURL    string
	Nutrition   *ItemNutrition
	Price       int64
	IsAvailable bool
	StockAmount int
}
