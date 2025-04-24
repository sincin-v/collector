package main

import (
	"net/http"

	"github.com/sincin-v/collector/internal/logger"
	"github.com/sincin-v/collector/internal/server/collector"
	"github.com/sincin-v/collector/internal/server/config"
	"github.com/sincin-v/collector/internal/server/router"
	"github.com/sincin-v/collector/internal/service"
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
	metricService := service.New(&memStorage)
	metricCollector := collector.New(metricService, serverConfig.FileStoragePath)

	if serverConfig.Restore {
		err := metricCollector.RestoreMetrics()
		if err != nil {
			logger.Log.Warnf("Cannot restore metrics from %s", serverConfig.FileStoragePath)
		}
	}

	var errSaveMetric error

	go func() {
		errSaveMetric = metricCollector.SaveMetrics(int(serverConfig.StoreInterval))
		if errSaveMetric != nil {
			logger.Log.Error("Error save metric: %s", errSaveMetric)
		}
	}()

	serverRouter := router.CreateRouter(&metricService)

	httpErr := http.ListenAndServe(serverConfig.Host, serverRouter)
	if httpErr != nil {
		panic(httpErr)
	}

}
