package model

import (
	pkgerrs "backend/pkg/errs"
	"time"

	"github.com/google/uuid"
)

type ItemCategory string

func (ic ItemCategory) String() string { return string(ic) }

const (
	ItemLunch     ItemCategory = "lunch"
	ItemBreakfast ItemCategory = "breakfast"
	ItemDrinks    ItemCategory = "drinks"
	ItemSnacks    ItemCategory = "snacks"
	ItemDesserts  ItemCategory = "desserts"
	ItemOther     ItemCategory = "other"
)

var categoryMap = map[string]ItemCategory{
	"lunch":     ItemLunch,
	"breakfast": ItemBreakfast,
	"drinks":    ItemDrinks,
	"snacks":    ItemSnacks,
	"desserts":  ItemDesserts,
	"other":     ItemOther,
}

// ================ Value Object - Nutrition ================

type Nutrition struct {
	calories *int
	proteins *float64
	fats     *float64
	carbs    *float64
}

func NewNutrition(calories *int, proteins, fats, carbs *float64) (*Nutrition, error) {
	if calories == nil && proteins == nil && fats == nil && carbs == nil {
		return nil, nil
	}

	if calories != nil && *calories < 0 {
		return nil, pkgerrs.NewValueInvalidError("calories")
	}
	if proteins != nil && *proteins < 0 {
		return nil, pkgerrs.NewValueInvalidError("proteins")
	}
	if fats != nil && *fats < 0 {
		return nil, pkgerrs.NewValueInvalidError("fats")
	}
	if carbs != nil && *carbs < 0 {
		return nil, pkgerrs.NewValueInvalidError("carbs")
	}

	return &Nutrition{
		calories: calories,
		proteins: proteins,
		fats:     fats,
		carbs:    carbs,
	}, nil
}

func RestoreNutrition(calories *int, proteins, fats, carbs *float64) *Nutrition {
	if calories == nil && proteins == nil && fats == nil && carbs == nil {
		return nil
	}

	return &Nutrition{
		calories: calories,
		proteins: proteins,
		fats:     fats,
		carbs:    carbs,
	}
}

func (n Nutrition) Calories() *int     { return n.calories }
func (n Nutrition) Proteins() *float64 { return n.proteins }
func (n Nutrition) Fats() *float64     { return n.fats }
func (n Nutrition) Carbs() *float64    { return n.carbs }

func (n Nutrition) Update(calories *int, proteins, fats, carbs *float64) error {
	if calories != nil && *calories < 0 {
		return pkgerrs.NewValueInvalidError("calories")
	}
	if proteins != nil && *proteins < 0 {
		return pkgerrs.NewValueInvalidError("proteins")
	}
	if fats != nil && *fats < 0 {
		return pkgerrs.NewValueInvalidError("fats")
	}
	if carbs != nil && *carbs < 0 {
		return pkgerrs.NewValueInvalidError("carbs")
	}

	if calories != nil {
		n.calories = calories
	}
	if proteins != nil {
		n.proteins = proteins
	}
	if fats != nil {
		n.fats = fats
	}
	if carbs != nil {
		n.carbs = carbs
	}

	return nil
}

// ================ Rich model for Item ================

type Item struct {
	id          uuid.UUID
	name        string
	description *string
	category    ItemCategory
	photoUrl    string
	nutrition   *Nutrition
	createdAt   time.Time
}

func NewItem(
	name string,
	description *string,
	category string,
	photoUrl string,
	nutrition *Nutrition,
) (*Item, error) {
	if len(name) < 5 {
		return nil, pkgerrs.NewValueInvalidError("locID")
	}
	if description != nil && len(*description) < 10 {
		return nil, pkgerrs.NewValueInvalidError("locID")
	}

	categoryMapped, ok := categoryMap[category]
	if !ok {
		return nil, pkgerrs.NewValueInvalidError("totalPrice")
	}

	if len(photoUrl) < 10 {
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
	photoUrl string,
	nutrition *Nutrition,
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
func (i *Item) PhotoURL() string       { return i.photoUrl }
func (i *Item) Nutrition() *Nutrition  { return i.nutrition }
func (i *Item) CreatedAt() time.Time   { return i.createdAt }

// ================ Mutation ================

func (i *Item) Update(
	name, desc, cat, photo *string,
	nutrition *Nutrition,
) error {
	if name != nil && len(*name) < 5 {
		return pkgerrs.NewValueInvalidError("locID")
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
			return pkgerrs.NewValueInvalidError("totalPrice")
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
		i.photoUrl = *photo
	}
	if nutrition != nil {
		i.nutrition = nutrition
	}

	return nil
}
