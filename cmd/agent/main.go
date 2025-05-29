package main

import (
	"context"
	"sync"
	"time"

	"github.com/sincin-v/collector/internal/agent/clients/rest"
	"github.com/sincin-v/collector/internal/agent/config"
	"github.com/sincin-v/collector/internal/agent/helpers/metrics"
	"github.com/sincin-v/collector/internal/logger"
	"github.com/sincin-v/collector/internal/service"
	"github.com/sincin-v/collector/internal/storage"
)

func worker(ctx context.Context, tasksChan chan int, wg *sync.WaitGroup, mc metrics.Collector) {
	defer wg.Done()
	for task := range tasksChan {
		logger.Log.Info("Send metrics from task %d", task)
		mc.SendMetricsJSON(ctx)
	}
}

func main() {

	agentConfig, err := config.GetAgentConfig()
	if err != nil {
		panic("Cannot get agent params for start")
	}
	logErr := logger.Initialize(agentConfig.LogLevel)
	if logErr != nil {
		panic(logErr)
	}
	ctx := context.Background()
	logger.Log.Info("Start agent work")
	logger.Log.Info("Send metrics to %s", agentConfig.ServerHost)
	memStorage := storage.NewMemStorage()
	service := service.New(&memStorage)
	hc := rest.New(agentConfig.ServerHost, agentConfig.RetryIntervals, agentConfig.SecretKey)
	metricsCollector := metrics.New(&service, hc)

	pollTiker := time.NewTicker(agentConfig.PollInterval)
	reportTiker := time.NewTicker(agentConfig.ReportInterval)
	if agentConfig.RateLimit == 0 {
		for {
			select {
			case <-pollTiker.C:
				metricsCollector.CollectMetrics(ctx)
			case <-reportTiker.C:
				metricsCollector.SendMetricsJSON(ctx)
			}
		}
	} else {
		for {
			tasksChan := make(chan int, int(agentConfig.RateLimit))
			var wg sync.WaitGroup
			wg.Add(2)

			for i := 0; i < int(agentConfig.RateLimit); i++ {

				go worker(ctx, tasksChan, &wg, metricsCollector)
			}

			go func() {
				defer wg.Done()
				for {
					metricsCollector.CollectMetrics(ctx)
					metricsCollector.CollectUtilizationMetric(ctx)
					tasksChan <- 1

					time.Sleep(agentConfig.PollInterval)
				}
			}()

			wg.Wait()
		}
	}
}
