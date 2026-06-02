package mapper

import (
	httpdto "backend/internal/adapter/in/http/dto"
	appdto "backend/internal/app/dto"
	pkgutils "backend/pkg/utils"

	"github.com/google/uuid"
)

// --- COMMON & ERRORS ---

func MapErrorToResponse(errStr string) httpdto.ErrorResponse {
	return httpdto.ErrorResponse{Error: errStr}
}

// --- ADMIN & AUTH ---

func MapRequestToAdminAuth(req httpdto.AdminAuthRequest) appdto.AdminAuthInput {
	return appdto.AdminAuthInput{Login: req.Login, Password: req.Password}
}

func MapOutputToAdminAuth(out appdto.AdminAuthOutput) httpdto.AdminAuthResponse {
	return httpdto.AdminAuthResponse{Token: out.Token}
}

// --- CATALOG ---

func MapRequestToGetCatalog(req httpdto.GetCatalogRequest) appdto.GetCatalogInput {
	return appdto.GetCatalogInput{Slug: req.Slug}
}

func mapOutputToCatalogItem(out appdto.CatalogItemOutput) httpdto.CatalogItemResponse {
	nutrition := mapOutputToNutrition(*out.Nutrition)
	return httpdto.CatalogItemResponse{
		ItemID:      out.ItemID.String(),
		Name:        out.Name,
		Description: out.Description,
		Category:    out.Category,
		PhotoURL:    out.PhotoURL,
		Nutrition:   &nutrition,
		Price:       out.Price,
		IsAvailable: out.IsAvailable,
		StockAmount: out.StockAmount,
	}
}

func MapOutputToGetCatalog(out appdto.GetCatalogOutput) httpdto.GetCatalogResponse {
	items := make([]httpdto.CatalogItemResponse, len(out.Items))
	for i := range out.Items {
		items[i] = mapOutputToCatalogItem(out.Items[i])
	}
	return httpdto.GetCatalogResponse{
		Location:   mapOutputToLocation(out.Location),
		Categories: out.Categories,
		Items:      items,
	}
}

// --- ORDERS ---

func mapRequestToOrderItem(req httpdto.OrderItemRequest) appdto.OrderItemInput {
	id, _ := uuid.Parse(req.ItemID)
	return appdto.OrderItemInput{
		ItemID: id,
		Amount: req.Amount,
		Price:  req.Price,
	}
}

func mapRequestToOrderItems(rawItems []httpdto.OrderItemRequest) []appdto.OrderItemInput {
	var items []appdto.OrderItemInput
	for i := range rawItems {
		items = append(items, mapRequestToOrderItem(rawItems[i]))
	}
	return items
}

func MapRequestToCreateOrder(req httpdto.CreateOrderRequest) appdto.CreateOrderInput {
	return appdto.CreateOrderInput{
		Slug:  req.Slug,
		Items: mapRequestToOrderItems(req.Items),
	}
}

func MapOutputToCreateOrder(out appdto.CreateOrderOutput) httpdto.CreateOrderResponse {
	return httpdto.CreateOrderResponse{
		OrderID:    out.OrderID.String(),
		TotalPrice: out.TotalPrice,
		PaymentURL: out.PaymentURL,
	}
}

func MapRequestToGetOrderStatus(req httpdto.GetOrderStatusRequest) appdto.GetOrderStatusInput {
	id, _ := uuid.Parse(req.OrderID)
	return appdto.GetOrderStatusInput{OrderID: id}
}

func MapOutputToGetOrderStatus(out appdto.GetOrderStatusOutput) httpdto.GetOrderStatusResponse {
	return httpdto.GetOrderStatusResponse{Status: out.Status}
}

// --- LOCATIONS ---

func MapRequestToCreateLocation(req httpdto.CreateLocationRequest) appdto.CreateLocationInput {
	return appdto.CreateLocationInput{Slug: req.Slug, Name: req.Name, Address: req.Address}
}

