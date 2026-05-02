package model

import (
	pkgerrs "backend/pkg/errs"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrCannotChangeStock = errors.New("stock amount cannot be less than zero")
)

// ================ Rich model for LocationItem ================

type LocationItem struct {
	id          uuid.UUID
	itemID      uuid.UUID
	locationID  uuid.UUID
	price       int64 // in kopecks
	isAvailable bool
	stockAmount int
}

func NewLocationItem(
	itemID, locID uuid.UUID,
	price int64,
	stockAmount int,
) (*LocationItem, error) {
	if itemID == uuid.Nil {
		return nil, pkgerrs.NewValueInvalidError("item_id")
	}
	if locID == uuid.Nil {
		return nil, pkgerrs.NewValueInvalidError("location_id")
	}
	if price < 0 {
		return nil, pkgerrs.NewValueInvalidError("total_price")
	}
	if stockAmount < 0 {
		return nil, pkgerrs.NewValueInvalidError("stock_amount")
	}

	var isAvailable bool
	if stockAmount != 0 {
		isAvailable = true
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

func RestoreLocationItem(
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

func (li *LocationItem) RestoreStock(amount int) error {
	if amount <= 0 {
		return ErrCannotChangeStock
	}

	li.stockAmount += amount
	li.isAvailable = true

	return nil
}

func (li *LocationItem) ReduceStock(amount int) error {
	if li.stockAmount-amount < 0 || amount <= 0 {
		return ErrCannotChangeStock
	}

	li.stockAmount -= amount
	if li.stockAmount == 0 {
		li.isAvailable = false
	}

	return nil
}

func (li *LocationItem) Update(
	price *int64,
	stockAmount *int,
) error {
	if price != nil && *price < 0 {
		return pkgerrs.NewValueInvalidError("totalPrice")
	}
	if stockAmount != nil && *stockAmount < 0 {
		return pkgerrs.NewValueInvalidError("stock_amount")
	}

	if price != nil {
		li.price = *price
	}
	if stockAmount != nil {
		li.stockAmount = *stockAmount
	}

	if li.stockAmount == 0 {
		li.isAvailable = false
	}

	return nil
}
