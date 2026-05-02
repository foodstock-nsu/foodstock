package dto

import (
	"time"

	"github.com/google/uuid"
)

type LocationDTO struct {
	ID        uuid.UUID
	Slug      string
	Name      string
	Address   string
	IsActive  bool
	CreatedAt time.Time
}

type NutritionDTO struct {
	Calories *int
	Proteins *float64
	Fats     *float64
	Carbs    *float64
}

type ItemDTO struct {
	ID          uuid.UUID
	Name        string
	Description *string
	Category    string
	PhotoURL    string
	Nutrition   *NutritionDTO
	CreatedAt   time.Time
}

type CatalogItemDTO struct {
	ID          uuid.UUID // here is a location item id, not an item id
	Name        string
	Description *string
	Category    string
	PhotoURL    string
	Nutrition   *NutritionDTO
	Price       int64
	IsAvailable bool
	StockAmount int
}

type InventoryItemDTO struct {
	ItemID      uuid.UUID
	Price       int64
	IsAvailable bool
	StockAmount int
}
