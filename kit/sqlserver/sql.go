package sqlserver

import (
	"database/sql"
	"fmt"

	_ "github.com/denisenkom/go-mssqldb"
)

type SqlServer struct {
	db *sql.DB
}

func NewSqlServer(config DBConfig) (*SqlServer, error) {
	if config.SSLMode == "" {
		config.SSLMode = "disable"
	}

	dns := fmt.Sprintf("server=%s; user id=%s; password=%s; port=%d; database=%s;",
		config.Server,
		config.User,
		config.Password,
		config.Port,
		config.DBName,
	)

	db, err := sql.Open("sqlserver", dns)
	return &SqlServer{db: db}, err
}

func (p *SqlServer) GetConnection() *sql.DB {
	return p.db
}
