package dto

import "github.com/google/uuid"

type GetLocationQRCodeInput struct {
	LocationID uuid.UUID
}

type GetLocationQRCodeOutput struct {
	QRCode []byte
}