func MapRequestToGetLocation(req httpdto.GetLocationRequest) appdto.GetLocationInput {
	return appdto.GetLocationInput{Slug: req.Slug}
}

func MapRequestToUpdateLocation(req httpdto.UpdateLocationRequest) appdto.UpdateLocationInput {
	return appdto.UpdateLocationInput{
		Slug:     req.Slug,
		Name:     req.Name,
		Address:  req.Address,
		IsActive: req.IsActive,
	}
}

func MapRequestToDeleteLocation(req httpdto.DeleteLocationRequest) appdto.DeleteLocationInput {
	return appdto.DeleteLocationInput{Slug: req.Slug}
}

func MapRequestToGetQRCode(req httpdto.GetQRCodeRequest) appdto.GetQRCodeInput {
	return appdto.GetQRCodeInput{Slug: req.Slug}
}

func mapOutputToLocation(out appdto.LocationOutput) httpdto.LocationResponse {
	var deletedAt *string
	if out.DeletedAt != nil {
		deletedAt = pkgutils.VPtr(out.DeletedAt.String())
	}
	return httpdto.LocationResponse{
		ID:        out.ID.String(),
		Slug:      out.Slug,
		Name:      out.Name,
		Address:   out.Address,
		IsActive:  out.IsActive,
		CreatedAt: out.CreatedAt.String(),
		DeletedAt: deletedAt,
	}
}

func MapOutputToCreateLocation(out appdto.CreateLocationOutput) httpdto.CreateLocationResponse {
	return httpdto.CreateLocationResponse{Location: mapOutputToLocation(out.Location)}
}

func MapOutputToGetLocation(out appdto.GetLocationOutput) httpdto.GetLocationResponse {
	return httpdto.GetLocationResponse{Location: mapOutputToLocation(out.Location)}
}

func MapOutputToUpdateLocation(out appdto.UpdateLocationOutput) httpdto.UpdateLocationResponse {
	return httpdto.UpdateLocationResponse{Location: mapOutputToLocation(out.Location)}
}

func MapOutputToListLocations(out appdto.ListLocationsOutput) httpdto.ListLocationsResponse {
	arr := make([]httpdto.LocationResponse, len(out.Locations))
	for i := range out.Locations {
		arr[i] = mapOutputToLocation(out.Locations[i])
	}
	return httpdto.ListLocationsResponse{Locations: arr}
}

// --- ITEMS & NUTRITION ---

func mapRequestToNutrition(req httpdto.NutritionRequest) appdto.NutritionOutput {
	return appdto.NutritionOutput{
		Calories: req.Calories,
		Proteins: req.Proteins,
		Fats:     req.Fats,
		Carbs:    req.Carbs,
	}
}

func mapOutputToNutrition(out appdto.NutritionOutput) httpdto.NutritionResponse {
	return httpdto.NutritionResponse{
		Calories: out.Calories,
		Proteins: out.Proteins,
		Fats:     out.Fats,
		Carbs:    out.Carbs,
	}
}

func MapRequestToCreateItem(req httpdto.CreateItemRequest) appdto.CreateItemInput {
	var nutrition *appdto.NutritionOutput
	if req.Nutrition != nil {
		mapped := mapRequestToNutrition(*req.Nutrition)
		nutrition = &mapped
	}
	return appdto.CreateItemInput{
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
		PhotoURL:    req.PhotoURL,
		Nutrition:   nutrition,
	}
}

func MapRequestToGetItem(req httpdto.GetItemRequest) appdto.GetItemInput {
	id, _ := uuid.Parse(req.ID)
	return appdto.GetItemInput{ID: id}
}

func MapRequestToUpdateItem(req httpdto.UpdateItemRequest) appdto.UpdateItemInput {
	id, _ := uuid.Parse(req.ID)

	var nutrition *appdto.NutritionOutput
	if req.Nutrition != nil {
		mapped := mapRequestToNutrition(*req.Nutrition)
		nutrition = &mapped
	}

	return appdto.UpdateItemInput{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
		PhotoURL:    req.PhotoURL,
		Nutrition:   nutrition,
	}
}

