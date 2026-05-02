package dto

type CreateLocationInput struct {
	Slug    string
	Name    string
	Address string
}

type CreateLocationOutput struct {
	Location LocationResponse
}
