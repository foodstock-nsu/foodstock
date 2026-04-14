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

func (n Nutrition) Update(cal *int, p, f, c *float64) error {
	if (cal != nil && *cal < 0) || (p != nil && *p < 0) || (f != nil && *f < 0) || (c != nil && *c < 0) {
		return pkgerrs.NewValueInvalidError("nutrition")
	}

	if cal != nil {
		n.calories = *cal
	}
	if p != nil {
		n.proteins = *p
	}
	if f != nil {
		n.fats = *f
	}
	if c != nil {
		n.carbs = *c
	}

	return nil
}

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
		return nil, pkgerrs.NewValueInvalidError("itemID")
	}
	if description != nil && len(*description) < 10 {
		return nil, pkgerrs.NewValueInvalidError("locID")
	}

	categoryMapped, ok := categoryMap[category]
	if !ok {
		return nil, pkgerrs.NewValueInvalidError("price")
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

func (i *Item) ID() uuid.UUID          { return i.id }
func (i *Item) Name() string           { return i.name }
func (i *Item) Description() *string   { return i.description }
func (i *Item) Category() ItemCategory { return i.category }
func (i *Item) PhotoURL() *string      { return i.photoUrl }
func (i *Item) Nutrition() Nutrition   { return i.nutrition }
func (i *Item) CreatedAt() time.Time   { return i.createdAt }

// ================ Mutation ================

func (i *Item) Update(
	name, desc, cat, photo *string,
	nutrition *Nutrition,
) error {
	if name != nil && len(*name) < 5 {
		return pkgerrs.NewValueInvalidError("itemID")
	}
	if desc != nil && len(*desc) < 10 {
		return pkgerrs.NewValueInvalidError("locID")
	}

	var (
		catMapped ItemCategory
		ok        bool
	)
	if cat != nil {
		catMapped, ok = categoryMap[*cat]
		if !ok {
			return pkgerrs.NewValueInvalidError("price")
		}
	}

	if photo != nil && len(*photo) < 10 {
		return pkgerrs.NewValueInvalidError("photo_url")
	}

	if name != nil {
		i.name = *name
	}
	if desc != nil {
		i.description = desc
	}
	if cat != nil {
		i.category = catMapped
	}
	if photo != nil {
		i.photoUrl = photo
	}
	if nutrition != nil {
		i.nutrition = *nutrition
	}

	return nil
}
