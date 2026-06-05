package mapper_test

import (
	"backend/internal/app/dto"
	"backend/internal/app/mapper"
	"backend/internal/domain/model"
	"backend/pkg/utils"
	"reflect"
	"testing"

	"github.com/google/uuid"
)

func TestMapDomainToLocationDTO(t *testing.T) {
	location, _ := model.NewLocation(
		"test_1",
		"Test Location",
		"Address Of Test Location",
	)
	expected := dto.LocationOutput{
		ID:        location.ID(),
		Slug:      "test_1",
		Name:      "Test Location",
		Address:   "Address Of Test Location",
		IsActive:  location.IsActive(),
		CreatedAt: location.CreatedAt(),
		DeletedAt: location.DeletedAt(),
	}

	mapped := mapper.MapDomainToLocationDTO(location)

	if !reflect.DeepEqual(mapped, expected) {
		t.Errorf("expected %v, got %v", expected, mapped)
	}
}

func TestMapDomainToLocationListDTO(t *testing.T) {
	location1, _ := model.NewLocation("test_1", "Test Location 1", "Address Of Test Location 1")
	location2, _ := model.NewLocation("test_2", "Test Location 2", "Address Of Test Location 2")
	locations := []*model.Location{location1, location2}

	expected := []dto.LocationOutput{
		{
			ID:        location1.ID(),
			Slug:      "test_1",
			Name:      "Test Location 1",
			Address:   "Address Of Test Location 1",
			IsActive:  location1.IsActive(),
			CreatedAt: location1.CreatedAt(),
			DeletedAt: location1.DeletedAt(),
		},
		{
			ID:        location2.ID(),
			Slug:      "test_2",
			Name:      "Test Location 2",
			Address:   "Address Of Test Location 2",
			IsActive:  location2.IsActive(),
			CreatedAt: location2.CreatedAt(),
			DeletedAt: location2.DeletedAt(),
		},
	}

	mapped := mapper.MapDomainToLocationListDTO(locations)

	if !reflect.DeepEqual(mapped, expected) {
		t.Errorf("expected %v, got %v", expected, mapped)
	}
}

func TestMapDomainToItemDTO(t *testing.T) {
	item, _ := model.NewItem("Сэндвич с рыбой", nil, "lunch", "http://photo-stock/a1.jpg", nil)
	expected := dto.ItemOutput{
		ID:          item.ID(),
		Name:        "Сэндвич с рыбой",
		Description: nil,
		Category:    "lunch",
		PhotoURL:    "http://photo-stock/a1.jpg",
		Nutrition:   nil,
		CreatedAt:   item.CreatedAt(),
		DeletedAt:   item.DeletedAt(),
	}

	mapped := mapper.MapDomainToItemDTO(item)

	if !reflect.DeepEqual(mapped, expected) {
		t.Errorf("expected %v, got %v", expected, mapped)
	}
}

func TestMapDomainToItemListDTO(t *testing.T) {
	item1, _ := model.NewItem("Сэндвич с рыбой", nil, "lunch", "http://photo-stock/a1.jpg", nil)
	item2, _ := model.NewItem("Сэндвич с сыром", nil, "lunch", "http://photo-stock/a2.jpg", nil)
	locations := []*model.Item{item1, item2}

	expected := []dto.ItemOutput{
		{
			ID:          item1.ID(),
			Name:        "Сэндвич с рыбой",
			Description: nil,
			Category:    "lunch",
			PhotoURL:    "http://photo-stock/a1.jpg",
			Nutrition:   nil,
			CreatedAt:   item1.CreatedAt(),
			DeletedAt:   item1.DeletedAt(),
		},
		{
			ID:          item2.ID(),
			Name:        "Сэндвич с сыром",
			Description: nil,
			Category:    "lunch",
			PhotoURL:    "http://photo-stock/a2.jpg",
			Nutrition:   nil,
			CreatedAt:   item2.CreatedAt(),
			DeletedAt:   item2.DeletedAt(),
		},
	}

	mapped := mapper.MapDomainToItemListDTO(locations)

	if !reflect.DeepEqual(mapped, expected) {
		t.Errorf("expected %v, got %v", expected, mapped)
	}
}

func TestMapDomainToCatalogItemDTO(t *testing.T) {
	nutrition, _ := model.NewNutrition(utils.VPtr(200), utils.VPtr(float64(20)), utils.VPtr(float64(20)), utils.VPtr(float64(20)))
	item, _ := model.NewItem("Сэндвич с рыбой", nil, "lunch", "http://photo-stock/a1.jpg", nutrition)
	locItem, _ := model.NewLocationItem(uuid.New(), item.ID(), 15000, 10)

	expected := dto.CatalogItemOutput{
		ItemID:      item.ID(),
		Name:        "Сэндвич с рыбой",
		Description: nil,
		Category:    "lunch",
		PhotoURL:    "http://photo-stock/a1.jpg",
		Nutrition: &dto.NutritionOutput{
			Calories: utils.VPtr(200),
			Proteins: utils.VPtr(float64(20)),
			Fats:     utils.VPtr(float64(20)),
			Carbs:    utils.VPtr(float64(20)),
		},
		Price:       locItem.Price(),
		IsAvailable: locItem.IsAvailable(),
		StockAmount: locItem.StockAmount(),
	}

	mapped := mapper.MapDomainToCatalogItemDTO(locItem, item)

	if !reflect.DeepEqual(mapped, expected) {
		t.Errorf("expected %v, got %v", expected, mapped)
	}
}

func TestMapDomainToInventoryItemListDTO(t *testing.T) {
	locID := uuid.New()
	locItem1, _ := model.NewLocationItem(uuid.New(), locID, 15000, 10)
	locItem2, _ := model.NewLocationItem(uuid.New(), locID, 20000, 5)
	items := []*model.LocationItem{locItem1, locItem2}

	expected := []dto.InventoryItemOutput{
		{
			ItemID:      locItem1.ItemID(),
			Price:       locItem1.Price(),
			IsAvailable: locItem1.IsAvailable(),
			StockAmount: locItem1.StockAmount(),
		},
		{
			ItemID:      locItem2.ItemID(),
			Price:       locItem2.Price(),
			IsAvailable: locItem2.IsAvailable(),
			StockAmount: locItem2.StockAmount(),
		},
	}

	mapped := mapper.MapDomainToInventoryItemListDTO(items)

	if !reflect.DeepEqual(mapped, expected) {
		t.Errorf("expected %v, got %v", expected, mapped)
	}
}
