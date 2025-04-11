package main

import (
	"net/http"

	"github.com/sincin-v/collector/internal/logger"
	"github.com/sincin-v/collector/internal/server/config"
	"github.com/sincin-v/collector/internal/server/router"
	"github.com/sincin-v/collector/internal/storage"
)

func main() {

	serverConfig, cfgErr := config.GetServerConfig()
	if cfgErr != nil {
		panic("Cannot get server params for start")
	}

	logErr := logger.Initialize(serverConfig.LogLevel)
	if logErr != nil {
		panic(logErr)
	}

	logger.Log.Infof("Start server work on %s", serverConfig.Host)

	memStorage := storage.New()

	serverRouter := router.CreateRouter(&memStorage)

	httpErr := http.ListenAndServe(serverConfig.Host, serverRouter)
	if httpErr != nil {
		panic(httpErr)
	}

}
