package dto

import (
	"time"

	"github.com/google/uuid"
)

type LocationOutput struct {
	ID        uuid.UUID
	Slug      string
	Name      string
	Address   string
	IsActive  bool
	CreatedAt time.Time
}

type NutritionOutput struct {
	Calories *int
	Proteins *float64
	Fats     *float64
	Carbs    *float64
}

type ItemOutput struct {
	ID          uuid.UUID
	Name        string
	Description *string
	Category    string
	PhotoURL    string
	Nutrition   *NutritionOutput
	CreatedAt   time.Time
}

type CatalogItemOutput struct {
	ID          uuid.UUID // here is a location item id, not an item id
	Name        string
	Description *string
	Category    string
	PhotoURL    string
	Nutrition   *NutritionOutput
	Price       int64
	IsAvailable bool
	StockAmount int
}
