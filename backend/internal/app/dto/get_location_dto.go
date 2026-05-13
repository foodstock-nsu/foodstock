package dto

type GetLocationInput struct {
	Slug string
}

type GetLocationOutput struct {
	Location LocationOutput
}
