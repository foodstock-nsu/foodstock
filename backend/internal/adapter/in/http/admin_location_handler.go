package http

import (
	httpdto "backend/internal/adapter/in/http/dto"
	"backend/internal/adapter/in/http/mapper"
	"backend/internal/app/usecase"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AdminLocationHandler struct {
	log              *slog.Logger
	createLocationUC usecase.CreateLocationUC
	updateLocationUC usecase.UpdateLocationUC
	deleteLocationUC usecase.DeleteLocationUC
	listLocationsUC  usecase.ListLocationsUC
	getQRCodeUC      usecase.GetQRCodeUC
}

func NewAdminLocationHandler(
	log *slog.Logger,
	createLocationUC usecase.CreateLocationUC,
	updateLocationUC usecase.UpdateLocationUC,
	deleteLocationUC usecase.DeleteLocationUC,
	listLocationsUC usecase.ListLocationsUC,
	getQRCodeUC usecase.GetQRCodeUC,
) *AdminLocationHandler {
	return &AdminLocationHandler{
		log:              log,
		createLocationUC: createLocationUC,
		updateLocationUC: updateLocationUC,
		deleteLocationUC: deleteLocationUC,
		listLocationsUC:  listLocationsUC,
		getQRCodeUC:      getQRCodeUC,
	}
}

func (h *AdminLocationHandler) CreateLocation(c echo.Context) error {
	var req httpdto.CreateLocationRequest

	if err := c.Bind(&req); err != nil {
		return h.returnErr(c, "invalid query", err)
	}

	out, err := h.createLocationUC.Execute(
		c.Request().Context(),
		mapper.MapRequestToCreateLocation(req),
	)
	if err != nil {
		return h.returnErr(c, "failed to create location", err)
	}

	return c.JSON(http.StatusCreated, mapper.MapOutputToCreateLocation(out))
}

func (h *AdminLocationHandler) UpdateLocation(c echo.Context) error {
	var req httpdto.UpdateLocationRequest

	if err := c.Bind(&req); err != nil {
		return h.returnErr(c, "invalid query or json", err)
	}

	out, err := h.updateLocationUC.Execute(
		c.Request().Context(),
		mapper.MapRequestToUpdateLocation(req),
	)
	if err != nil {
		return h.returnErr(c, "failed to update location", err)
	}

	return c.JSON(http.StatusOK, mapper.MapOutputToUpdateLocation(out))
}

func (h *AdminLocationHandler) DeleteLocation(c echo.Context) error {
	var req httpdto.DeleteLocationRequest

	err := c.Bind(&req)
	if err != nil {
		return h.returnErr(c, "invalid query", err)
	}

	err = h.deleteLocationUC.Execute(
		c.Request().Context(),
		mapper.MapRequestToDeleteLocation(req),
	)
	if err != nil {
		return h.returnErr(c, "failed to delete location", err)
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *AdminLocationHandler) ListLocations(c echo.Context) error {
	out, err := h.listLocationsUC.Execute(c.Request().Context())
	if err != nil {
		return h.returnErr(c, "failed to get a list of locations", err)
	}

	return c.JSON(http.StatusOK, mapper.MapOutputToListLocations(out))
}

func (h *AdminLocationHandler) GetQRCode(c echo.Context) error {
	var req httpdto.GetQRCodeRequest

	if err := c.Bind(&req); err != nil {
		return h.returnErr(c, "invalid query", err)
	}

	out, err := h.getQRCodeUC.Execute(
		c.Request().Context(),
		mapper.MapRequestToGetQRCode(req),
	)
	if err != nil {
		return h.returnErr(c, "failed to get a qr code", err)
	}

	return c.Blob(http.StatusOK, "image/png", out.QRCode)
}

func (h *AdminLocationHandler) returnErr(c echo.Context, msg string, err error) error {
	outErr := mapper.HttpError(err)

	h.log.ErrorContext(c.Request().Context(), msg,
		slog.Int("code", outErr.Code),
		slog.String("public_msg", outErr.Message),
		slog.Any("cause", outErr.Reason),
	)

	return c.JSON(outErr.Code, map[string]string{"error": msg})
}
