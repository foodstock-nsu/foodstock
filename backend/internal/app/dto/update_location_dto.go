package dto

type UpdateLocationInput struct {
	Slug     string
	Name     *string
	Address  *string
	IsActive *bool
}

type UpdateLocationOutput struct {
	Location LocationOutput
}
