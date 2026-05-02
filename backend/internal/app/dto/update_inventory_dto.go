package dto

import "github.com/google/uuid"

type UpdateInventoryInput struct {
	LocationID uuid.UUID
	Inventory  []InventoryItemRequest
}

type UpdateInventoryOutput struct {
	Inventory []InventoryItemResponse
}
