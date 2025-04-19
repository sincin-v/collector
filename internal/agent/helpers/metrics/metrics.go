package metrics

import (
	"bytes"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"time"

	"github.com/sincin-v/collector/internal/compress"
	"github.com/sincin-v/collector/internal/models"
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

type MemMetrics struct {
}

type MetricsService interface {
	CreateCounterMetric(string, int64)
	CreateGaugeMetric(string, float64)
	GetMetric(string, string) (string, error)
	GetAllMetrics() (map[string]int64, map[string]float64)
}

type HTTPClient interface {
	SendPostRequest(string, bytes.Buffer) (*http.Response, error)
}

type Collector struct {
	service        MetricsService
	httpClient     HTTPClient
	memStatsMetric map[string]float64
}

func New(s MetricsService, hc HTTPClient) Collector {
	return Collector{service: s, httpClient: hc, memStatsMetric: make(map[string]float64)}
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

func (c *Collector) GetMetricsFromMemStats() {
	var metrics runtime.MemStats
	runtime.ReadMemStats(&metrics)
	c.memStatsMetric = map[string]float64{
		"Alloc":         float64(metrics.Alloc),
		"TotalAlloc":    float64(metrics.TotalAlloc),
		"Sys":           float64(metrics.Sys),
		"Lookups":       float64(metrics.Lookups),
		"Mallocs":       float64(metrics.Mallocs),
		"Frees":         float64(metrics.Frees),
		"HeapAlloc":     float64(metrics.HeapAlloc),
		"HeapSys":       float64(metrics.HeapSys),
		"HeapIdle":      float64(metrics.HeapIdle),
		"HeapInuse":     float64(metrics.HeapInuse),
		"HeapReleased":  float64(metrics.HeapReleased),
		"HeapObjects":   float64(metrics.HeapObjects),
		"StackInuse":    float64(metrics.StackInuse),
		"StackSys":      float64(metrics.StackSys),
		"MSpanInuse":    float64(metrics.MSpanInuse),
		"MSpanSys":      float64(metrics.MSpanSys),
		"MCacheInuse":   float64(metrics.MCacheInuse),
		"MCacheSys":     float64(metrics.MCacheSys),
		"BuckHashSys":   float64(metrics.BuckHashSys),
		"GCSys":         float64(metrics.GCSys),
		"OtherSys":      float64(metrics.OtherSys),
		"NextGC":        float64(metrics.NextGC),
		"LastGC":        float64(metrics.LastGC),
		"PauseTotalNs":  float64(metrics.PauseTotalNs),
		"NumGC":         float64(metrics.NumGC),
		"NumForcedGC":   float64(metrics.NumForcedGC),
		"GCCPUFraction": float64(metrics.GCCPUFraction),
	}
}

func (c Collector) CollectMetrics() {

	log.Printf("Start collect metrics")
	c.GetMetricsFromMemStats()

	for metricName := range c.memStatsMetric {
		metricValue := c.memStatsMetric[metricName]
		log.Printf("Filed %s, value %v", metricName, metricValue)
		c.service.CreateGaugeMetric(metricName, metricValue)
	}

	c.service.CreateCounterMetric("PollCount", 1)
	c.service.CreateGaugeMetric("RandomValue", rand.Float64())

	log.Printf("Finish collect metrics")
}

func (c Collector) SendMetrics() {
	log.Printf("Send metric")
	counterMetrics, gaugeMetrics := c.service.GetAllMetrics()
	var methodURL = "/update/"

	for metricName := range gaugeMetrics {
		metricValue := gaugeMetrics[metricName]
		log.Printf("Send metric %s", metricName)
		metricData := models.Metrics{
			ID:    metricName,
			MType: "gauge",
			Value: &metricValue,
		}

		var buf bytes.Buffer
		encoder := json.NewEncoder(&buf)
		encoder.Encode(metricData)
		res, err := c.httpClient.SendPostRequest(methodURL, buf)
		if err != nil {
			log.Printf("Cannot send request to server to set metric %s", metricName)
			continue
		}
		defer res.Body.Close()

	}
	for metricName := range counterMetrics {
		metricValue := counterMetrics[metricName]
		log.Printf("Send metric %s", metricName)
		metricData := models.Metrics{
			ID:    metricName,
			MType: "counter",
			Delta: &metricValue,
		}
		var buf bytes.Buffer
		encoder := json.NewEncoder(&buf)
		encoder.Encode(metricData)

		metricsData, errCompress := compress.Compress(buf)
		if errCompress != nil {
			log.Printf("Cannot compress data of metric %s", metricName)
			continue
		}

		res, err := c.httpClient.SendPostRequest(methodURL, *metricsData)
		if err != nil {
			log.Printf("Cannot send request to server to set metric %s", metricName)
			continue
		}
		defer res.Body.Close()

	}
	log.Printf("Finish send metrics")

}
