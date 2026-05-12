package dto

type GetCatalogInput struct {
	Slug string
}

type GetCatalogOutput struct {
	Location   LocationResponse
	Categories []string
	Items      []CatalogItemResponse
}
