package dto

import "github.com/google/uuid"

type OrderItemInput struct {
	ItemID uuid.UUID
	Amount int
	Price  int64
}

type CreateOrderInput struct {
	LocationID uuid.UUID
	Items      []OrderItemInput
}

type CreateOrderOutput struct {
	OrderID    uuid.UUID
	TotalPrice int64
	PaymentURL string
}
