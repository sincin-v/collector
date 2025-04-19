package main

import (
	"time"

	"github.com/sincin-v/collector/internal/agent/clients/rest"
	"github.com/sincin-v/collector/internal/agent/config"
	"github.com/sincin-v/collector/internal/agent/helpers/metrics"
	"github.com/sincin-v/collector/internal/logger"
	"github.com/sincin-v/collector/internal/service"
	"github.com/sincin-v/collector/internal/storage"
)

func main() {

	agentConfig, err := config.GetAgentConfig()
	if err != nil {
		panic("Cannot get agent params for start")
	}
	logErr := logger.Initialize(agentConfig.LogLevel)
	if logErr != nil {
		panic(logErr)
	}
	logger.Log.Info("Start agent work")
	logger.Log.Info("Send metrics to %s", agentConfig.ServerHost)
	memStorage := storage.New()
	service := service.New(&memStorage)
	hc := rest.New(agentConfig.ServerHost)
	metricsCollector := metrics.New(&service, hc)
	go metricsCollector.StartSendMetrics(agentConfig.ReportInterval)
	for {
		go metricsCollector.CollectMetrics()

		time.Sleep(agentConfig.PollInterval)
	}
}
