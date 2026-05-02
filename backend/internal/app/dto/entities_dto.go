package dto

import (
	"time"

	"github.com/google/uuid"
)

type LocationResponse struct {
	ID        uuid.UUID
	Slug      string
	Name      string
	Address   string
	IsActive  bool
	CreatedAt time.Time
}

type NutritionResponse struct {
	Calories *int
	Proteins *float64
	Fats     *float64
	Carbs    *float64
}

type ItemResponse struct {
	ID          uuid.UUID
	Name        string
	Description *string
	Category    string
	PhotoURL    string
	Nutrition   *NutritionResponse
	CreatedAt   time.Time
}

type CatalogItemResponse struct {
	ID          uuid.UUID // here is a location item id, not an item id
	Name        string
	Description *string
	Category    string
	PhotoURL    string
	Nutrition   *NutritionResponse
	Price       int64
	IsAvailable bool
	StockAmount int
}

type InventoryItemRequest struct {
	ItemID      uuid.UUID
	Price       *int64
	IsAvailable *bool
	StockAmount *int
}

type InventoryItemResponse struct {
	ItemID      uuid.UUID
	Price       int64
	IsAvailable bool
	StockAmount int
}
