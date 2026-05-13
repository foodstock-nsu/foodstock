package mapper_test

import (
	httpdto "backend/internal/adapter/in/http/dto"
	"backend/internal/adapter/in/http/mapper"
	appdto "backend/internal/app/dto"
	"backend/pkg/utils"
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
	req := httpdto.GetCatalogRequest{Slug: "new_1"}
	expected := appdto.GetCatalogInput{Slug: "new_1"}
	result := mapper.MapRequestToGetCatalog(req)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestMapOutputToGetCatalog(t *testing.T) {
	id := uuid.New()
	locID := uuid.New()
	now := time.Now().UTC()
	nutritionOut := appdto.NutritionOutput{
		Calories: utils.VPtr(100),
		Proteins: utils.VPtr(10.5),
		Fats:     utils.VPtr(5.5),
		Carbs:    utils.VPtr(20.0),
	}
	out := appdto.GetCatalogOutput{
		Location: appdto.LocationOutput{
			ID:        locID,
			Slug:      "test_1",
			Name:      "test name of location",
			Address:   "address of test name of location",
			IsActive:  true,
			CreatedAt: now,
			DeletedAt: nil,
		},
		Categories: []string{"drinks"},
		Items: []appdto.CatalogItemOutput{
			{
				ItemID:      id,
				Name:        "name",
				Description: utils.VPtr("desc"),
				Category:    "drinks",
				PhotoURL:    "url",
				Nutrition:   &nutritionOut,
				Price:       100,
				IsAvailable: true,
				StockAmount: 10,
			},
		},
	}
	nutritionRes := httpdto.NutritionResponse{
		Calories: utils.VPtr(100),
		Proteins: utils.VPtr(10.5),
		Fats:     utils.VPtr(5.5),
		Carbs:    utils.VPtr(20.0),
	}
	expected := httpdto.GetCatalogResponse{
		Location: httpdto.LocationResponse{
			ID:        locID.String(),
			Slug:      "test_1",
			Name:      "test name of location",
			Address:   "address of test name of location",
			IsActive:  true,
			CreatedAt: now.String(),
			DeletedAt: nil,
		},
		Categories: []string{"drinks"},
		Items: []httpdto.CatalogItemResponse{
			{
				ItemID:      id.String(),
				Name:        "name",
				Description: utils.VPtr("desc"),
				Category:    "drinks",
				PhotoURL:    "url",
				Nutrition:   &nutritionRes,
				Price:       100,
				IsAvailable: true,
				StockAmount: 10,
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
			{ItemID: itemID1.String(), Amount: 1, Price: 200},
			{ItemID: itemID2.String(), Amount: 2, Price: 300},
		},
	}

	expected := appdto.CreateOrderInput{
		LocationID: locID,
		Items: []appdto.OrderItemInput{
			{ItemID: itemID1, Amount: 1, Price: 200},
			{ItemID: itemID2, Amount: 2, Price: 300},
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
		Location: httpdto.LocationResponse{
			ID: id.String(), Slug: "slug", Name: "name", Address: "address", IsActive: true, CreatedAt: now.String(),
		},
	}
	result := mapper.MapOutputToCreateLocation(out)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestMapRequestToUpdateLocation(t *testing.T) {
	req := httpdto.UpdateLocationRequest{
		Slug:     "slug",
		Name:     utils.VPtr("name"),
		Address:  utils.VPtr("address"),
		IsActive: utils.VPtr(true),
	}

	expected := appdto.UpdateLocationInput{
		Slug:     "slug",
		Name:     utils.VPtr("name"),
		Address:  utils.VPtr("address"),
		IsActive: utils.VPtr(true),
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
			ID:        id,
			Slug:      "slug",
			Name:      "name",
			Address:   "address",
			IsActive:  true,
			CreatedAt: now,
		},
	}

	expected := httpdto.UpdateLocationResponse{
		Location: httpdto.LocationResponse{
			ID:        id.String(),
			Slug:      "slug",
			Name:      "name",
			Address:   "address",
			IsActive:  true,
			CreatedAt: now.String(),
		},
	}

	result := mapper.MapOutputToUpdateLocation(out)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestMapRequestToDeleteLocation(t *testing.T) {
	req := httpdto.DeleteLocationRequest{Slug: "new_1"}
	expected := appdto.DeleteLocationInput{Slug: "new_1"}

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
			{
				ID:        id,
				Slug:      "slug",
				Name:      "name",
				Address:   "address",
				IsActive:  true,
				CreatedAt: now,
			},
		},
	}

	expected := httpdto.ListLocationsResponse{
		Locations: []httpdto.LocationResponse{
			{
				ID:        id.String(),
				Slug:      "slug",
				Name:      "name",
				Address:   "address",
				IsActive:  true,
				CreatedAt: now.String(),
			},
		},
	}

	result := mapper.MapOutputToListLocations(out)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestMapRequestToGetQRCode(t *testing.T) {
	req := httpdto.GetQRCodeRequest{Slug: "new_1"}
	expected := appdto.GetQRCodeInput{Slug: "new_1"}

	result := mapper.MapRequestToGetQRCode(req)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

// --- ITEMS & NUTRITION ---

func TestMapRequestToCreateItem(t *testing.T) {
	req := httpdto.CreateItemRequest{
		Name:        "name",
		Description: utils.VPtr("description"),
		Category:    "cat",
		PhotoURL:    "url",
		Nutrition:   &httpdto.NutritionRequest{Calories: utils.VPtr(100)},
	}

	expected := appdto.CreateItemInput{
		Name:        "name",
		Description: utils.VPtr("description"),
		Category:    "cat",
		PhotoURL:    "url",
		Nutrition:   &appdto.NutritionOutput{Calories: utils.VPtr(100)},
	}

	result := mapper.MapRequestToCreateItem(req)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestMapRequestToGetItem(t *testing.T) {
	id := uuid.New()

	req := httpdto.GetItemRequest{ID: id.String()}
	expected := appdto.GetItemInput{ID: id}

	result := mapper.MapRequestToGetItem(req)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestMapRequestToUpdateItem(t *testing.T) {
	id := uuid.New()

	req := httpdto.UpdateItemRequest{
		ID:        id.String(),
		Name:      utils.VPtr("new name"),
		Nutrition: &httpdto.NutritionRequest{Calories: utils.VPtr(200)},
	}

	expected := appdto.UpdateItemInput{
		ID:        id,
		Name:      utils.VPtr("new name"),
		Nutrition: &appdto.NutritionOutput{Calories: utils.VPtr(200)},
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
	out := appdto.CreateItemOutput{
		Item: appdto.ItemOutput{
			ID:        id,
			Name:      "name",
			Category:  "cat",
			Nutrition: &appdto.NutritionOutput{Calories: utils.VPtr(640)},
			CreatedAt: now,
		},
	}

	expected := httpdto.CreateItemResponse{
		Item: httpdto.ItemResponse{
			ID:        id.String(),
			Name:      "name",
			Category:  "cat",
			Nutrition: &httpdto.NutritionResponse{Calories: utils.VPtr(640)},
			CreatedAt: now.String(),
		},
	}

	result := mapper.MapOutputToCreateItem(out)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestMapOutputToGetItem(t *testing.T) {
	id := uuid.New()
	now := time.Now()

	out := appdto.GetItemOutput{
		Item: appdto.ItemOutput{
			ID:        id,
			Name:      "updated",
			Category:  "cat",
			CreatedAt: now,
		},
	}

	expected := httpdto.GetItemResponse{
		Item: httpdto.ItemResponse{
			ID:        id.String(),
			Name:      "updated",
			Category:  "cat",
			CreatedAt: now.String(),
		},
	}

	result := mapper.MapOutputToGetItem(out)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestMapOutputToUpdateItem(t *testing.T) {
	id := uuid.New()
	now := time.Now()

	out := appdto.UpdateItemOutput{
		Item: appdto.ItemOutput{
			ID:        id,
			Name:      "updated",
			Category:  "cat",
			CreatedAt: now,
		},
	}

	expected := httpdto.UpdateItemResponse{
		Item: httpdto.ItemResponse{
			ID:        id.String(),
			Name:      "updated",
			Category:  "cat",
			CreatedAt: now.String(),
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
			{
				ID:        id,
				Name:      "name",
				Category:  "cat",
				CreatedAt: now,
			},
		},
	}

	expected := httpdto.ListItemsResponse{
		Items: []httpdto.ItemResponse{
			{
				ID:        id.String(),
				Name:      "name",
				Category:  "cat",
				CreatedAt: now.String(),
			},
		},
	}

	result := mapper.MapOutputToListItems(out)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}
