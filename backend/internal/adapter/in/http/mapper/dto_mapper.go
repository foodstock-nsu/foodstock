package mapper

import (
	httpdto "backend/internal/adapter/in/http/dto"
	appdto "backend/internal/app/dto"

	"github.com/google/uuid"
)

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

func MapOutputToLocation(out appdto.LocationOutput) httpdto.LocationResponse {
	return httpdto.LocationResponse{
		ID:        out.ID.String(),
		Slug:      out.Slug,
		Name:      out.Name,
		Address:   out.Name,
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
	return appdto.UpdateLocationInput{
		ID:       uuid.MustParse(req.ID),
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
	return appdto.DeleteLocationInput{ID: uuid.MustParse(req.ID)}
}

func MapOutputToListLocations(out appdto.ListLocationsOutput) httpdto.ListLocationsResponse {
	arr := make([]httpdto.LocationResponse, len(out.Locations))
	for i := range arr {
		arr[i] = MapOutputToLocation(out.Locations[i])
	}

	return httpdto.ListLocationsResponse{Locations: arr}
}

func MapRequestToGetQRCode(req httpdto.GetQRCodeRequest) appdto.GetQRCodeInput {
	return appdto.GetQRCodeInput{LocationID: uuid.MustParse(req.ID)}
}
