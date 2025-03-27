package main

import (
	"log"
	"time"

	"github.com/sincin-v/collector/internal/agent/clients/rest"
	"github.com/sincin-v/collector/internal/storage"
	"github.com/sincin-v/collector/internal/agent/helpers/metrics"
)

var pollInterval int64 = 2

func main() {
	log.Printf("Start agent work")
	metricCh := make(chan storage.MetricStorage)
	go rest.SendMetric(metricCh)
	for {
		go metrics.GetMetrics(metricCh)
		time.Sleep(time.Duration(pollInterval) * time.Second)
	}
}
