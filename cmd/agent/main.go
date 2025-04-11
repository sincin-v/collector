package main

import (
	"log"
	"time"

	"github.com/sincin-v/collector/internal/agent/clients/rest"
	"github.com/sincin-v/collector/internal/agent/config"
	"github.com/sincin-v/collector/internal/agent/helpers/metrics"
	"github.com/sincin-v/collector/internal/service"
	"github.com/sincin-v/collector/internal/storage"
)

func main() {
	log.Printf("Start agent work")
	agentConfig, err := config.GetAgentConfig()
	if err != nil {
		log.Fatalf("Cannot get agent params for start")
	}

	log.Printf("Send metrics to %s", agentConfig.ServerHost)
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
