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

type NutritionRequest struct {
	Calories *int     `json:"calories"`
	Proteins *float64 `json:"proteins"`
	Fats     *float64 `json:"fats"`
	Carbs    *float64 `json:"carbs"`
}

type CreateItemRequest struct {
	Name        string
	Description *string
	Category    string
	PhotoURL    string
	Nutrition   *NutritionRequest
}

type DeleteItemRequest struct {
	ID string `param:"id"`
}
