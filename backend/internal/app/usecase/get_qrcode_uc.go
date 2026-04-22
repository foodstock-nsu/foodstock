package usecase

import (
	"backend/internal/app/dto"
	ucerrs "backend/internal/app/errs"
	"backend/internal/domain/port"
	pkgerrs "backend/pkg/errs"
	"context"
	"errors"
)

type GetQRCodeUC struct {
	location port.LocationRepository
	qrcode   port.QRCodeGenerator
}

func NewGetQRCodeUC(
	location port.LocationRepository,
	qrcode port.QRCodeGenerator,
) *GetQRCodeUC {
	return &GetQRCodeUC{
		location: location,
		qrcode:   qrcode,
	}

}

func (uc *GetQRCodeUC) Execute(ctx context.Context, in dto.GetQRCodeInput) (dto.GetQRCodeOutput, error) {
	// Get a location by id
	location, err := uc.location.GetByID(ctx, in.LocationID)
	if err != nil {
		if errors.Is(err, pkgerrs.ErrObjectNotFound) {
			return dto.GetQRCodeOutput{}, ucerrs.ErrLocationNotFound
		}
		return dto.GetQRCodeOutput{}, ucerrs.Wrap(
			ucerrs.ErrGetLocationByIDDB, err,
		)
	}

	// Validation
	if !location.IsOperational() {
		return dto.GetQRCodeOutput{}, ucerrs.ErrCannotGetLocationQRCode
	}

	// Generate QR-code
	qrcode, err := uc.qrcode.Generate(location.Slug())
	if err != nil {
		return dto.GetQRCodeOutput{}, ucerrs.Wrap(
			ucerrs.ErrGenerateQRCode, err,
		)
	}

	return dto.GetQRCodeOutput{QRCode: qrcode}, nil
}
