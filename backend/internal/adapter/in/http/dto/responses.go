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

type ItemResponse struct {
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	Description *string            `json:"description"`
	Category    string             `json:"category"`
	PhotoURL    string             `json:"photo_url"`
	Nutrition   *NutritionResponse `json:"nutrition"`
	CreatedAt   string             `json:"created_at"`
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
	Item ItemResponse `json:"item"`
}

type UpdateItemResponse struct {
	Item ItemResponse `json:"item"`
}

type ListItemsResponse struct {
	Items []ItemResponse `json:"items"`
}

type CreateOrderResponse struct {
	OrderID    string `json:"order_id"`
	TotalPrice int64  `json:"total_price"`
	PaymentURL string `json:"payment_url"`
}

type InventoryItemResponse struct {
	ItemID      string `json:"item_id"`
	Price       int64  `json:"price"`
	IsAvailable bool   `json:"is_available"`
	StockAmount int    `json:"stock_amount"`
}

type GetInventoryResponse struct {
	Inventory []InventoryItemResponse `json:"inventory"`
}

type UpdateInventoryResponse struct {
	Inventory []InventoryItemResponse `json:"inventory"`
}
