package http

import (
	httpdto "backend/internal/adapter/in/http/dto"
	"backend/internal/adapter/in/http/mapper"
	"backend/internal/app/usecase"
	"context"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ItemHandler struct {
	log          *slog.Logger
	createItemUC *usecase.CreateItemUC
	deleteItemUC *usecase.DeleteItemUC
}

func NewItemHandler(
	log *slog.Logger,
	createItemUC *usecase.CreateItemUC,
	deleteItemUC *usecase.DeleteItemUC,
) *ItemHandler {
	return &ItemHandler{
		log:          log,
		createItemUC: createItemUC,
		deleteItemUC: deleteItemUC,
	}
}

func (h *ItemHandler) Create(c echo.Context) error {
	var req httpdto.CreateItemRequest

	if err := c.Bind(&req); err != nil {
		return h.returnErr(c, "invalid json", err)
	}

	out, err := h.createItemUC.Execute(
		c.Request().Context(),
		mapper.MapRequestToCreateItem(req),
	)
	if err != nil {
		return h.returnErr(c, "failed to create item", err)
	}

	return c.JSON(http.StatusCreated, mapper.MapOutputToCreateItem(out))
}

func (h *ItemHandler) Delete(c echo.Context) error {
	var req httpdto.DeleteItemRequest

	err := c.Bind(&req)
	if err != nil {
		return h.returnErr(c, "binding failed", err)
	}

	if _, err = uuid.Parse(req.ID); err != nil {
		return h.returnErr(c, "invalid uuid", echo.NewHTTPError(http.StatusBadRequest, "invalid identifier format"))
	}

	err = h.deleteItemUC.Execute(
		c.Request().Context(),
		mapper.MapRequestToDeleteItem(req),
	)
	if err != nil {
		return h.returnErr(c, "failed to delete item", err)
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *ItemHandler) returnErr(c echo.Context, msg string, err error) error {
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
