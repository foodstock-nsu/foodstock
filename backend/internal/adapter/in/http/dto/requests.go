package dto

type AdminAuthRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type GetCatalogRequest struct {
	Slug string `param:"slug"`
}

type CreateLocationRequest struct {
	Slug    string `json:"slug"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

type GetLocationRequest struct {
	Slug string `param:"slug"`
}

type UpdateLocationRequest struct {
	Slug     string  `param:"slug"`
	Name     *string `json:"name"`
	Address  *string `json:"address"`
	IsActive *bool   `json:"is_active"`
}

type DeleteLocationRequest struct {
	Slug string `param:"slug"`
}

type GetQRCodeRequest struct {
	Slug string `param:"slug"`
}

type NutritionRequest struct {
	Calories *int     `json:"calories"`
	Proteins *float64 `json:"proteins"`
	Fats     *float64 `json:"fats"`
	Carbs    *float64 `json:"carbs"`
}

type CreateItemRequest struct {
	Name        string            `json:"name"`
	Description *string           `json:"description"`
	Category    string            `json:"category"`
	PhotoURL    string            `json:"photo_url"`
	Nutrition   *NutritionRequest `json:"nutrition"`
}

type GetItemRequest struct {
	ID string `param:"id"`
}

type UpdateItemRequest struct {
	ID          string            `param:"id"`
	Name        *string           `json:"name"`
	Description *string           `json:"description"`
	Category    *string           `json:"category"`
	PhotoURL    *string           `json:"photo_url"`
	Nutrition   *NutritionRequest `json:"nutrition"`
}

type DeleteItemRequest struct {
	ID string `param:"id"`
}

type OrderItemRequest struct {
	ItemID string `json:"item_id"`
	Amount int    `json:"amount"`
	Price  int64  `json:"price"`
}

type CreateOrderRequest struct {
	Slug  string             `json:"slug"`
	Items []OrderItemRequest `json:"items"`
}

type GetOrderStatusRequest struct {
	OrderID string `param:"id"`
}

type GetInventoryRequest struct {
	Slug string `param:"slug"`
}

type InventoryItemRequest struct {
	ItemID      string `json:"item_id"`
	Price       *int64 `json:"price"`
	IsAvailable *bool  `json:"is_available"`
	StockAmount *int   `json:"stock_amount"`
}

type UpdateInventoryRequest struct {
	Slug      string                 `param:"slug"`
	Inventory []InventoryItemRequest `json:"inventory"`
}
