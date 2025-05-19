package collector

import (
	"encoding/json"
	"errors"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/sincin-v/collector/internal/models"
	"github.com/sincin-v/collector/internal/server/config"
)

type MetricsService interface {
	GetAllCountersMetrics(context.Context) map[string]int64
	GetAllGaugeMetrics(context.Context) map[string]float64
	UpdateCounterMetric(context.Context, string, int64) error
	UpdateGaugeMetric(context.Context, string, float64) error
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

func (mc MetricCollector) SaveMetrics(ctx context.Context, interval int) error {
	for {
		time.Sleep(time.Duration(interval) * time.Second)

		metricsMap := map[string]interface{}{
			config.CounterMetricType: mc.service.GetAllCountersMetrics(ctx),
			config.GaugeMetricType:   mc.service.GetAllGaugeMetrics(ctx),
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

func (mc MetricCollector) RestoreMetrics(ctx context.Context) error {

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
		if errCounterMetric := mc.service.UpdateCounterMetric(ctx, counterMetric, metricsData.Counter[counterMetric]); errCounterMetric != nil {
			err = errors.Join(err, fmt.Errorf("update metric %s error: %w", counterMetric, errCounterMetric))
		}
	}

	for gaugeMetric := range metricsData.Gauge {
		if errGaugeMetric := mc.service.UpdateGaugeMetric(ctx, gaugeMetric, metricsData.Gauge[gaugeMetric]); errGaugeMetric != nil {
			err = errors.Join(err, fmt.Errorf("update metric %s error: %w", gaugeMetric, errGaugeMetric))
		}
	}

	if err != nil {
		return err
	}

	return nil

}
