package dto

type CreateItemInput struct {
	Name        string
	Description *string
	Category    string
	PhotoURL    string
	Nutrition   *ItemNutrition
}

type CreateItemOutput struct {
	Item Item
}
