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

type ItemHandler struct {
	log          *slog.Logger
	createItemUC *usecase.CreateItemUC
	updateItemUC *usecase.UpdateItemUC
	deleteItemUC *usecase.DeleteItemUC
	listItemsUC  *usecase.ListItemsUC
}

func NewItemHandler(
	log *slog.Logger,
	createItemUC *usecase.CreateItemUC,
	updateItemUC *usecase.UpdateItemUC,
	deleteItemUC *usecase.DeleteItemUC,
	listItemsUC *usecase.ListItemsUC,
) *ItemHandler {
	return &ItemHandler{
		log:          log,
		createItemUC: createItemUC,
		updateItemUC: updateItemUC,
		deleteItemUC: deleteItemUC,
		listItemsUC:  listItemsUC,
	}
}

func (h *ItemHandler) Create(c echo.Context) error {
	var req httpdto.CreateItemRequest

	if err := c.Bind(&req); err != nil {
		return h.returnErr(c, "binding failed", pkgerrs.ErrInvalidJSON)
	}

	out, err := h.createItemUC.Execute(
		c.Request().Context(),
		mapper.MapRequestToCreateItem(req),
	)
	if err != nil {
		return h.returnErr(c, "failed to create item", err)
	}

	h.log.InfoContext(
		c.Request().Context(), "item created",
		slog.Any("name", req.Name),
		slog.Any("description", req.Description),
		slog.Any("category", req.Category),
		slog.Any("photo_url", req.PhotoURL),
		slog.Any("nutrition", req.Nutrition),
	)

	return c.JSON(http.StatusCreated, mapper.MapOutputToCreateItem(out))
}

func (h *ItemHandler) Update(c echo.Context) error {
	var req httpdto.UpdateItemRequest

	if err := c.Bind(req); err != nil {
		return h.returnErr(c, "binding failed", pkgerrs.ErrInvalidJSON)
	}

	if _, err := uuid.Parse(req.ID); err != nil {
		return h.returnErr(c, "failed to parse uuid", pkgerrs.ErrInvalidIdentifier)
	}

	out, err := h.updateItemUC.Execute(
		c.Request().Context(),
		mapper.MapRequestToUpdateItem(req),
	)

	if err != nil {
		return h.returnErr(c, "failed to update item", err)
	}

	h.log.InfoContext(
		c.Request().Context(), "item updated",
		slog.String("id", req.ID),
		slog.Any("name", req.Name),
		slog.Any("description", req.Description),
		slog.Any("category", req.Category),
		slog.Any("photo_url", req.PhotoURL),
		slog.Any("nutrition", req.Nutrition),
	)

	return c.JSON(http.StatusOK, mapper.MapOutputToUpdateItem(out))
}

func (h *ItemHandler) Delete(c echo.Context) error {
	var req httpdto.DeleteItemRequest

	err := c.Bind(&req)
	if err != nil {
		return h.returnErr(c, "binding failed", pkgerrs.ErrInvalidJSON)
	}

	if _, err = uuid.Parse(req.ID); err != nil {
		return h.returnErr(c, "failed to parse uuid", pkgerrs.ErrInvalidIdentifier)
	}

	err = h.deleteItemUC.Execute(
		c.Request().Context(),
		mapper.MapRequestToDeleteItem(req),
	)
	if err != nil {
		return h.returnErr(c, "failed to delete item", err)
	}

	h.log.InfoContext(
		c.Request().Context(), "item deleted",
		slog.String("id", req.ID),
	)

	return c.NoContent(http.StatusNoContent)
}

func (h *ItemHandler) List(c echo.Context) error {
	out, err := h.listItemsUC.Execute(c.Request().Context())
	if err != nil {
		return h.returnErr(c, "failed to get a list of items", err)
	}

	return c.JSON(http.StatusOK, mapper.MapOutputToListItems(out))
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
