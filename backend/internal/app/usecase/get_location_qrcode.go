package usecase

import (
	"backend/internal/app/dto"
	ucerrs "backend/internal/app/errs"
	"backend/internal/domain/port"
	"context"
)

type GetLocationQRCodeUC struct {
	location port.LocationRepository
	qrcode   port.QRCodeGenerator
}

func (uc *GetLocationQRCodeUC) Execute(ctx context.Context, in dto.GetLocationQRCodeInput) (dto.GetLocationQRCodeOutput, error) {
	// Get a location by id
	location, err := uc.location.GetByID(ctx, in.LocationID)
	if err != nil {
		return dto.GetLocationQRCodeOutput{}, ucerrs.Wrap(
			ucerrs.ErrGetLocationByIDDB, err,
		)
	}

	// Validation
	if !location.IsOperational() {
		return dto.GetLocationQRCodeOutput{}, ucerrs.ErrCannotGetLocationQRCode
	}

	// Generate QR-code
	qrcode, err := uc.qrcode.Generate(location.Slug())
	if err != nil {
		return dto.GetLocationQRCodeOutput{}, ucerrs.Wrap(
			ucerrs.ErrGenerateQRCode, err,
		)
	}

	return dto.GetLocationQRCodeOutput{QRCode: qrcode}, nil
}
