package sqlserver

type DBConfig struct {
	Driver   string
	User     string
	Password string
	Server   string
	Port     uint
	DBName   string
	SSLMode  string
}

func NewDBConfig(driver string, user string, password string, server string, port uint, DBName string, SSLMode string) DBConfig {
	return DBConfig{Driver: driver, User: user, Password: password, Server: server, Port: port, DBName: DBName, SSLMode: SSLMode}
}
