package dto

import "github.com/google/uuid"

type GetCatalogInput struct {
	LocationID uuid.UUID
}

type GetCatalogOutput struct {
	Categories []string
	Items      []CatalogItemResponse
}
