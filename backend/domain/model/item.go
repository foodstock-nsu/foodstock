package model

import (
	pkgerrs "backend/pkg/errs"
	"time"

	"github.com/google/uuid"
)

type ItemCategory string

const (
	ItemLunch     ItemCategory = "lunch"
	ItemBreakfast ItemCategory = "breakfast"
	ItemDrinks    ItemCategory = "drinks"
	ItemSnacks    ItemCategory = "snacks"
	ItemDesserts  ItemCategory = "desserts"
)

var categoryMap = map[string]ItemCategory{
	"lunch":     ItemLunch,
	"breakfast": ItemBreakfast,
	"drinks":    ItemDrinks,
	"snacks":    ItemSnacks,
	"desserts":  ItemDesserts,
}

// ================ Value Object - Nutrition ================

type Nutrition struct {
	calories int
	proteins float64
	fats     float64
	carbs    float64
}

func NewNutrition(cal int, p, f, c float64) (Nutrition, error) {
	if cal < 0 || p < 0 || f < 0 || c < 0 {
		return Nutrition{}, pkgerrs.NewValueInvalidError("nutrition")
	}
	return Nutrition{
		calories: cal,
		proteins: p,
		fats:     f,
		carbs:    c,
	}, nil
}

func RestoreNutrition(cal int, p, f, c float64) Nutrition {
	return Nutrition{
		calories: cal,
		proteins: p,
		fats:     f,
		carbs:    c,
	}
}

func (n Nutrition) Calories() int     { return n.calories }
func (n Nutrition) Proteins() float64 { return n.proteins }
func (n Nutrition) Fats() float64     { return n.fats }
func (n Nutrition) Carbs() float64    { return n.carbs }

// ================ Rich model for Item ================

type Item struct {
	id          uuid.UUID
	name        string
	description *string
	category    ItemCategory
	photoUrl    *string
	nutrition   Nutrition
	createdAt   time.Time
}

func NewItem(
	name string,
	description *string,
	category string,
	photoUrl *string,
	nutrition Nutrition,
) (*Item, error) {
	if len(name) < 5 {
		return nil, pkgerrs.NewValueInvalidError("name")
	}
	if description != nil && len(*description) < 10 {
		return nil, pkgerrs.NewValueInvalidError("description")
	}

	categoryMapped, ok := categoryMap[category]
	if !ok {
		return nil, pkgerrs.NewValueInvalidError("category")
	}

	if photoUrl != nil && len(*photoUrl) < 10 {
		return nil, pkgerrs.NewValueInvalidError("photo_url")
	}

	return &Item{
		id:          uuid.New(),
		name:        name,
		description: description,
		category:    categoryMapped,
		photoUrl:    photoUrl,
		nutrition:   nutrition,
		createdAt:   time.Now().UTC(),
	}, nil
}

func RestoreItem(
	id uuid.UUID,
	name string,
	description *string,
	category ItemCategory,
	photoUrl *string,
	nutrition Nutrition,
	createdAt time.Time,
) *Item {
	return &Item{
		id:          id,
		name:        name,
		description: description,
		category:    category,
		photoUrl:    photoUrl,
		nutrition:   nutrition,
		createdAt:   createdAt,
	}
}

// ================ Read-Only ================

func (p *Item) ID() uuid.UUID          { return p.id }
func (p *Item) Name() string           { return p.name }
func (p *Item) Description() *string   { return p.description }
func (p *Item) Category() ItemCategory { return p.category }
func (p *Item) PhotoURL() *string      { return p.photoUrl }
func (p *Item) Nutrition() Nutrition   { return p.nutrition }
func (p *Item) CreatedAt() time.Time   { return p.createdAt }
