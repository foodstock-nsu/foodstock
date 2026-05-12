package dto

type UpdateInventoryInput struct {
	Slug      string
	Inventory []InventoryItemRequest
}

type UpdateInventoryOutput struct {
	Inventory []InventoryItemResponse
}
