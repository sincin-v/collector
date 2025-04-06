package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

type MetricsService interface {
	CreateCounterMetric(string, int64)
	CreateGaugeMetric(string, float64)
	GetMetric(string, string) (string, error)
	GetAllMetrics() (map[string]int64, map[string]float64)
}

type Handler struct {
	service MetricsService
}

func New(s MetricsService) Handler {
	return Handler{
		service: s,
	}
}

func (h Handler) UpdateMetricHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		log.Printf("Error: %d", http.StatusMethodNotAllowed)
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	metricType := req.PathValue("metricType")
	metricName := req.PathValue("metricName")
	metricValue := req.PathValue("metricValue")
	log.Printf("Method: %s Url: %s, metricType: %s, metricName: %s, metricValue: %s", req.Method, req.URL.Path, metricType, metricName, metricValue)

	switch metricType {
	case "gauge":
		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			log.Printf("Invalid value (%s) for type (%s)", metricValue, metricType)
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		h.service.CreateGaugeMetric(metricName, value)
	case "counter":
		value, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			log.Printf("Invalid value (%s) for type (%s)", metricValue, metricType)
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		h.service.CreateCounterMetric(metricName, value)
	default:
		log.Printf("Invalid type of new metric (%s)", metricType)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	newMetricvalue, err := h.service.GetMetric(metricType, metricName)
	if err != nil {
		log.Printf("Cannot set new value (%s) for metric '%s' Error: %s", metricValue, metricName, err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Printf("New value of metric %s (type: %s) = %s", metricName, metricType, newMetricvalue)
	res.WriteHeader(http.StatusOK)
}

func (h Handler) GetMetricHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		log.Printf("Error: %d", http.StatusMethodNotAllowed)
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	metricType := req.PathValue("metricType")
	metricName := req.PathValue("metricName")
	metric, err := h.service.GetMetric(metricType, metricName)
	if err != nil {
		log.Printf("Metric %s not found. Error: %s", metricName, err)
		res.WriteHeader(http.StatusNotFound)
		return
	}
	res.WriteHeader(http.StatusOK)
	io.WriteString(res, metric)
}

func (h Handler) GetAllMetricsHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		log.Printf("Error: %d", http.StatusMethodNotAllowed)
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	counterMetric, gaugeMetrics := h.service.GetAllMetrics()
	for metricName := range counterMetric {
		metricValue, err := h.service.GetMetric("counter", metricName)
		if err != nil {
			log.Printf("Cannot get value of metric '%s' . Error: %s", metricName, err)
			continue
		}
		io.WriteString(res, fmt.Sprintf("%s = %s\n", metricName, metricValue))
	}
	for metricName := range gaugeMetrics {
		metricValue, err := h.service.GetMetric("gauge", metricName)
		if err != nil {
			log.Printf("Cannot get value of metric '%s' . Error: %s", metricName, err)
			continue
		}
		io.WriteString(res, fmt.Sprintf("%s = %s\n", metricName, metricValue))
	}
	res.WriteHeader(http.StatusOK)
}
