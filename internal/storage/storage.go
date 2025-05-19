package storage

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"

	"github.com/sincin-v/collector/internal/models"
	"github.com/sincin-v/collector/internal/server/config"
)

type MemStorage struct {
	mu      sync.RWMutex
	gauge   map[string]float64
	counter map[string]int64
}

func NewMemStorage() MemStorage {
	return MemStorage{
		gauge:   map[string]float64{},
		counter: map[string]int64{},
	}
}
func (ms *MemStorage) UpdateMetricsByBatch(_ context.Context, metrics []models.Metrics) error {

	ms.mu.Lock()
	defer ms.mu.Unlock()
	for _, metricObj := range metrics {
		switch metricObj.MType {
		case config.GaugeMetricType:
			ms.gauge[metricObj.ID] = *metricObj.Value
		case config.CounterMetricType:
			_, ok := ms.counter[metricObj.ID]
			if !ok {
				ms.counter[metricObj.ID] = *metricObj.Delta
				continue
			}
			ms.counter[metricObj.ID] += *metricObj.Delta
		default:
			continue
		}
	}
	return nil
}

func (ms *MemStorage) UpdateGaugeMetric(_ context.Context, name string, value float64) error {
	ms.mu.Lock()
	ms.gauge[name] = value
	ms.mu.Unlock()
	return nil
}

func (ms *MemStorage) UpdateCounterMetric(_ context.Context, name string, value int64) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	_, ok := ms.counter[name]
	if !ok {
		ms.counter[name] = value
		return nil
	}
	ms.counter[name] += value
	return nil
}

func (ms *MemStorage) GetMetric(_ context.Context, metricType string, metricName string) (string, error) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	switch metricType {
	case config.GaugeMetricType:
		value, ok := ms.gauge[metricName]
		if !ok {
			return "", errors.New("metric does not exists")
		}
		return strconv.FormatFloat(value, 'f', -1, 64), nil
	case config.CounterMetricType:
		value, ok := ms.counter[metricName]
		if !ok {
			return "", errors.New("metric does not exists")
		}
		return fmt.Sprintf("%d", value), nil
	default:
		return "", fmt.Errorf("there is no metric type %s", metricType)
	}
}

func (ms *MemStorage) GetAllCountersMetrics(_ context.Context) map[string]int64 {
	return ms.counter
}

func (ms *MemStorage) GetAllGaugeMetrics(_ context.Context) map[string]float64 {
	return ms.gauge
}

func (ms *MemStorage) HealthCheck(ctx context.Context) error {
	return nil
}

