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

// --- COMMON & ERRORS ---

func TestMapErrorToResponse(t *testing.T) {
	expected := httpdto.ErrorResponse{Error: "test error"}
	result := mapper.MapErrorToResponse("test error")
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

// --- ADMIN & AUTH ---

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

// --- CATALOG ---

func TestMapRequestToGetCatalog(t *testing.T) {
	id := uuid.New()
	req := httpdto.GetCatalogRequest{ID: id.String()}
	expected := appdto.GetCatalogInput{LocationID: id}
	result := mapper.MapRequestToGetCatalog(req)
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

// --- ORDERS ---

func TestMapRequestToCreateOrder(t *testing.T) {
	locID := uuid.New()
	itemID1 := uuid.New()
	itemID2 := uuid.New()

	req := httpdto.CreateOrderRequest{
		LocationID: locID.String(),
		Items: []httpdto.OrderItemRequest{
			{itemID1.String(), 1, 200},
			{itemID2.String(), 2, 300},
		},
	}

	expected := appdto.CreateOrderInput{
		LocationID: locID,
		Items: []appdto.OrderItemInput{
			{itemID1, 1, 200},
			{itemID2, 2, 300},
		},
	}

	result := mapper.MapRequestToCreateOrder(req)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestMapOutputToCreateOrder(t *testing.T) {
	orderID := uuid.New()

	out := appdto.CreateOrderOutput{
		OrderID:    orderID,
		TotalPrice: 2000,
		PaymentURL: "https://pay.easy/1",
	}

	expected := httpdto.CreateOrderResponse{
		OrderID:    orderID.String(),
		TotalPrice: 2000,
		PaymentURL: "https://pay.easy/1",
	}

	result := mapper.MapOutputToCreateOrder(out)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

// --- LOCATIONS ---

func TestMapRequestToCreateLocation(t *testing.T) {
	req := httpdto.CreateLocationRequest{Slug: "slug", Name: "name", Address: "address"}
	expected := appdto.CreateLocationInput{Slug: "slug", Name: "name", Address: "address"}
	result := mapper.MapRequestToCreateLocation(req)
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

// --- ITEMS & NUTRITION ---

func TestMapRequestToCreateItem(t *testing.T) {
	desc := "description"
	cal := 100
	req := httpdto.CreateItemRequest{
		Name:        "name",
		Description: &desc,
		Category:    "cat",
		PhotoURL:    "url",
		Nutrition:   &httpdto.NutritionRequest{Calories: &cal},
	}
	expected := appdto.CreateItemInput{
		Name:        "name",
		Description: &desc,
		Category:    "cat",
		PhotoURL:    "url",
		Nutrition:   &appdto.NutritionOutput{Calories: &cal},
	}
	result := mapper.MapRequestToCreateItem(req)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestMapRequestToUpdateItem(t *testing.T) {
	id := uuid.New()
	name := "new name"
	cal := 200
	req := httpdto.UpdateItemRequest{
		ID:        id.String(),
		Name:      &name,
		Nutrition: &httpdto.NutritionRequest{Calories: &cal},
	}
	expected := appdto.UpdateItemInput{
		ID:        id,
		Name:      &name,
		Nutrition: &appdto.NutritionOutput{Calories: &cal},
	}
	result := mapper.MapRequestToUpdateItem(req)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestMapRequestToDeleteItem(t *testing.T) {
	id := uuid.New()
	req := httpdto.DeleteItemRequest{ID: id.String()}
	expected := appdto.DeleteItemInput{ID: id}
	result := mapper.MapRequestToDeleteItem(req)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestMapOutputToCreateItem(t *testing.T) {
	id := uuid.New()
	now := time.Now()
	calories := 640
	out := appdto.CreateItemOutput{
		Item: appdto.ItemOutput{
			ID:       id,
			Name:     "name",
			Category: "cat",
			Nutrition: &appdto.NutritionOutput{
				Calories: &calories,
				Proteins: nil,
				Fats:     nil,
				Carbs:    nil,
			},
			CreatedAt: now,
		},
	}
	expected := httpdto.CreateItemResponse{
		Item: httpdto.ItemResponse{
			ID:       id.String(),
			Name:     "name",
			Category: "cat",
			Nutrition: &httpdto.NutritionResponse{
				Calories: &calories,
				Proteins: nil,
				Fats:     nil,
				Carbs:    nil,
			},
			CreatedAt: now.String(),
		},
	}
	result := mapper.MapOutputToCreateItem(out)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestMapOutputToUpdateItem(t *testing.T) {
	id := uuid.New()
	now := time.Now()
	out := appdto.UpdateItemOutput{
		Item: appdto.ItemOutput{
			ID: id, Name: "updated", Category: "cat", CreatedAt: now,
		},
	}
	expected := httpdto.UpdateItemResponse{
		Item: httpdto.ItemResponse{
			ID: id.String(), Name: "updated", Category: "cat", CreatedAt: now.String(),
		},
	}
	result := mapper.MapOutputToUpdateItem(out)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestMapOutputToListItems(t *testing.T) {
	id := uuid.New()
	now := time.Now()
	out := appdto.ListItemsOutput{
		Items: []appdto.ItemOutput{
			{ID: id, Name: "name", Category: "cat", CreatedAt: now},
		},
	}
	expected := httpdto.ListItemsResponse{
		Items: []httpdto.ItemResponse{
			{ID: id.String(), Name: "name", Category: "cat", CreatedAt: now.String()},
		},
	}
	result := mapper.MapOutputToListItems(out)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}
