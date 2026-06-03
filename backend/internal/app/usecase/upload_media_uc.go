package usecase

import (
	"backend/internal/app/dto"
	ucerrs "backend/internal/app/errs"
	"backend/internal/domain/port"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type UploadMediaUC struct {
	mediaStorageGateway port.MediaStorageGateway
}

func NewUploadMediaUC(mediaStorageGateway port.MediaStorageGateway) *UploadMediaUC {
	return &UploadMediaUC{mediaStorageGateway: mediaStorageGateway}
}

func (uc *UploadMediaUC) Execute(ctx context.Context, in dto.UploadMediaInput) (dto.UploadMediaOutput, error) {
	const maxFileSize = 5 * 1024 * 1024
	limitedReader := io.LimitReader(in.File, maxFileSize)

	// Validation: defense against non-image files
	buf := make([]byte, 512)
	n, err := limitedReader.Read(buf)
	if err != nil && !errors.Is(err, io.EOF) {
		return dto.UploadMediaOutput{}, ucerrs.Wrap(
			ucerrs.ErrInvalidInput, err,
		)
	}

	fileType := http.DetectContentType(buf)
	if fileType != "image/jpeg" && fileType != "image/png" && fileType != "image/webp" {
		return dto.UploadMediaOutput{}, ucerrs.Wrap(
			ucerrs.ErrInvalidInput,
			fmt.Errorf("invalid file type: %s, allowed only jpeg, png and webp", fileType),
		)
	}

	fullReader := io.MultiReader(bytes.NewReader(buf[:n]), limitedReader)

	// Get an image key
	s3Key, err := uc.mediaStorageGateway.Upload(
		ctx, fullReader,
		in.Filename, in.ContentType,
	)
	if err != nil {
		return dto.UploadMediaOutput{}, ucerrs.Wrap(
			ucerrs.ErrUploadMedia, err,
		)
	}

	return dto.UploadMediaOutput{MediaKey: s3Key}, nil
}
