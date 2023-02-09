package handler

import (
	"net/http"

	"github.com/labstack/echo"

	http_helper "github.com/cityos-dev/Cornelius-David-Herianto/helper/http"
	"github.com/cityos-dev/Cornelius-David-Herianto/internal/health/service"
)

type healthHTTPHandler struct {
	healthSvc service.Service
}

func New(healthSvc service.Service) healthHTTPHandler {
	return healthHTTPHandler{
		healthSvc: healthSvc,
	}
}

func (h *healthHTTPHandler) GetHealth(ctx echo.Context) error {
	if err := h.healthSvc.GetServiceHealth(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, http_helper.NewHTTPErrorMessage(err))
	}

	return ctx.String(http.StatusOK, "OK")
}
