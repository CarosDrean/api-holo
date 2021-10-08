package bootstrap

import (
	"database/sql"
	"log"

	"api-holo/kit/sqlserver"
)

func newSQLDatabase(conf Configuration) *sql.DB {
	psqlDB, err := sqlserver.NewSqlServer(
		sqlserver.NewDBConfig(
			conf.Database.Engine,
			conf.Database.User,
			conf.Database.Password,
			conf.Database.Server,
			conf.Database.Port,
			conf.Database.Name,
			conf.Database.SSLMode,
		),
	)
	if err != nil {
		log.Fatalf("no se pudo obtener una conexion a la base de datos de extraccion: %v", err)
	}

	return psqlDB.GetConnection()
}
