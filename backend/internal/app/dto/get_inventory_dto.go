package dto

type GetInventoryInput struct {
	Slug string
}

type GetInventoryOutput struct {
	Inventory []InventoryItemResponse
}
