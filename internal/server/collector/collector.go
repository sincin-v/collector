package collector

import (
	"encoding/json"
	"os"
	"time"

	"github.com/sincin-v/collector/internal/models"
)

type MetricsService interface {
	GetAllCountersMetrics() map[string]int64
	GetAllGaugeMetrics() map[string]float64
	UpdateCounterMetric(string, int64)
	UpdateGaugeMetric(string, float64)
}

type MetricCollector struct {
	service MetricsService
	path    string
}

func New(s MetricsService, path string) MetricCollector {
	return MetricCollector{
		service: s,
		path:    path,
	}
}

func (mc MetricCollector) SaveMetrics(interval int) error {
	for {
		time.Sleep(time.Duration(interval) * time.Second)

		metricsMap := map[string]interface{}{
			"counter": mc.service.GetAllCountersMetrics(),
			"gauge":   mc.service.GetAllGaugeMetrics(),
		}
		resultData, errJSON := json.MarshalIndent(metricsMap, "", "   ")
		if errJSON != nil {
			return errJSON
		}

		err := os.WriteFile(mc.path, resultData, 0666)
		if err != nil {
			return err
		}

	}
}

func (mc MetricCollector) RestoreMetrics() error {

	savedData, errOpenFile := os.ReadFile(mc.path)
	if errOpenFile != nil {
		return errOpenFile
	}
	metricsData := models.RestoredDataModel{}
	if err := json.Unmarshal(savedData, &metricsData); err != nil {
		return err
	}

	for counterMetric := range metricsData.Counter {
		mc.service.UpdateCounterMetric(counterMetric, metricsData.Counter[counterMetric])
	}

	for gaugeMetric := range metricsData.Gauge {
		mc.service.UpdateGaugeMetric(gaugeMetric, metricsData.Gauge[gaugeMetric])
	}

	return nil

}
