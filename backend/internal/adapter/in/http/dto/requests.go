package dto

type AdminAuthRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type GetCatalogRequest struct {
	ID string `query:"id"`
}

type CreateLocationRequest struct {
	Slug    string `json:"slug"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

type UpdateLocationRequest struct {
	ID       string  `param:"id"`
	Slug     *string `json:"slug"`
	Name     *string `json:"name"`
	Address  *string `json:"address"`
	IsActive *bool   `json:"is_active"`
}

type DeleteLocationRequest struct {
	ID string `param:"id"`
}

type GetQRCodeRequest struct {
	ID string `param:"id"`
}
