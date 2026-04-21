package dto

type CreateItemInput struct {
	Name        string
	Description *string
	Category    string
	PhotoURL    string
	Nutrition   *NutritionOutput
}

type CreateItemOutput struct {
	Item ItemOutput
}
