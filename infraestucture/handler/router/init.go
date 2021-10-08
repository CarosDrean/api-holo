package router

import (
	"api-holo/model"
	"database/sql"

	"api-holo/infraestucture/handler/health"

	"github.com/labstack/echo/v4"
)

type Configuration interface {
	DBEngine() string
}

func InitRoutes(app *echo.Echo, db *sql.DB, config Configuration, logger model.Logger) {
	//authMiddleware := authorization.NewAuthServiceValidator(logger)

	// H
	health.NewRouter(app)

	//executor.NewRouter(app, authMiddleware, permissionMiddleware, logger, param, queueUseCase)
}
