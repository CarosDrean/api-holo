package router

import (
	"api-holo/infraestucture/handler/cie10"
	"api-holo/kit/authorization"
	"api-holo/model"
	"database/sql"

	"api-holo/infraestucture/handler/health"

	"github.com/labstack/echo/v4"
)

func InitRoutes(app *echo.Echo, db *sql.DB, logger model.Logger) {
	authMiddleware := authorization.NewAuthServiceValidator(logger)

	// H
	health.NewRouter(app)

	cie10.NewRouter(app, db, authMiddleware, logger)
}
