package dto

type ErrorResponse struct {
	Error string `json:"error"`
}

type AdminAuthResponse struct {
	Token string `json:"token"`
}

type Location struct {
	ID        string `json:"id"`
	Slug      string `json:"slug"`
	Name      string `json:"name"`
	Address   string `json:"address"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
}

type NutritionResponse struct {
	Calories *int     `json:"calories"`
	Proteins *float64 `json:"proteins"`
	Fats     *float64 `json:"fats"`
	Carbs    *float64 `json:"carbs"`
}

type CatalogItemResponse struct {
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	Description *string            `json:"description"`
	Category    string             `json:"category"`
	PhotoURL    string             `json:"photo_url"`
	Nutrition   *NutritionResponse `json:"nutrition"`
	Price       int64              `json:"price"`
	IsAvailable bool               `json:"is_available"`
	StockAmount int                `json:"stock_amount"`
}

type GetCatalogResponse struct {
	Categories []string              `json:"categories"`
	Items      []CatalogItemResponse `json:"items"`
}

type CreateLocationResponse struct {
	Location Location `json:"location"`
}

type UpdateLocationResponse struct {
	Location Location `json:"location"`
}

type ListLocationsResponse struct {
	Locations []Location `json:"locations"`
}

type CreateItemResponse struct {
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	Description *string            `json:"description"`
	Category    string             `json:"category"`
	PhotoURL    string             `json:"photo_url"`
	Nutrition   *NutritionResponse `json:"nutrition"`
	CreatedAt   string             `json:"created_at"`
}
