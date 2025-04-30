package storage

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/sincin-v/collector/internal/models"
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
func (ms *MemStorage) UpdateMetricsByBatch(metrics []models.Metrics) error {

	ms.mu.Lock()
	defer ms.mu.Unlock()
	for _, metricObj := range metrics {
		switch metricObj.MType {
		case "gauge":
			ms.gauge[metricObj.ID] = *metricObj.Value
		case "counter":
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

func (ms *MemStorage) UpdateGaugeMetric(name string, value float64) error {
	ms.mu.Lock()
	ms.gauge[name] = value
	ms.mu.Unlock()
	return nil
}

func (ms *MemStorage) UpdateCounterMetric(name string, value int64) error {
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

func (ms *MemStorage) GetMetric(metricType string, metricName string) (string, error) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	switch metricType {
	case "gauge":
		value, ok := ms.gauge[metricName]
		if !ok {
			return "", fmt.Errorf("there is no gauge metric %s", metricName)
		}
		return strconv.FormatFloat(value, 'f', -1, 64), nil
	case "counter":
		value, ok := ms.counter[metricName]
		if !ok {
			return "", fmt.Errorf("there is no counter metric %s", metricName)
		}
		return fmt.Sprintf("%d", value), nil
	default:
		return "", fmt.Errorf("there is no metric type %s", metricType)
	}
}

func (ms *MemStorage) GetAllCountersMetrics() map[string]int64 {
	return ms.counter
}

func (ms *MemStorage) GetAllGaugeMetrics() map[string]float64 {
	return ms.gauge
}
