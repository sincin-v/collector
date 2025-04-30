package collector

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/sincin-v/collector/internal/models"
)

type MetricsService interface {
	GetAllCountersMetrics() map[string]int64
	GetAllGaugeMetrics() map[string]float64
	UpdateCounterMetric(string, int64) error
	UpdateGaugeMetric(string, float64) error
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
	var err error
	metricsData := models.RestoredDataModel{}
	if err = json.Unmarshal(savedData, &metricsData); err != nil {
		return err
	}

	for counterMetric := range metricsData.Counter {
		if errCounterMetric := mc.service.UpdateCounterMetric(counterMetric, metricsData.Counter[counterMetric]); errCounterMetric != nil {
			err = errors.Join(err, fmt.Errorf("update metric %s error: %w", counterMetric, errCounterMetric))
		}
	}

	for gaugeMetric := range metricsData.Gauge {
		if errGaugeMetric := mc.service.UpdateGaugeMetric(gaugeMetric, metricsData.Gauge[gaugeMetric]); errGaugeMetric != nil {
			err = errors.Join(err, fmt.Errorf("update metric %s error: %w", gaugeMetric, errGaugeMetric))
		}
	}

	if err != nil {
		return err
	}

	return nil

}
