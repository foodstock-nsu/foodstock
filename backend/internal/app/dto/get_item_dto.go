package dto

import "github.com/google/uuid"

type GetItemInput struct {
	ID uuid.UUID
}

type GetItemOutput struct {
	Item ItemOutput
}
