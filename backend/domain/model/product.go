package model

import (
	pkgerrs "backend/pkg/errs"
	"time"

	"github.com/google/uuid"
)

type ProductCategory string

const (
	ProductDrinks   ProductCategory = "drinks"
	ProductSnacks   ProductCategory = "snacks"
	ProductDesserts ProductCategory = "desserts"
)

var categoryMap = map[string]ProductCategory{
	"drinks":   ProductDrinks,
	"snacks":   ProductSnacks,
	"desserts": ProductDesserts,
}

// ================ Rich model for Location (e.g. Fridge) ================

type Product struct {
	id          uuid.UUID
	name        string
	description *string
	category    ProductCategory
	photoUrl    *string
	calories    int
	createdAt   time.Time
}

func NewProduct(
	name string,
	description *string,
	category string,
	photoUrl *string,
	calories int,
) (*Product, error) {
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
	if calories < 0 {
		return nil, pkgerrs.NewValueInvalidError("calories")
	}

	return &Product{
		id:          uuid.New(),
		name:        name,
		description: description,
		category:    categoryMapped,
		photoUrl:    photoUrl,
		calories:    calories,
		createdAt:   time.Now().UTC(),
	}, nil
}

func RestoreProduct(
	id uuid.UUID,
	name string,
	description *string,
	category ProductCategory,
	photoUrl *string,
	calories int,
	createdAt time.Time,
) *Product {
	return &Product{
		id:          id,
		name:        name,
		description: description,
		category:    category,
		photoUrl:    photoUrl,
		calories:    calories,
		createdAt:   createdAt,
	}
}

// ================ Read-Only ================

func (p *Product) ID() uuid.UUID             { return p.id }
func (p *Product) Name() string              { return p.name }
func (p *Product) Description() *string      { return p.description }
func (p *Product) Category() ProductCategory { return p.category }
func (p *Product) PhotoURL() *string         { return p.photoUrl }
func (p *Product) Calories() int             { return p.calories }
func (p *Product) CreatedAt() time.Time      { return p.createdAt }
