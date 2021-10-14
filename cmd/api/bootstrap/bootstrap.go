package bootstrap

import (
	"flag"
	"fmt"

	"api-holo/infraestucture/handler/response"
	"api-holo/infraestucture/handler/router"
)

func Run() error {
	pathConfig := flag.String("config", "configuration.json", "Configuration file location. You must include the file name: Ex: /tu/path/configuration.json")
	flag.Parse()

	config := newConfiguration(*pathConfig)
	logger := newLogrus(config.LogFolder, false)
	db := newSQLDatabase(config)
	loadSignatures(config, logger)
	api := newEcho(config, response.HTTPErrorHandler)

	router.InitRoutes(api, db, logger)

	port := fmt.Sprintf(":%d", config.PortHttp)
	return api.Start(port)

}
