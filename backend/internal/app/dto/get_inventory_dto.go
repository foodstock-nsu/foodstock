package dto

import "github.com/google/uuid"

type GetInventoryInput struct {
	LocationID uuid.UUID
}

type GetInventoryOutput struct {
	Inventory []InventoryItemResponse
}
