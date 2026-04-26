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

func MapOutputToCatalogItem(out appdto.CatalogItemOutput) httpdto.CatalogItemResponse {
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
		items[i] = MapOutputToCatalogItem(out.Items[i])
	}
	return httpdto.GetCatalogResponse{Categories: out.Categories, Items: items}
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

func mapOutputToLocation(out appdto.LocationOutput) httpdto.Location {
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