func MapRequestToDeleteItem(req httpdto.DeleteItemRequest) appdto.DeleteItemInput {
	id, _ := uuid.Parse(req.ID)
	return appdto.DeleteItemInput{ID: id}
}

func mapOutputToItemResponse(out appdto.ItemOutput) httpdto.ItemResponse {
	var nutrition *httpdto.NutritionResponse
	if out.Nutrition != nil {
		mapped := mapOutputToNutrition(*out.Nutrition)
		nutrition = &mapped
	}
	return httpdto.ItemResponse{
		ID:          out.ID.String(),
		Name:        out.Name,
		Description: out.Description,
		Category:    out.Category,
		PhotoURL:    out.PhotoURL,
		Nutrition:   nutrition,
		CreatedAt:   out.CreatedAt.String(),
	}
}

func MapOutputToCreateItem(out appdto.CreateItemOutput) httpdto.CreateItemResponse {
	return httpdto.CreateItemResponse{Item: mapOutputToItemResponse(out.Item)}
}

func MapOutputToGetItem(out appdto.GetItemOutput) httpdto.GetItemResponse {
	return httpdto.GetItemResponse{Item: mapOutputToItemResponse(out.Item)}
}

func MapOutputToUpdateItem(out appdto.UpdateItemOutput) httpdto.UpdateItemResponse {
	return httpdto.UpdateItemResponse{Item: mapOutputToItemResponse(out.Item)}
}

func MapOutputToListItems(out appdto.ListItemsOutput) httpdto.ListItemsResponse {
	arr := make([]httpdto.ItemResponse, len(out.Items))
	for i := range out.Items {
		arr[i] = mapOutputToItemResponse(out.Items[i])
	}
	return httpdto.ListItemsResponse{Items: arr}
}

func MapRequestToGetInventory(req httpdto.GetInventoryRequest) appdto.GetInventoryInput {
	return appdto.GetInventoryInput{Slug: req.Slug}
}

func mapOutputToInventoryItem(out appdto.InventoryItemOutput) httpdto.InventoryItemResponse {
	return httpdto.InventoryItemResponse{
		ItemID:      out.ItemID.String(),
		Price:       out.Price,
		IsAvailable: out.IsAvailable,
		StockAmount: out.StockAmount,
	}
}

func MapOutputToGetInventory(out appdto.GetInventoryOutput) httpdto.GetInventoryResponse {
	arr := make([]httpdto.InventoryItemResponse, len(out.Inventory))
	for i := range out.Inventory {
		arr[i] = mapOutputToInventoryItem(out.Inventory[i])
	}
	return httpdto.GetInventoryResponse{Inventory: arr}
}

func mapRequestToInventoryItem(req httpdto.InventoryItemRequest) appdto.InventoryItemInput {
	itemID, _ := uuid.Parse(req.ItemID)
	return appdto.InventoryItemInput{
		ItemID:      itemID,
		Price:       req.Price,
		IsAvailable: req.IsAvailable,
		StockAmount: req.StockAmount,
	}
}

func MapRequestToUpdateInventory(req httpdto.UpdateInventoryRequest) appdto.UpdateInventoryInput {
	arr := make([]appdto.InventoryItemInput, len(req.Inventory))
	for i := range req.Inventory {
		arr[i] = mapRequestToInventoryItem(req.Inventory[i])
	}

	return appdto.UpdateInventoryInput{
		Slug:      req.Slug,
		Inventory: arr,
	}
}

func MapOutputToUpdateInventory(out appdto.UpdateInventoryOutput) httpdto.UpdateInventoryResponse {
	arr := make([]httpdto.InventoryItemResponse, len(out.Inventory))
	for i := range out.Inventory {
		arr[i] = mapOutputToInventoryItem(out.Inventory[i])
	}
	return httpdto.UpdateInventoryResponse{Inventory: arr}
}
