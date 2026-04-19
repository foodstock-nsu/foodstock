package dto

import "github.com/google/uuid"

type GetCatalogInput struct {
	LocationID uuid.UUID
	Limit      int
	Offset     int
}

type GetCatalogOutput struct {
	Categories []ItemCategory
	Items      []CatalogItem
}
