package http

import (
	httpdto "backend/internal/adapter/in/http/dto"
	"backend/internal/adapter/in/http/mapper"
	"backend/internal/app/usecase"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ClientHandler struct {
	log          *slog.Logger
	getCatalogUC *usecase.GetCatalogUC
}

func NewClientHandler(
	log *slog.Logger,
	getCatalogUC *usecase.GetCatalogUC,
) *ClientHandler {
	return &ClientHandler{
		log:          log,
		getCatalogUC: getCatalogUC,
	}
}

func (h *ClientHandler) GetCatalog(c echo.Context) error {
	var req httpdto.GetCatalogRequest
	if err := c.Bind(&req); err != nil {
		return h.returnErr(c, "invalid query", err)
	}

	out, err := h.getCatalogUC.Execute(
		c.Request().Context(),
		mapper.MapRequestToGetCatalog(req),
	)
	if err != nil {
		return h.returnErr(c, "failed to get catalog", err)
	}

	return c.JSON(http.StatusOK, mapper.MapOutputToGetCatalog(out))
}

func (h *ClientHandler) returnErr(c echo.Context, msg string, err error) error {
	outErr := mapper.HttpError(err)

	h.log.ErrorContext(c.Request().Context(), msg,
		slog.Int("code", outErr.Code),
		slog.String("public_msg", outErr.Message),
		slog.Any("cause", outErr.Reason),
	)

	return c.JSON(outErr.Code, mapper.MapErrorToResponse(outErr.Message))
}
