package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

type MetricsService interface {
	CreateMetric(string, string, string) (string, error)
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
	}

	metricType := req.PathValue("metricType")
	metricName := req.PathValue("metricName")
	metricValue := req.PathValue("metricValue")
	log.Printf("Method: %s Url: %s, metricType: %s, metricName: %s, metricValue: %s", req.Method, req.URL.Path, metricType, metricName, metricValue)
	newMetricvalue, err := h.service.CreateMetric(metricType, metricName, metricValue)
	if err != nil {
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
	}
	metricType := req.PathValue("metricType")
	metricName := req.PathValue("metricName")
	var metric, err = h.service.GetMetric(metricType, metricName)
	if err != nil {
		log.Printf("Error. Metric %s not found", metricName)
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
	}
	counterMetric, gaugeMetrics := h.service.GetAllMetrics()
	for metricName := range counterMetric {
		metricValue, _ := h.service.GetMetric("counter", metricName)
		io.WriteString(res, fmt.Sprintf("%s = %s\n", metricName, metricValue))
	}
	for metricName := range gaugeMetrics {
		metricValue, _ := h.service.GetMetric("gauge", metricName)
		io.WriteString(res, fmt.Sprintf("%s = %s\n", metricName, metricValue))
	}
	res.WriteHeader(http.StatusOK)
}
