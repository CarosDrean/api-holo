package login

import (
	"api-holo/domain/login"
	"api-holo/domain/user"
	"api-holo/infraestucture/handler/response"
	userStorage "api-holo/infraestucture/sqlserver/user"
	"api-holo/model"
	"crypto/rsa"
	"database/sql"

	"github.com/labstack/echo/v4"
)

const routeVersionPrefix = "api/v1/login"

func NewRouter(app *echo.Echo, db *sql.DB, logger model.Logger, privateKey *rsa.PrivateKey) {
	useCaseUser := user.New(userStorage.New(db))
	useCase := login.New(useCaseUser, privateKey)
	responser := response.New(logger)

	handler := NewHandler(useCase, responser)

	routes(app, handler)
}

func routes(app *echo.Echo, handler Handler) {
	api := app.Group(routeVersionPrefix)

	api.POST("", handler.ValidateLogin)
}
