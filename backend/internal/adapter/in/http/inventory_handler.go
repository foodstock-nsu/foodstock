package http

import (
	httpdto "backend/internal/adapter/in/http/dto"
	"backend/internal/adapter/in/http/mapper"
	"backend/internal/app/usecase"
	pkgerrs "backend/pkg/errs"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type InventoryHandler struct {
	log               *slog.Logger
	getInventoryUC    *usecase.GetInventoryUC
	updateInventoryUC *usecase.UpdateInventoryUC
}

func NewInventoryHandler(
	log *slog.Logger,
	getInventoryUC *usecase.GetInventoryUC,
	updateInventoryUC *usecase.UpdateInventoryUC,
) *InventoryHandler {
	return &InventoryHandler{
		log:               log,
		getInventoryUC:    getInventoryUC,
		updateInventoryUC: updateInventoryUC,
	}
}

func (h *InventoryHandler) Get(c echo.Context) error {
	var req httpdto.GetInventoryRequest

	if err := c.Bind(&req); err != nil {
		return h.returnErr(c, "binding failed", pkgerrs.ErrInvalidJSON)
	}

	if _, err := uuid.Parse(req.LocationID); err != nil {
		return h.returnErr(c, "failed to parse uuid", pkgerrs.ErrInvalidIdentifier)
	}

	out, err := h.getInventoryUC.Execute(
		c.Request().Context(),
		mapper.MapRequestToGetInventory(req),
	)
	if err != nil {
		return h.returnErr(c, "failed to get inventory", err)
	}

	return c.JSON(http.StatusOK, mapper.MapOutputToGetInventory(out))
}

func (h *InventoryHandler) Update(c echo.Context) error {
	var req httpdto.UpdateInventoryRequest

	if err := c.Bind(&req); err != nil {
		return h.returnErr(c, "binding failed", pkgerrs.ErrInvalidJSON)
	}

	if _, err := uuid.Parse(req.LocationID); err != nil {
		return h.returnErr(c, "failed to parse uuid", pkgerrs.ErrInvalidIdentifier)
	}

	out, err := h.updateInventoryUC.Execute(
		c.Request().Context(),
		mapper.MapRequestToUpdateInventory(req),
	)
	if err != nil {
		return h.returnErr(c, "failed to update inventory", err)
	}

	h.log.InfoContext(
		c.Request().Context(), "inventory updated",
		slog.String("location_id", req.LocationID),
		slog.Int("items affected", len(out.Inventory)),
	)

	return c.JSON(http.StatusOK, mapper.MapOutputToUpdateInventory(out))
}

func (h *InventoryHandler) returnErr(c echo.Context, msg string, err error) error {
	if err == nil {
		return c.NoContent(http.StatusBadRequest)
	}

	outErr := mapper.HttpError(err)

	h.log.ErrorContext(c.Request().Context(), msg,
		slog.Int("code", outErr.Code),
		slog.String("public_msg", outErr.Message),
		slog.Any("cause", outErr.Reason),
	)

	return c.JSON(outErr.Code, mapper.MapErrorToResponse(outErr.Message))
}
