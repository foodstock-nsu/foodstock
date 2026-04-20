package dto

type AdminAuthRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type GetCatalogRequest struct {
	LocationID string `query:"location_id"`
}
