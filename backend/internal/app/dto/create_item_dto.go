package dto

type CreateItemInput struct {
	Name        string
	Description *string
	Category    string
	PhotoURL    string
	Nutrition   *NutritionDTO
}

type CreateItemOutput struct {
	Item ItemDTO
}
