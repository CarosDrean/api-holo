package login

import (
	"crypto/rsa"
	"database/sql"
	"fmt"

	"api-holo/domain/login"
	"api-holo/domain/user"
	"api-holo/infraestucture/handler/response"
	userStorage "api-holo/infraestucture/sqlserver/user"
	"api-holo/kit/authorization"
	"api-holo/model"

	"github.com/labstack/echo/v4"
)

const routeVersionPrefix = "api/v1/login"

func NewRouter(app *echo.Echo, db *sql.DB, authMiddleware authorization.AuthMiddleware, logger model.Logger, privateKey *rsa.PrivateKey) {
	useCaseUser := user.New(userStorage.New(db))
	useCase := login.New(useCaseUser, privateKey)
	responser := response.New(logger)

	handler := NewHandler(useCase, responser)

	adminRoutes(app, handler, authMiddleware.ValidateJWT)
	privateRoutes(app, handler, authMiddleware.ValidateJWT)
}

func adminRoutes(app *echo.Echo, handler Handler, middlewares ...echo.MiddlewareFunc) {
	api := app.Group(fmt.Sprintf("%sadmin/", routeVersionPrefix), middlewares...)

	api.POST("", handler.ValidateLogin)
}

func privateRoutes(app *echo.Echo, handler Handler, middlewares ...echo.MiddlewareFunc) {
	api := app.Group(fmt.Sprintf("%sprivate/", routeVersionPrefix), middlewares...)

	api.POST("", handler.ValidateLogin)
}
