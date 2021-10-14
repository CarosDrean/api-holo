package service

import (
	"database/sql"
	"fmt"

	"api-holo/domain/service"
	"api-holo/infraestucture/handler/response"
	serviceStorage "api-holo/infraestucture/sqlserver/service"
	"api-holo/kit/authorization"
	"api-holo/model"

	"github.com/labstack/echo/v4"
)

const routeVersionPrefix = "api/v1/service"

func NewRouter(app *echo.Echo, db *sql.DB, authMiddleware authorization.AuthMiddleware, logger model.Logger) {
	useCaseCie10 := service.New(serviceStorage.New(db))
	responser := response.New(logger)

	handler := NewHandler(useCaseCie10, responser)

	adminRoutes(app, handler, authMiddleware.ValidateJWT)
	privateRoutes(app, handler, authMiddleware.ValidateJWT)
}

func adminRoutes(app *echo.Echo, handler Handler, middlewares ...echo.MiddlewareFunc) {
	api := app.Group(fmt.Sprintf("%sadmin/", routeVersionPrefix), middlewares...)

	api.POST("", handler.Create)
	api.PUT("/:id", handler.Update)
	api.GET("", handler.GetWhere)
	api.GET("all", handler.GetAllWhere)
}

func privateRoutes(app *echo.Echo, handler Handler, middlewares ...echo.MiddlewareFunc) {
	api := app.Group(fmt.Sprintf("%sprivate/", routeVersionPrefix), middlewares...)

	api.GET("", handler.GetWhere)
	api.GET("all", handler.GetAllWhere)
}
