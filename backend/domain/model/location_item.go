package model

import (
	pkgerrs "backend/pkg/errs"
	"errors"

	"github.com/google/uuid"
)

var ErrCannotReduceStock = errors.New("stock cannot be less than zero")

// ================ Rich model for LocationItem ================

type LocationItem struct {
	id          uuid.UUID
	itemID      uuid.UUID
	locationID  uuid.UUID
	price       int64 // in kopecks
	isAvailable bool
	stockAmount int
}

func NewItemPrice(
	itemID, locID uuid.UUID,
	price int64,
	isAvailable bool,
	stockAmount int,
) (*LocationItem, error) {
	if itemID == uuid.Nil {
		return nil, pkgerrs.NewValueInvalidError("item_id")
	}
	if locID == uuid.Nil {
		return nil, pkgerrs.NewValueInvalidError("location_id")
	}
	if price < 0 {
		return nil, pkgerrs.NewValueInvalidError("price")
	}
	if stockAmount == 0 && isAvailable {
		return nil, pkgerrs.NewValueInvalidError("is_available")
	}
	if stockAmount < 0 {
		return nil, pkgerrs.NewValueInvalidError("stock_amount")
	}

	return &LocationItem{
		id:          uuid.New(),
		itemID:      itemID,
		locationID:  locID,
		price:       price,
		isAvailable: isAvailable,
		stockAmount: stockAmount,
	}, nil
}

func RestoreItemPrice(
	id, itemID, locID uuid.UUID,
	price int64,
	isAvailable bool,
	stockAmount int,
) *LocationItem {
	return &LocationItem{
		id:          id,
		itemID:      itemID,
		locationID:  locID,
		price:       price,
		isAvailable: isAvailable,
		stockAmount: stockAmount,
	}
}

// ================ Read-Only ================

func (li *LocationItem) ID() uuid.UUID         { return li.id }
func (li *LocationItem) ItemID() uuid.UUID     { return li.itemID }
func (li *LocationItem) LocationID() uuid.UUID { return li.locationID }
func (li *LocationItem) Price() int64          { return li.price }
func (li *LocationItem) IsAvailable() bool     { return li.isAvailable }
func (li *LocationItem) StockAmount() int      { return li.stockAmount }

// ================ Business Logic ================

func (li *LocationItem) CanBeSold() bool {
	return li.isAvailable && li.stockAmount > 0
}

// ================ Mutation ================

func (li *LocationItem) ReduceStock(amount int) error {
	if li.stockAmount-amount < 0 {
		return ErrCannotReduceStock
	}
	li.stockAmount -= amount
	return nil
}
