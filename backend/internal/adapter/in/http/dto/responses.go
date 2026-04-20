package dto

type AdminAuthResponse struct {
	Token string `json:"token"`
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
