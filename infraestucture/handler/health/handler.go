package health

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type Health struct{}

func NewHandler() Health {
	return Health{}
}

func (h *Health) Health(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"hora": time.Now().String(), "respuesta": "Hello World, scripts!"})
}
