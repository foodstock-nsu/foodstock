package dto

import "github.com/google/uuid"

type UpdateItemInput struct {
	ID          uuid.UUID
	Name        *string
	Description *string
	Category    *string
	PhotoURL    *string
	Nutrition   *NutritionDTO
}

type UpdateItemOutput struct {
	Item ItemDTO
}
