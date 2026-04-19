package dto

import (
	"time"

	"github.com/google/uuid"
)

type Location struct {
	ID        uuid.UUID
	Slug      string
	Name      string
	Address   string
	IsActive  bool
	CreatedAt time.Time
}

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
