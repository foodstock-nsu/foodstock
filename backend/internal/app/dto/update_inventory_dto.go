package dto

type UpdateInventoryInput struct {
	Slug      string
	Inventory []InventoryItemInput
}

type UpdateInventoryOutput struct {
	Inventory []InventoryItemOutput
}
