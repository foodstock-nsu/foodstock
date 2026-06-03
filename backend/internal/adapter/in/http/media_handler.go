package http

import (
	"backend/internal/adapter/in/http/mapper"
	appdto "backend/internal/app/dto"
	"backend/internal/app/usecase"
	pkgerrs "backend/pkg/errs"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
)

type MediaHandler struct {
	log           *slog.Logger
	uploadMediaUC *usecase.UploadMediaUC
}

func NewMediaHandler(log *slog.Logger, uploadMediaUC *usecase.UploadMediaUC) *MediaHandler {
	return &MediaHandler{
		log:           log,
		uploadMediaUC: uploadMediaUC,
	}
}

func (h *MediaHandler) Upload(c echo.Context) error {
	fileHeader, err := c.FormFile("image")
	if err != nil {
		return h.returnErr(c, "missing or invalid image field", pkgerrs.ErrInvalidImage)
	}

	file, err := fileHeader.Open()
	if err != nil {
		return h.returnErr(c, "failed to open uploaded file", pkgerrs.ErrInvalidImage)
	}
	defer file.Close()

	out, err := h.uploadMediaUC.Execute(
		c.Request().Context(),
		appdto.UploadMediaInput{
			File:        file,
			Filename:    fileHeader.Filename,
			ContentType: fileHeader.Header.Get("Content-Type"),
		})
	if err != nil {
		return h.returnErr(c, "failed to upload media", err)
	}

	h.log.InfoContext(
		c.Request().Context(), "media uploaded successfully",
		slog.String("filename", fileHeader.Filename),
		slog.String("s3_key", out.MediaKey),
	)

	return c.JSON(http.StatusCreated, mapper.MapOutputToUploadMedia(out))
}

func (h *MediaHandler) returnErr(c echo.Context, msg string, err error) error {
	outErr := mapper.HttpError(err)

	h.log.ErrorContext(c.Request().Context(), msg,
		slog.Int("code", outErr.Code),
		slog.String("public_msg", outErr.Message),
		slog.Any("cause", outErr.Reason),
	)

	return c.JSON(outErr.Code, mapper.MapErrorToResponse(outErr.Message))
}
