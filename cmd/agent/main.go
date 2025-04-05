package main

import (
	"log"
	"time"

	"github.com/sincin-v/collector/internal/agent/clients/rest"
	"github.com/sincin-v/collector/internal/agent/config"
	"github.com/sincin-v/collector/internal/agent/helpers/metrics"
	"github.com/sincin-v/collector/internal/common/service"
	"github.com/sincin-v/collector/internal/common/storage"
)

func main() {
	log.Printf("Start agent work")
	agentConfig := config.GetAgentConfig()

	log.Printf("Send metrics to %s", agentConfig.ServerHost)
	metricStorage := storage.New()
	service := service.New(&metricStorage)
	hc := rest.New(agentConfig.ServerHost)
	metricsCollector := metrics.New(&service, hc)
	go metricsCollector.StartSendMetrics(agentConfig.ReportInterval)
	for {
		// go metricsCollector.StartCollectMetrics(pollInterval)
		go metricsCollector.CollectMetrics()

		time.Sleep(agentConfig.PollInterval)
	}
}
