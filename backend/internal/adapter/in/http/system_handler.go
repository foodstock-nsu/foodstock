package http

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

var startTime = time.Now()

type SystemHandler struct {
	env        string
	apiVersion string
}

func NewSystemHandler(env, apiVersion string) *SystemHandler {
	return &SystemHandler{env: env, apiVersion: apiVersion}
}

func (h *SystemHandler) Health(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}

func (h *SystemHandler) Info(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"app":         "foodstock-api",
		"environment": h.env,
		"uptime":      time.Since(startTime).String(),
		"api_version": h.apiVersion,
	})
}
