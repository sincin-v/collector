package service

import (
	"errors"
	"fmt"
	"log"
	"strconv"
)

type MetricStorage interface {
	CreateCounterMetric(string, int64)
	CreateGaugeMetric(string, float64)
	GetMetric(string, string) (string, error)
	GetAllCountersMetrics() map[string]int64
	GetAllGaugeMetrics() map[string]float64
}

type MetricsService struct {
	metricStorage MetricStorage
}

func New(s MetricStorage) MetricsService {
	return MetricsService{
		metricStorage: s,
	}
}

func (s MetricsService) CreateMetric(metricType string, metricName string, metricValue string) (string, error) {
	switch metricType {
	case "gauge":
		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			log.Printf("Invalid value (%s) for type (%s)", metricValue, metricType)
			return "", err
		}
		s.metricStorage.CreateGaugeMetric(metricName, value)
	case "counter":
		value, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			log.Printf("Invalid value (%s) for type (%s)", metricValue, metricType)
			return "", err
		}
		s.metricStorage.CreateCounterMetric(metricName, value)
	default:
		log.Printf("Invalid type of new metric (%s)", metricType)
		return "", errors.New(fmt.Sprintf("Invalid type of new metric (%s)", metricType))
	}
	newSetValue, _ := s.metricStorage.GetMetric(metricType, metricName)
	return newSetValue, nil
}

func (s MetricsService) GetMetric(metricType string, metricName string) (string, error) {
	metricValue, err := s.metricStorage.GetMetric(metricType, metricName)
	return metricValue, err
}

func (s MetricsService) GetAllMetrics() (map[string]int64, map[string]float64) {
	counterMetric := s.metricStorage.GetAllCountersMetrics()
	gaugeMetrics := s.metricStorage.GetAllGaugeMetrics()
	return counterMetric, gaugeMetrics
}

func (s MetricsService) GetStorage() MetricStorage {
	return s.metricStorage
}
