package dto

type CreateItemInput struct {
	Name        string
	Description *string
	Category    string
	PhotoURL    string
	Nutrition   *NutritionResponse
}

type CreateItemOutput struct {
	Item ItemResponse
}
