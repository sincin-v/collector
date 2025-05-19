package main

import (
	"context"
	"net/http"

	"github.com/sincin-v/collector/internal/database/migrator"
	"github.com/sincin-v/collector/internal/logger"
	"github.com/sincin-v/collector/internal/server/clients/db"
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

	var baseCtx = context.Background()

	var metricService service.MetricsService
	var dbClient *db.DBClient
	var err error

	if serverConfig.DBDns != "" {
		dbClient, err = db.New(serverConfig.DBDns, serverConfig.RetryIntervals)
		if err != nil {
			logger.Log.Panic("Error connect to DB %s Error: %s", serverConfig.DBDns, err)
		}

		migrateErr := migrator.ApplyMigrations(serverConfig.DBDns, serverConfig.MigrationPath)
		if migrateErr != nil {
			logger.Log.Panic(migrateErr)
		}
		storage := storage.NewDBStorage(baseCtx, *dbClient)
		metricService = service.New(storage)

		defer dbClient.Close()
	} else {
		storage := storage.NewMemStorage()
		metricService = service.New(&storage)

		metricCollector := collector.New(metricService, serverConfig.FileStoragePath)

		if serverConfig.Restore {
			err := metricCollector.RestoreMetrics(baseCtx)
			if err != nil {
				logger.Log.Warnf("Cannot restore metrics from %s", serverConfig.FileStoragePath)
			}
		}

		var errSaveMetric error

		go func() {
			errSaveMetric = metricCollector.SaveMetrics(baseCtx, int(serverConfig.StoreInterval))
			if errSaveMetric != nil {
				logger.Log.Error("Error save metric: %s", errSaveMetric)
			}
		}()
	}

	serverRouter, errCreateRouter := router.CreateRouter(&metricService, *serverConfig)

	if errCreateRouter != nil {
		panic(errCreateRouter)
	}

	httpErr := http.ListenAndServe(serverConfig.Host, serverRouter)
	if httpErr != nil {
		panic(httpErr)
	}

}
