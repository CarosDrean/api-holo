package health

import "github.com/labstack/echo/v4"

func NewRouter(api *echo.Echo) {
	handle := NewHandler()
	api.GET("/health", handle.Health)
}
