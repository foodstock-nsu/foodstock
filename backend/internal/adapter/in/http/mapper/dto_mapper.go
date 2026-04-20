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

func MapAdminAuthToResponse(out appdto.AdminAuthOutput) httpdto.AdminAuthResponse {
	return httpdto.AdminAuthResponse{Token: out.Token}
}

func MapRequestToGetCatalog(req httpdto.GetCatalogRequest) appdto.GetCatalogInput {
	return appdto.GetCatalogInput{
		LocationID: uuid.MustParse(req.LocationID),
	}
}

func MapResponseToNutrition(out appdto.NutritionOutput) httpdto.NutritionResponse {
	return httpdto.NutritionResponse{
		Calories: out.Calories,
		Proteins: out.Proteins,
		Fats:     out.Fats,
		Carbs:    out.Carbs,
	}
}

func MapResponseToCatalogItem(out appdto.CatalogItemOutput) httpdto.CatalogItemResponse {
	nutrition := MapResponseToNutrition(*out.Nutrition)
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

func MapResponseToGetCatalog(out appdto.GetCatalogOutput) httpdto.GetCatalogResponse {
	items := make([]httpdto.CatalogItemResponse, len(out.Items))
	for i := range items {
		items[i] = MapResponseToCatalogItem(out.Items[i])
	}
	return httpdto.GetCatalogResponse{
		Categories: out.Categories,
		Items:      items,
	}
}
