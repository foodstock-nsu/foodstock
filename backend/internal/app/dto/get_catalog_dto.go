package dto

type GetCatalogInput struct {
	Slug string
}

type GetCatalogOutput struct {
	Location   LocationOutput
	Categories []string
	Items      []CatalogItemOutput
}
