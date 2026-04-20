package dto

import "github.com/google/uuid"

type GetQRCodeInput struct {
	LocationID uuid.UUID
}

type GetQRCodeOutput struct {
	QRCode []byte
}
