package model

import (
	pkgerrs "backend/pkg/errs"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrItemNotAvailable       = errors.New("item not available")
	ErrOrderCannotBePaid      = errors.New("order is either already paid or cancelled")
	ErrOrderCannotBeCancelled = errors.New("order is either already cancelled or paid")
)

// ================ Value Objects ================

type OrderStatus string

const (
	OrderPending   OrderStatus = "PENDING"
	OrderPaid      OrderStatus = "PAID"
	OrderCancelled OrderStatus = "CANCELLED"
)

type OrderItem struct {
	id              uuid.UUID
	itemID          uuid.UUID
	itemAmount      int
	priceAtPurchase int64
}

func NewOrderItem(
	itemID uuid.UUID,
	amount int,
	price int64,
) (*OrderItem, error) {
	if itemID == uuid.Nil {
		return nil, pkgerrs.NewValueInvalidError("item_id")
	}
	if amount < 0 {
		return nil, pkgerrs.NewValueInvalidError("item_amount")
	}
	if price < 0 {
		return nil, pkgerrs.NewValueInvalidError("price_at_purchase")
	}
	return &OrderItem{
		id:              uuid.New(),
		itemID:          itemID,
		itemAmount:      amount,
		priceAtPurchase: price,
	}, nil
}

func RestoreOrderItem(
	id uuid.UUID,
	itemID uuid.UUID,
	itemAmount int,
	priceAtPurchase int64,
) *OrderItem {
	return &OrderItem{
		id:              id,
		itemID:          itemID,
		itemAmount:      itemAmount,
		priceAtPurchase: priceAtPurchase,
	}
}

// ================ Rich model for Order ================

type Order struct {
	id         uuid.UUID
	locationID uuid.UUID
	items      []OrderItem
	status     OrderStatus
	totalPrice int64
	createdAt  time.Time
	paidAt     *time.Time
}

func NewOrder(locID uuid.UUID, items []OrderItem, totalPrice int64) (*Order, error) {
	if locID == uuid.Nil {
		return nil, pkgerrs.NewValueInvalidError("location_id")
	}
	if items == nil {
		return nil, pkgerrs.NewValueRequiredError("order_items")
	}
	if totalPrice < 0 {
		return nil, pkgerrs.NewValueInvalidError("total_price")
	}
	return &Order{
		id:         uuid.New(),
		locationID: locID,
		items:      items,
		status:     OrderPending,
		totalPrice: totalPrice,
		createdAt:  time.Now().UTC(),
		paidAt:     nil,
	}, nil
}

func RestoreOrder(
	id uuid.UUID,
	locationID uuid.UUID,
	items []OrderItem,
	status OrderStatus,
	totalPrice int64,
	createdAt time.Time,
	paidAt *time.Time,
) *Order {
	return &Order{
		id:         id,
		locationID: locationID,
		items:      items,
		status:     status,
		totalPrice: totalPrice,
		createdAt:  createdAt,
		paidAt:     paidAt,
	}
}

// ================ Read-Only ================

func (o *Order) ID() uuid.UUID         { return o.id }
func (o *Order) LocationID() uuid.UUID { return o.locationID }
func (o *Order) Items() []OrderItem    { return o.items }
func (o *Order) Status() OrderStatus   { return o.status }
func (o *Order) TotalPrice() int64     { return o.totalPrice }
func (o *Order) CreatedAt() time.Time  { return o.createdAt }
func (o *Order) PaidAt() *time.Time    { return o.paidAt }

// ================ Business Logic ================

func (o *Order) calculateTotal() {
	var total int64
	for i := range o.items {
		total += o.items[i].priceAtPurchase * int64(o.items[i].itemAmount)
	}
	o.totalPrice = total
}

func (o *Order) AddItem(locationItem LocationItem, quantity int) error {
	if locationItem.LocationID() != o.locationID || !locationItem.CanBeSold() {
		return ErrItemNotAvailable
	}

	item := OrderItem{
		id:              uuid.New(),
		itemID:          locationItem.ItemID(),
		itemAmount:      quantity,
		priceAtPurchase: locationItem.Price(),
	}

	o.items = append(o.items, item)
	o.calculateTotal()

	return nil
}

// ================ Mutation ================

func (o *Order) Pay() error {
	if o.status != OrderPending {
		return ErrOrderCannotBePaid
	}
	o.status = OrderPaid

	now := time.Now().UTC()
	o.paidAt = &now

	return nil
}

func (o *Order) Cancel() error {
	if o.status != OrderPending {
		return ErrOrderCannotBeCancelled
	}
	o.status = OrderCancelled
	return nil
}
