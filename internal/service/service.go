package service

import (
	"context"
	"github.com/sincin-v/collector/internal/models"
)

type metricStorage interface {
	UpdateCounterMetric(context.Context, string, int64) error
	UpdateGaugeMetric(context.Context, string, float64) error
	GetMetric(context.Context, string, string) (string, error)
	GetAllCountersMetrics(context.Context) map[string]int64
	GetAllGaugeMetrics(context.Context) map[string]float64
	UpdateMetricsByBatch(context.Context, []models.Metrics) error
}

type MetricsService struct {
	metricStorage metricStorage
}

func New(s metricStorage) MetricsService {
	return MetricsService{
		metricStorage: s,
	}
}

func (s MetricsService) UpdateMetricsByBatch(ctx context.Context, metrics []models.Metrics) error {
	return s.metricStorage.UpdateMetricsByBatch(ctx, metrics)
}

func (s MetricsService) UpdateGaugeMetric(ctx context.Context, metricName string, value float64) error {
	return s.metricStorage.UpdateGaugeMetric(ctx, metricName, value)
}

func (s MetricsService) UpdateCounterMetric(ctx context.Context, metricName string, value int64) error {
	return s.metricStorage.UpdateCounterMetric(ctx, metricName, value)
}

func (s MetricsService) GetMetric(ctx context.Context, metricType string, metricName string) (string, error) {
	metricValue, err := s.metricStorage.GetMetric(ctx, metricType, metricName)
	return metricValue, err
}

func (s MetricsService) GetAllMetrics(ctx context.Context, ) (map[string]int64, map[string]float64) {
	counterMetric := s.metricStorage.GetAllCountersMetrics(ctx)
	gaugeMetrics := s.metricStorage.GetAllGaugeMetrics(ctx)
	return counterMetric, gaugeMetrics
}

func (s MetricsService) GetAllCountersMetrics(ctx context.Context, ) map[string]int64 {
	return s.metricStorage.GetAllCountersMetrics(ctx)
}

func (s MetricsService) GetAllGaugeMetrics(ctx context.Context, ) map[string]float64 {
	return s.metricStorage.GetAllGaugeMetrics(ctx)
}
