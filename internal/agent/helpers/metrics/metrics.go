package metrics

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"time"

	"github.com/sincin-v/collector/internal/agent/config"
	"github.com/sincin-v/collector/internal/compress"
	"github.com/sincin-v/collector/internal/models"
)

var PollCountValue int = 0

type MemMetrics struct {
}

type MetricsService interface {
	UpdateCounterMetric(context.Context, string, int64) error
	UpdateGaugeMetric(context.Context, string, float64) error
	GetMetric(context.Context, string, string) (string, error)
	GetAllMetrics(context.Context) (map[string]int64, map[string]float64)
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

func (c Collector) StartCollectMetrics(ctx context.Context, pollInterval time.Duration) {
	for {
		c.CollectMetrics(ctx)
		time.Sleep(pollInterval)
	}
}

func (c Collector) StartSendMetrics(reportInterval time.Duration) {
	ctx := context.Background()
	for {
		time.Sleep(reportInterval)
		// c.SendMetrics()
		c.SendMetricsJSON(ctx)
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

func (c Collector) CollectMetrics(ctx context.Context) {

	log.Printf("Start collect metrics")
	c.GetMetricsFromMemStats()
	var err error
	for metricName := range c.memStatsMetric {
		metricValue := c.memStatsMetric[metricName]
		log.Printf("Filed %s, value %v", metricName, metricValue)
		if errLoopMetric := c.service.UpdateGaugeMetric(ctx, metricName, metricValue); errLoopMetric != nil {
			err = errors.Join(err, fmt.Errorf("update metric %s error: %w", metricName, errLoopMetric))
		}

	}

	if errCounterMetric := c.service.UpdateCounterMetric(ctx, "PollCount", 1); errCounterMetric != nil {
		err = errors.Join(err, fmt.Errorf("update metric PollCount error: %w", errCounterMetric))
	}
	if errGaugeMetric := c.service.UpdateGaugeMetric(ctx, "RandomValue", rand.Float64()); errGaugeMetric != nil {
		err = errors.Join(err, fmt.Errorf("update metric RandomValue error: %w", errGaugeMetric))
	}
	if err != nil {
		log.Printf("Update finish with error:  %s", err)
	}
	log.Printf("Finish collect metrics")
}

func (c Collector) SendMetricsJSON(ctx context.Context) {
	log.Printf("Start send metrics")
	counterMetrics, gaugeMetrics := c.service.GetAllMetrics(ctx)
	var methodURL = "/updates/"
	var metricsArr []models.Metrics

	for metricName := range gaugeMetrics {
		metricValue := gaugeMetrics[metricName]

		metricObj := models.Metrics{
			ID:    metricName,
			MType: config.GaugeMetricType,
			Value: &metricValue,
		}
		metricsArr = append(metricsArr, metricObj)
	}
	for metricName := range counterMetrics {
		metricValue := counterMetrics[metricName]

		metricObj := models.Metrics{
			ID:    metricName,
			MType: config.CounterMetricType,
			Delta: &metricValue,
		}
		metricsArr = append(metricsArr, metricObj)
	}

	sendObj := metricsArr // models.MetricsArray{Metrics: metricsArr}

	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	errEncode := encoder.Encode(sendObj)
	if errEncode != nil {
		log.Printf("Cannot encode data err: %s", errEncode)
		return
	}

	sendData, errCompress := compress.Compress(buf)
	if errCompress != nil {
		log.Printf("Cannot compress data of metrics")
		return
	}
	log.Printf("SEND DATA")
	res, err := c.httpClient.SendPostRequest(methodURL, *sendData)
	if err != nil {
		log.Printf("Cannot send request to server to set metrics Error: %s", err)
		return
	}
	defer func() {
		if errBodyClose := res.Body.Close(); errBodyClose != nil {
			err = errors.Join(err, fmt.Errorf("close body error: %w", errBodyClose))
		}
	}()
}

func (c Collector) SendMetrics(ctx context.Context) {
	log.Printf("Send metric")
	counterMetrics, gaugeMetrics := c.service.GetAllMetrics(ctx)
	var methodURL = "/update/"

	for metricName := range gaugeMetrics {
		metricValue := gaugeMetrics[metricName]
		log.Printf("Send metric %s", metricName)
		metricData := models.Metrics{
			ID:    metricName,
			MType: config.GaugeMetricType,
			Value: &metricValue,
		}

		var buf bytes.Buffer
		encoder := json.NewEncoder(&buf)
		errEncode := encoder.Encode(metricData)
		if errEncode != nil {
			log.Printf("Cannot encode data err: %s", errEncode)
			continue
		}

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
		defer func() {
			if errBodyClose := res.Body.Close(); errBodyClose != nil {
				err = errors.Join(err, fmt.Errorf("close body error: %w", errBodyClose))
			}
		}()

	}
	for metricName := range counterMetrics {
		metricValue := counterMetrics[metricName]
		log.Printf("Send metric %s", metricName)
		metricData := models.Metrics{
			ID:    metricName,
			MType: config.CounterMetricType,
			Delta: &metricValue,
		}
		var buf bytes.Buffer
		encoder := json.NewEncoder(&buf)
		errEncode := encoder.Encode(metricData)
		if errEncode != nil {
			log.Printf("Cannot encode data err: %s", errEncode)
			continue
		}

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
		defer func() {
			if errBodyClose := res.Body.Close(); errBodyClose != nil {
				err = errors.Join(err, fmt.Errorf("close body error: %w", errBodyClose))
			}
		}()
	}
	log.Printf("Finish send metrics")

}
