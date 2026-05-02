package mapper

import (
	httpdto "backend/internal/adapter/in/http/dto"
	appdto "backend/internal/app/dto"

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
	return appdto.GetCatalogInput{LocationID: uuid.MustParse(req.ID)}
}

func mapOutputToCatalogItem(out appdto.CatalogItemResponse) httpdto.CatalogItemResponse {
	nutrition := mapOutputToNutrition(*out.Nutrition)
	return httpdto.CatalogItemResponse{
		ID:          out.ID.String(),
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
	return httpdto.GetCatalogResponse{Categories: out.Categories, Items: items}
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
	id, _ := uuid.Parse(req.LocationID)
	return appdto.CreateOrderInput{
		LocationID: id,
		Items:      mapRequestToOrderItems(req.Items),
	}
}

func MapOutputToCreateOrder(out appdto.CreateOrderOutput) httpdto.CreateOrderResponse {
	return httpdto.CreateOrderResponse{
		OrderID:    out.OrderID.String(),
		TotalPrice: out.TotalPrice,
		PaymentURL: out.PaymentURL,
	}
}

// --- LOCATIONS ---

func MapRequestToCreateLocation(req httpdto.CreateLocationRequest) appdto.CreateLocationInput {
	return appdto.CreateLocationInput{Slug: req.Slug, Name: req.Name, Address: req.Address}
}

func MapRequestToUpdateLocation(req httpdto.UpdateLocationRequest) appdto.UpdateLocationInput {
	id, _ := uuid.Parse(req.ID)
	return appdto.UpdateLocationInput{
		ID:       id,
		Slug:     req.Slug,
		Name:     req.Name,
		Address:  req.Address,
		IsActive: req.IsActive,
	}
}

func MapRequestToDeleteLocation(req httpdto.DeleteLocationRequest) appdto.DeleteLocationInput {
	id, _ := uuid.Parse(req.ID)
	return appdto.DeleteLocationInput{ID: id}
}

func MapRequestToGetQRCode(req httpdto.GetQRCodeRequest) appdto.GetQRCodeInput {
	id, _ := uuid.Parse(req.ID)
	return appdto.GetQRCodeInput{LocationID: id}
}

func mapOutputToLocation(out appdto.LocationResponse) httpdto.Location {
	return httpdto.Location{
		ID:        out.ID.String(),
		Slug:      out.Slug,
		Name:      out.Name,
		Address:   out.Address,
		IsActive:  out.IsActive,
		CreatedAt: out.CreatedAt.String(),
	}
}

func MapOutputToCreateLocation(out appdto.CreateLocationOutput) httpdto.CreateLocationResponse {
	return httpdto.CreateLocationResponse{Location: mapOutputToLocation(out.Location)}
}

func MapOutputToUpdateLocation(out appdto.UpdateLocationOutput) httpdto.UpdateLocationResponse {
	return httpdto.UpdateLocationResponse{Location: mapOutputToLocation(out.Location)}
}

func MapOutputToListLocations(out appdto.ListLocationsOutput) httpdto.ListLocationsResponse {
	arr := make([]httpdto.Location, len(out.Locations))
	for i := range out.Locations {
		arr[i] = mapOutputToLocation(out.Locations[i])
	}
	return httpdto.ListLocationsResponse{Locations: arr}
}

// --- ITEMS & NUTRITION ---

func mapRequestToNutrition(req httpdto.NutritionRequest) appdto.NutritionResponse {
	return appdto.NutritionResponse{
		Calories: req.Calories,
		Proteins: req.Proteins,
		Fats:     req.Fats,
		Carbs:    req.Carbs,
	}
}

func mapOutputToNutrition(out appdto.NutritionResponse) httpdto.NutritionResponse {
	return httpdto.NutritionResponse{
		Calories: out.Calories,
		Proteins: out.Proteins,
		Fats:     out.Fats,
		Carbs:    out.Carbs,
	}
}

func MapRequestToCreateItem(req httpdto.CreateItemRequest) appdto.CreateItemInput {
	var nutrition *appdto.NutritionResponse
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

func MapRequestToUpdateItem(req httpdto.UpdateItemRequest) appdto.UpdateItemInput {
	id, _ := uuid.Parse(req.ID)

	var nutrition *appdto.NutritionResponse
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

func mapOutputToItemResponse(out appdto.ItemResponse) httpdto.ItemResponse {
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
	locID, _ := uuid.Parse(req.LocationID)
	return appdto.GetInventoryInput{LocationID: locID}
}

func mapOutputToInventoryItem(out appdto.InventoryItemResponse) httpdto.InventoryItemResponse {
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

func mapRequestToInventoryItem(req httpdto.InventoryItemRequest) appdto.InventoryItemRequest {
	itemID, _ := uuid.Parse(req.ItemID)
	return appdto.InventoryItemRequest{
		ItemID:      itemID,
		Price:       req.Price,
		IsAvailable: req.IsAvailable,
		StockAmount: req.StockAmount,
	}
}

func MapRequestToUpdateInventory(req httpdto.UpdateInventoryRequest) appdto.UpdateInventoryInput {
	locID, _ := uuid.Parse(req.LocationID)

	arr := make([]appdto.InventoryItemRequest, len(req.Inventory))
	for i := range req.Inventory {
		arr[i] = mapRequestToInventoryItem(req.Inventory[i])
	}

	return appdto.UpdateInventoryInput{
		LocationID: locID,
		Inventory:  arr,
	}
}

func MapOutputToUpdateInventory(out appdto.UpdateInventoryOutput) httpdto.UpdateInventoryResponse {
	arr := make([]httpdto.InventoryItemResponse, len(out.Inventory))
	for i := range out.Inventory {
		arr[i] = mapOutputToInventoryItem(out.Inventory[i])
	}
	return httpdto.UpdateInventoryResponse{Inventory: arr}
}
