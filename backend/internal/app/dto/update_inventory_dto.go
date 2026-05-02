package dto

import "github.com/google/uuid"

type UpdateInventoryInput struct {
	LocationID uuid.UUID
	Inventory  []InventoryItemDTO
}

type UpdateInventoryOutput struct {
	Inventory []InventoryItemDTO
}
