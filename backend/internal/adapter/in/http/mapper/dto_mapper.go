package mapper

import (
	httpdto "backend/internal/adapter/in/http/dto"
	appdto "backend/internal/app/dto"

	"github.com/google/uuid"
)

func MapErrorToResponse(errStr string) httpdto.ErrorResponse {
	return httpdto.ErrorResponse{Error: errStr}
}

func MapRequestToAdminAuth(req httpdto.AdminAuthRequest) appdto.AdminAuthInput {
	return appdto.AdminAuthInput{
		Login:    req.Login,
		Password: req.Password,
	}
}

func MapOutputToAdminAuth(out appdto.AdminAuthOutput) httpdto.AdminAuthResponse {
	return httpdto.AdminAuthResponse{Token: out.Token}
}

func MapRequestToGetCatalog(req httpdto.GetCatalogRequest) appdto.GetCatalogInput {
	return appdto.GetCatalogInput{
		LocationID: uuid.MustParse(req.ID),
	}
}

func MapOutputToNutrition(out appdto.NutritionOutput) httpdto.NutritionResponse {
	return httpdto.NutritionResponse{
		Calories: out.Calories,
		Proteins: out.Proteins,
		Fats:     out.Fats,
		Carbs:    out.Carbs,
	}
}

func MapOutputToCatalogItem(out appdto.CatalogItemOutput) httpdto.CatalogItemResponse {
	nutrition := MapOutputToNutrition(*out.Nutrition)
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
	for i := range items {
		items[i] = MapOutputToCatalogItem(out.Items[i])
	}
	return httpdto.GetCatalogResponse{
		Categories: out.Categories,
		Items:      items,
	}
}

func MapRequestToCreateLocation(req httpdto.CreateLocationRequest) appdto.CreateLocationInput {
	return appdto.CreateLocationInput{
		Slug:    req.Slug,
		Name:    req.Name,
		Address: req.Address,
	}
}

func MapOutputToLocation(out appdto.LocationOutput) httpdto.Location {
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
	return httpdto.CreateLocationResponse{
		Location: MapOutputToLocation(out.Location),
	}
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

func MapOutputToUpdateLocation(out appdto.UpdateLocationOutput) httpdto.UpdateLocationResponse {
	return httpdto.UpdateLocationResponse{
		Location: MapOutputToLocation(out.Location),
	}
}

func MapRequestToDeleteLocation(req httpdto.DeleteLocationRequest) appdto.DeleteLocationInput {
	id, _ := uuid.Parse(req.ID)
	return appdto.DeleteLocationInput{ID: id}
}

func MapOutputToListLocations(out appdto.ListLocationsOutput) httpdto.ListLocationsResponse {
	arr := make([]httpdto.Location, len(out.Locations))
	for i := range arr {
		arr[i] = MapOutputToLocation(out.Locations[i])
	}

	return httpdto.ListLocationsResponse{Locations: arr}
}

func MapRequestToGetQRCode(req httpdto.GetQRCodeRequest) appdto.GetQRCodeInput {
	id, _ := uuid.Parse(req.ID)
	return appdto.GetQRCodeInput{LocationID: id}
}

func MapRequestToNutrition(req httpdto.NutritionRequest) appdto.NutritionOutput {
	return appdto.NutritionOutput{
		Calories: req.Calories,
		Proteins: req.Proteins,
		Fats:     req.Fats,
		Carbs:    req.Carbs,
	}
}

func MapRequestToCreateItem(req httpdto.CreateItemRequest) appdto.CreateItemInput {
	nutrition := MapRequestToNutrition(*req.Nutrition)
	return appdto.CreateItemInput{
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
		PhotoURL:    req.PhotoURL,
		Nutrition:   &nutrition,
	}
}

func MapOutputToCreateItem(out appdto.CreateItemOutput) httpdto.CreateItemResponse {
	nutrition := MapOutputToNutrition(*out.Item.Nutrition)
	return httpdto.CreateItemResponse{
		ID:          out.Item.ID.String(),
		Name:        out.Item.Name,
		Description: out.Item.Description,
		Category:    out.Item.Category,
		PhotoURL:    out.Item.PhotoURL,
		Nutrition:   &nutrition,
		CreatedAt:   out.Item.CreatedAt.String(),
	}
}

func MapRequestToDeleteItem(req httpdto.DeleteItemRequest) appdto.DeleteItemInput {
	id, _ := uuid.Parse(req.ID)
	return appdto.DeleteItemInput{ID: id}
}
