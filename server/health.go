package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func HealthCheck(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}
