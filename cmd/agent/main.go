package main

import (
	"flag"
	"log"
	"time"

	"github.com/sincin-v/collector/internal/agent/clients/rest"
	"github.com/sincin-v/collector/internal/agent/helpers/metrics"
	"github.com/sincin-v/collector/internal/storage"
)

func main() {
	log.Printf("Start agent work")

	serverHost := flag.String("a", "localhost:8080", "Metric server host and port")
	reportInterval := flag.Int("r", 10, "Report interval")
	pollInterval := flag.Int("p", 2, "Poll interval")

	flag.Parse()

	log.Printf("Send metrics to %s", *serverHost)
	var metricCh = make(chan storage.MetricStorage)
	go rest.SendMetric(metricCh, *serverHost, *reportInterval)
	for {
		go metrics.GetMetrics(metricCh)
		time.Sleep(time.Duration(*pollInterval) * time.Second)
	}
}
