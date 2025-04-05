package main

import (
	"flag"
	"log"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/sincin-v/collector/internal/agent/clients/rest"
	"github.com/sincin-v/collector/internal/agent/config"
	"github.com/sincin-v/collector/internal/agent/helpers/metrics"
	"github.com/sincin-v/collector/internal/common/service"
	"github.com/sincin-v/collector/internal/common/storage"
)

func main() {
	log.Printf("Start agent work")

	var argServerHost = flag.String("a", "localhost:8080", "Metric server host and port")
	var argReportInterval = flag.Int("r", 10, "Report interval")
	var argPollInterval = flag.Int("p", 2, "Poll interval")

	flag.Parse()

	var cfg config.Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Printf("Env is empty")
	}

	var serverHost string
	if cfg.Address != "" {
		serverHost = cfg.Address
	} else {
		serverHost = *argServerHost
	}
	var reportInterval time.Duration
	if cfg.ReportInterval != 0 {
		reportInterval = cfg.ReportInterval
	} else {
		reportInterval = time.Duration(time.Duration(*argReportInterval) * time.Second)
	}

	var pollInterval time.Duration
	if cfg.PollInterval != 0 {
		pollInterval = cfg.PollInterval
	} else {
		pollInterval = time.Duration(time.Duration(*argPollInterval) * time.Second)

	}

	log.Printf("Send metrics to %s", serverHost)
	metricStorage := storage.New()
	service := service.New(&metricStorage)
	hc := rest.New(serverHost)
	metricsCollector := metrics.New(&service, hc)
	go metricsCollector.StartSendMetrics(reportInterval)
	for {
		// go metricsCollector.StartCollectMetrics(pollInterval)
		go metricsCollector.CollectMetrics()

		time.Sleep(pollInterval)
	}
}
