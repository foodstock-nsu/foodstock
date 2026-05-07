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
	// Get a location by slug and validate it
	location, err := uc.location.GetBySlug(ctx, in.Slug)
	if err != nil {
		if errors.Is(err, pkgerrs.ErrObjectNotFound) {
			return dto.GetQRCodeOutput{}, ucerrs.ErrLocationNotFound
		}
		return dto.GetQRCodeOutput{}, ucerrs.Wrap(
			ucerrs.ErrGetLocationBySlugDB, err,
		)
	}

	if location.IsDeleted() {
		return dto.GetQRCodeOutput{}, ucerrs.ErrLocationAlreadyDeleted
	}

	if !location.IsOperational() {
		return dto.GetQRCodeOutput{}, ucerrs.ErrLocationIsNotOperational
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
