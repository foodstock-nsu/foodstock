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
	DeletedAt *time.Time
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
	DeletedAt   *time.Time
}

type CatalogItemOutput struct {
	ItemID      uuid.UUID
	Name        string
	Description *string
	Category    string
	PhotoURL    string
	Nutrition   *NutritionOutput
	Price       int64
	IsAvailable bool
	StockAmount int
}

type InventoryItemInput struct {
	ItemID      uuid.UUID
	Price       *int64
	IsAvailable *bool
	StockAmount *int
}

type InventoryItemOutput struct {
	ItemID      uuid.UUID
	Price       int64
	IsAvailable bool
	StockAmount int
}
