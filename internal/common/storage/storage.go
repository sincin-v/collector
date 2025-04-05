package storage

import (
	"fmt"
	"strconv"
	"sync"
)

type MetricStorage struct {
	mu      sync.RWMutex
	gauge   map[string]float64
	counter map[string]int64
}

func New() MetricStorage {
	return MetricStorage{
		gauge:   map[string]float64{},
		counter: map[string]int64{},
	}
}

func (ms *MetricStorage) CreateGaugeMetric(name string, value float64) {
	ms.mu.Lock()
	ms.gauge[name] = value
	ms.mu.Unlock()
}

func (ms *MetricStorage) CreateCounterMetric(name string, value int64) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	_, ok := ms.counter[name]
	if !ok {
		ms.counter[name] = value
		return
	}
	ms.counter[name] += value

}

func (ms *MetricStorage) GetMetric(metricType string, metricName string) (string, error) {
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

func (ms *MetricStorage) GetAllCountersMetrics() map[string]int64 {
	return ms.counter
}
func (ms *MetricStorage) GetAllGaugeMetrics() map[string]float64 {
	return ms.gauge
}
