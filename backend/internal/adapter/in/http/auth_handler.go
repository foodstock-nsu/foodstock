package http

import (
	httpdto "backend/internal/adapter/in/http/dto"
	"backend/internal/adapter/in/http/mapper"
	"backend/internal/app/usecase"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	log         *slog.Logger
	adminAuthUC *usecase.AdminAuthUC
}

func NewAuthHandler(
	log *slog.Logger,
	adminAuthUC *usecase.AdminAuthUC,
) *AuthHandler {
	return &AuthHandler{
		log:         log,
		adminAuthUC: adminAuthUC,
	}
}

func (h *AuthHandler) AdminAuth(c echo.Context) error {
	var req httpdto.AdminAuthRequest
	if err := c.Bind(&req); err != nil {
		return h.returnErr(c, "invalid json", err)
	}

	out, err := h.adminAuthUC.Execute(
		c.Request().Context(),
		mapper.MapRequestToAdminAuth(req),
	)
	if err != nil {
		return h.returnErr(c, "failed to auth admin", err)
	}

	return c.JSON(http.StatusOK, mapper.MapOutputToAdminAuth(out))
}

func (h *AuthHandler) returnErr(c echo.Context, msg string, err error) error {
	outErr := mapper.HttpError(err)

	h.log.ErrorContext(c.Request().Context(), msg,
		slog.Int("code", outErr.Code),
		slog.String("public_msg", outErr.Message),
		slog.Any("cause", outErr.Reason),
	)

	return c.JSON(outErr.Code, map[string]string{"error": msg})
}
