package cie10

import (
	"database/sql"
	"fmt"

	"api-holo/domain/cie10"
	"api-holo/infraestucture/handler/response"
	cie10Storage "api-holo/infraestucture/sqlserver/cie10"
	"api-holo/kit/authorization"
	"api-holo/model"

	"github.com/labstack/echo/v4"
)

const routeVersionPrefix = "api/v1/"

func NewRouter(app *echo.Echo, db *sql.DB, authMiddleware authorization.AuthMiddleware, logger model.Logger) {
	useCaseCie10 := cie10.New(cie10Storage.New(db), logger)
	responser := response.New(logger)

	handler := NewHandler(useCaseCie10, responser)

	adminRoutes(app, handler, authMiddleware.ValidateJWT)
	privateRoutes(app, handler, authMiddleware.ValidateJWT)
}

func adminRoutes(app *echo.Echo, handler Handler, middlewares ...echo.MiddlewareFunc) {
	api := app.Group(fmt.Sprintf("%sadmin/cie10/", routeVersionPrefix), middlewares...)

	api.POST("", handler.Create)
	api.PUT("/:id", handler.Update)
	api.GET("all", handler.GetAllWhere)
	api.GET("", handler.GetWhere)
}

func privateRoutes(app *echo.Echo, handler Handler, middlewares ...echo.MiddlewareFunc) {
	api := app.Group(fmt.Sprintf("%sprivate/cie10/", routeVersionPrefix), middlewares...)

	api.GET("all", handler.GetAllWhere)
	api.GET("", handler.GetWhere)
}
