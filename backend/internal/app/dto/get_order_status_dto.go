package dto

import "github.com/google/uuid"

type GetOrderStatusInput struct {
	OrderID uuid.UUID
}

type GetOrderStatusOutput struct {
	Status string
}
