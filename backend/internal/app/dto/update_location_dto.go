package dto

import "github.com/google/uuid"

type UpdateLocationInput struct {
	ID       uuid.UUID
	Slug     *string
	Name     *string
	Address  *string
	IsActive *bool
}

type UpdateLocationOutput struct {
	Location LocationOutput
}
