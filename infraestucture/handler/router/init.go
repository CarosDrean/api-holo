package router

import (
	"crypto/rsa"
	"database/sql"

	"api-holo/infraestucture/handler/cie10"
	"api-holo/infraestucture/handler/health"
	"api-holo/infraestucture/handler/login"
	"api-holo/infraestucture/handler/service"
	"api-holo/kit/authorization"
	"api-holo/model"

	"github.com/labstack/echo/v4"
)

func InitRoutes(app *echo.Echo, db *sql.DB, logger model.Logger, privateKey *rsa.PrivateKey) {
	authMiddleware := authorization.NewAuthServiceValidator(logger)

	cie10.NewRouter(app, db, authMiddleware, logger)

	// H
	health.NewRouter(app)

	// L
	login.NewRouter(app, db, logger, privateKey)

	// S
	service.NewRouter(app, db, authMiddleware, logger)
}
