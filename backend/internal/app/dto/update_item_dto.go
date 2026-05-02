package dto

import "github.com/google/uuid"

type UpdateItemInput struct {
	ID          uuid.UUID
	Name        *string
	Description *string
	Category    *string
	PhotoURL    *string
	Nutrition   *NutritionResponse
}

type UpdateItemOutput struct {
	Item ItemResponse
}
