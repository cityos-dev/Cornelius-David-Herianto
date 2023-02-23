package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	httpHelper "github.com/cityos-dev/Cornelius-David-Herianto/helper/http"
	"github.com/cityos-dev/Cornelius-David-Herianto/internal/health/service"
)

type healthHTTPHandler struct {
	healthSvc service.Service
}

// New returns new instance of healthHTTPHandler
func New(healthSvc service.Service) healthHTTPHandler {
	return healthHTTPHandler{
		healthSvc: healthSvc,
	}
}

// GetHealth handles HTTP request for getting the service's health
func (h *healthHTTPHandler) GetHealth(ctx echo.Context) error {
	if err := h.healthSvc.GetServiceHealth(ctx.Request().Context()); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, httpHelper.NewErrorMessage("service is not healthy", err))
	}

	return ctx.String(http.StatusOK, "OK")
}
