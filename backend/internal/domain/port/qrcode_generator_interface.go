package port

import "github.com/google/uuid"

type QRCodeGenerator interface {
	Generate(locationID uuid.UUID) ([]byte, error)
}
