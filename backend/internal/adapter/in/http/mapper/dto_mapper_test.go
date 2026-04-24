package mapper_test

import (
	httpdto "backend/internal/adapter/in/http/dto"
	"backend/internal/adapter/in/http/mapper"
	appdto "backend/internal/app/dto"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMapErrorToResponse(t *testing.T) {
	expected := httpdto.ErrorResponse{Error: "test error"}
	result := mapper.MapErrorToResponse("test error")
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestMapRequestToAdminAuth(t *testing.T) {
	req := httpdto.AdminAuthRequest{Login: "admin", Password: "password"}
	expected := appdto.AdminAuthInput{Login: "admin", Password: "password"}
	result := mapper.MapRequestToAdminAuth(req)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestMapOutputToAdminAuth(t *testing.T) {
	out := appdto.AdminAuthOutput{Token: "token"}
	expected := httpdto.AdminAuthResponse{Token: "token"}
	result := mapper.MapOutputToAdminAuth(out)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestMapRequestToGetCatalog(t *testing.T) {
	id := uuid.New()
	req := httpdto.GetCatalogRequest{ID: id.String()}
	expected := appdto.GetCatalogInput{LocationID: id}
	result := mapper.MapRequestToGetCatalog(req)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestMapOutputToNutrition(t *testing.T) {
	c := 100
	p := 10.5
	f := 5.5
	carbs := 20.0
	out := appdto.NutritionOutput{Calories: &c, Proteins: &p, Fats: &f, Carbs: &carbs}
	expected := httpdto.NutritionResponse{Calories: &c, Proteins: &p, Fats: &f, Carbs: &carbs}
	result := mapper.MapOutputToNutrition(out)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestMapOutputToCatalogItem(t *testing.T) {
	id := uuid.New()
	desc := "desc"
	c := 100
	p := 10.5
	f := 5.5
	carbs := 20.0
	nutritionOut := appdto.NutritionOutput{Calories: &c, Proteins: &p, Fats: &f, Carbs: &carbs}
	out := appdto.CatalogItemOutput{
		ID: id, Name: "name", Description: &desc, Category: "drinks",
		PhotoURL: "url", Nutrition: &nutritionOut, Price: 100, IsAvailable: true, StockAmount: 10,
	}
	nutritionRes := httpdto.NutritionResponse{Calories: &c, Proteins: &p, Fats: &f, Carbs: &carbs}
	expected := httpdto.CatalogItemResponse{
		ID: id.String(), Name: "name", Description: &desc, Category: "drinks",
		PhotoURL: "url", Nutrition: &nutritionRes, Price: 100, IsAvailable: true, StockAmount: 10,
	}
	result := mapper.MapOutputToCatalogItem(out)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestMapOutputToGetCatalog(t *testing.T) {
	id := uuid.New()
	desc := "desc"
	c := 100
	p := 10.5
	f := 5.5
	carbs := 20.0
	nutritionOut := appdto.NutritionOutput{Calories: &c, Proteins: &p, Fats: &f, Carbs: &carbs}
	out := appdto.GetCatalogOutput{
		Categories: []string{"drinks"},
		Items: []appdto.CatalogItemOutput{
			{
				ID: id, Name: "name", Description: &desc, Category: "drinks",
				PhotoURL: "url", Nutrition: &nutritionOut, Price: 100, IsAvailable: true, StockAmount: 10,
			},
		},
	}
	nutritionRes := httpdto.NutritionResponse{Calories: &c, Proteins: &p, Fats: &f, Carbs: &carbs}
	expected := httpdto.GetCatalogResponse{
		Categories: []string{"drinks"},
		Items: []httpdto.CatalogItemResponse{
			{
				ID: id.String(), Name: "name", Description: &desc, Category: "drinks",
				PhotoURL: "url", Nutrition: &nutritionRes, Price: 100, IsAvailable: true, StockAmount: 10,
			},
		},
	}
	result := mapper.MapOutputToGetCatalog(out)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestMapRequestToCreateLocation(t *testing.T) {
	req := httpdto.CreateLocationRequest{Slug: "slug", Name: "name", Address: "address"}
	expected := appdto.CreateLocationInput{Slug: "slug", Name: "name", Address: "address"}
	result := mapper.MapRequestToCreateLocation(req)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestMapOutputToLocation(t *testing.T) {
	id := uuid.New()
	now := time.Now()
	out := appdto.LocationOutput{
		ID: id, Slug: "slug", Name: "name", Address: "address", IsActive: true, CreatedAt: now,
	}
	expected := httpdto.Location{
		ID: id.String(), Slug: "slug", Name: "name", Address: "address", IsActive: true, CreatedAt: now.String(),
	}
	result := mapper.MapOutputToLocation(out)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestMapOutputToCreateLocation(t *testing.T) {
	id := uuid.New()
	now := time.Now()
	out := appdto.CreateLocationOutput{
		Location: appdto.LocationOutput{
			ID: id, Slug: "slug", Name: "name", Address: "address", IsActive: true, CreatedAt: now,
		},
	}
	expected := httpdto.CreateLocationResponse{
		Location: httpdto.Location{
			ID: id.String(), Slug: "slug", Name: "name", Address: "address", IsActive: true, CreatedAt: now.String(),
		},
	}
	result := mapper.MapOutputToCreateLocation(out)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestMapRequestToUpdateLocation(t *testing.T) {
	id := uuid.New()
	slug := "slug"
	name := "name"
	address := "address"
	isActive := true
	req := httpdto.UpdateLocationRequest{
		ID: id.String(), Slug: &slug, Name: &name, Address: &address, IsActive: &isActive,
	}
	expected := appdto.UpdateLocationInput{
		ID: id, Slug: &slug, Name: &name, Address: &address, IsActive: &isActive,
	}
	result := mapper.MapRequestToUpdateLocation(req)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestMapOutputToUpdateLocation(t *testing.T) {
	id := uuid.New()
	now := time.Now()
	out := appdto.UpdateLocationOutput{
		Location: appdto.LocationOutput{
			ID: id, Slug: "slug", Name: "name", Address: "address", IsActive: true, CreatedAt: now,
		},
	}
	expected := httpdto.UpdateLocationResponse{
		Location: httpdto.Location{
			ID: id.String(), Slug: "slug", Name: "name", Address: "address", IsActive: true, CreatedAt: now.String(),
		},
	}
	result := mapper.MapOutputToUpdateLocation(out)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestMapRequestToDeleteLocation(t *testing.T) {
	id := uuid.New()
	req := httpdto.DeleteLocationRequest{ID: id.String()}
	expected := appdto.DeleteLocationInput{ID: id}
	result := mapper.MapRequestToDeleteLocation(req)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestMapOutputToListLocations(t *testing.T) {
	id := uuid.New()
	now := time.Now()
	out := appdto.ListLocationsOutput{
		Locations: []appdto.LocationOutput{
			{ID: id, Slug: "slug", Name: "name", Address: "address", IsActive: true, CreatedAt: now},
		},
	}
	expected := httpdto.ListLocationsResponse{
		Locations: []httpdto.Location{
			{ID: id.String(), Slug: "slug", Name: "name", Address: "address", IsActive: true, CreatedAt: now.String()},
		},
	}
	result := mapper.MapOutputToListLocations(out)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestMapRequestToGetQRCode(t *testing.T) {
	id := uuid.New()
	req := httpdto.GetQRCodeRequest{ID: id.String()}
	expected := appdto.GetQRCodeInput{LocationID: id}
	result := mapper.MapRequestToGetQRCode(req)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}
