package metrics

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"reflect"
	"runtime"
	"time"
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

type MetricsService interface {
	CreateMetric(string, string, string) (string, error)
	GetMetric(string, string) (string, error)
	GetAllMetrics() (map[string]int64, map[string]float64)
}

type HttpClient interface {
	SendPostRequest(string) (*http.Response, error)
}

type Collector struct {
	service    MetricsService
	httpClient HttpClient
}

func New(s MetricsService, hc HttpClient) Collector {
	return Collector{service: s, httpClient: hc}
}

func (c Collector) StartCollectMetrics(pollInterval time.Duration) {
	for {
		c.CollectMetrics()
		time.Sleep(pollInterval)
	}
}

func (c Collector) StartSendMetrics(reportInterval time.Duration) {
	for {
		time.Sleep(reportInterval)
		c.SendMetrics()
	}
}

func (c Collector) CollectMetrics() {

	log.Printf("Start collect metrics")

	var metrics runtime.MemStats
	runtime.ReadMemStats(&metrics)
	metricsValues := reflect.ValueOf(metrics)

	for _, metricName := range collectMetricsNames {
		mv := reflect.Indirect(metricsValues).FieldByName(metricName)
		c.service.CreateMetric("gauge", metricName, fmt.Sprintf("%v", mv))
	}
	c.service.CreateMetric("counter", "PollCount", "1")
	c.service.CreateMetric("gauge", "randomValue", fmt.Sprintf("%f", rand.Float64()))

	log.Printf("Finish collect metrics")
}

func (c Collector) SendMetrics() {
	log.Printf("Send metric")
	counterMetrics, gaugeMetrics := c.service.GetAllMetrics()
	for metricName := range gaugeMetrics {
		metricValue := gaugeMetrics[metricName]

		methodURL := fmt.Sprintf("/update/gauge/%s/%f", metricName, metricValue)
		_, err := c.httpClient.SendPostRequest(methodURL)
		if err != nil {
			log.Fatalf("Cannot send request to server to set metric %s", metricName)
			continue
		}
	}
	for metricName := range counterMetrics {
		metricValue := counterMetrics[metricName]

		methodURL := fmt.Sprintf("/update/counter/%s/%d", metricName, metricValue)
		_, err := c.httpClient.SendPostRequest(methodURL)
		if err != nil {
			log.Fatalf("Cannot send request to server to set metric %s", metricName)
			continue
		}
	}
	log.Printf("Finish send metrics")

}
