package usecase

import "backend/internal/domain/port"

type DeleteLocationUC struct {
	location     port.LocationRepository
	locationItem port.LocationItemRepository
}

func NewDeleteLocationUC(
	location port.LocationRepository,
	locationItem port.LocationItemRepository,
) *DeleteLocationUC {
	return &DeleteLocationUC{
		location:     location,
		locationItem: locationItem,
	}
}
