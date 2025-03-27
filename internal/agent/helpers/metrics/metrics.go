package metrics

import (
	"fmt"
	"log"
	"math/rand"
	"reflect"
	"runtime"

	"github.com/sincin-v/collector/internal/storage"
)

var PollCountValue int = 0

var collectMetricsNames = [...]string{
	"Alloc",
	"TotalAlloc",
	"Sys",
	"Lookups",
	"Mallocs",
	"Frees",
	"HeapAlloc",
	"HeapSys",
	"HeapIdle",
	"HeapInuse",
	"HeapReleased",
	"HeapObjects",
	"StackInuse",
	"StackSys",
	"MSpanInuse",
	"MSpanSys",
	"MCacheInuse",
	"MCacheSys",
	"BuckHashSys",
	"GCSys",
	"OtherSys",
	"NextGC",
	"LastGC",
	"PauseTotalNs",
	"NumGC",
	"NumForcedGC",
	"GCCPUFraction",
}

var metricsStorage = storage.MetricStorage{Metrics: make(map[string]storage.Metric)}

func GetMetrics(ch chan<- storage.MetricStorage) {

	log.Printf("Start collect metrics")

	var metrics runtime.MemStats
	runtime.ReadMemStats(&metrics)
	metrics_values := reflect.ValueOf(metrics)

	for _, metricName := range collectMetricsNames {
		mv := reflect.Indirect(metrics_values).FieldByName(metricName)
		var metric = metricsStorage.GetMetrics(metricName)
		if metric == nil {
			metric = &storage.GaugeMetric{Name: metricName}
		}
		storage.SetMetricValue(metric, fmt.Sprintf("%v", mv), `gauge`)
		metricsStorage.CreateMetric(fmt.Sprintf("%s", metricName), metric)
	}
	pollCountMetric := metricsStorage.GetMetrics("PollCount")
	if pollCountMetric == nil {
		pollCountMetric = &storage.CounterMetric{Name: "PollCount"}
	}
	storage.SetMetricValue(pollCountMetric, `1`, `counter`)
	metricsStorage.CreateMetric("PollCount", pollCountMetric)

	randomValueMetric := metricsStorage.GetMetrics("randomValue")
	if randomValueMetric == nil {
		randomValueMetric = &storage.GaugeMetric{Name: "randomValue"}
	}
	storage.SetMetricValue(randomValueMetric, fmt.Sprintf("%f", rand.Float64()), `gauge`)
	metricsStorage.CreateMetric("randomValue", randomValueMetric)

	ch <- metricsStorage
}
