package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/sincin-v/collector/internal/logger"
	"github.com/sincin-v/collector/internal/models"
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
		logger.Log.Errorf("Error: %d", http.StatusMethodNotAllowed)
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	metricType := req.PathValue("metricType")
	metricName := req.PathValue("metricName")
	metricValue := req.PathValue("metricValue")
	logger.Log.Infof("Method: %s Url: %s, metricType: %s, metricName: %s, metricValue: %s", req.Method, req.URL.Path, metricType, metricName, metricValue)

	switch metricType {
	case "gauge":
		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			logger.Log.Errorf("Invalid value (%s) for type (%s)", metricValue, metricType)
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		h.service.CreateGaugeMetric(metricName, value)
	case "counter":
		value, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			logger.Log.Errorf("Invalid value (%s) for type (%s)", metricValue, metricType)
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		h.service.CreateCounterMetric(metricName, value)
	default:
		logger.Log.Errorf("Invalid type of new metric (%s)", metricType)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	newMetricvalue, err := h.service.GetMetric(metricType, metricName)
	if err != nil {
		logger.Log.Errorf("Cannot set new value (%s) for metric '%s' Error: %s", metricValue, metricName, err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	logger.Log.Infof("New value of metric %s (type: %s) = %s", metricName, metricType, newMetricvalue)
	res.WriteHeader(http.StatusOK)
}

func (h Handler) GetMetricHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		logger.Log.Errorf("Error: %d", http.StatusMethodNotAllowed)
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	metricType := req.PathValue("metricType")
	metricName := req.PathValue("metricName")
	metric, err := h.service.GetMetric(metricType, metricName)
	if err != nil {
		logger.Log.Errorf("Metric %s not found. Error: %s", metricName, err)
		res.WriteHeader(http.StatusNotFound)
		return
	}
	res.WriteHeader(http.StatusOK)
	io.WriteString(res, metric)
}

func (h Handler) GetAllMetricsHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/html")
	if req.Method != http.MethodGet {
		logger.Log.Errorf("Error: %d", http.StatusMethodNotAllowed)
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	counterMetric, gaugeMetrics := h.service.GetAllMetrics()
	for metricName := range counterMetric {
		metricValue, err := h.service.GetMetric("counter", metricName)
		if err != nil {
			logger.Log.Errorf("Cannot get value of metric '%s' . Error: %s", metricName, err)
			continue
		}
		io.WriteString(res, fmt.Sprintf("%s = %s\n", metricName, metricValue))
	}
	for metricName := range gaugeMetrics {
		metricValue, err := h.service.GetMetric("gauge", metricName)
		if err != nil {
			logger.Log.Errorf("Cannot get value of metric '%s' . Error: %s", metricName, err)
			continue
		}
		io.WriteString(res, fmt.Sprintf("%s = %s\n", metricName, metricValue))
	}
	res.WriteHeader(http.StatusOK)
}

func (h Handler) UpdateMetricJSONHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		logger.Log.Errorf("Error: %d", http.StatusMethodNotAllowed)
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var inputData models.Metrics
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&inputData); err != nil {
		logger.Log.Errorf("Cannot decode input body Error: %s", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	resp := models.Metrics{
		ID:    inputData.ID,
		MType: inputData.MType,
	}

	switch inputData.MType {
	case "gauge":
		h.service.CreateGaugeMetric(inputData.ID, *inputData.Value)
	case "counter":
		h.service.CreateCounterMetric(inputData.ID, *inputData.Delta)
	default:
		logger.Log.Infof("Invalid type of new metric (%s)", inputData.MType)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	newMetricvalue, err := h.service.GetMetric(inputData.MType, inputData.ID)
	if err != nil {
		logger.Log.Errorf("Cannot set new value for metric '%s' Error: %s", inputData.ID, err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	switch inputData.MType {
	case "gauge":
		value, err := strconv.ParseFloat(newMetricvalue, 64)
		if err != nil {
			logger.Log.Debugf("Invalid value (%s) for type (gauge)", value)
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		resp.Value = &value
	case "counter":
		value, err := strconv.ParseInt(newMetricvalue, 10, 64)
		if err != nil {
			logger.Log.Debugf("Invalid value (%s) for type (counter)", value)
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		resp.Delta = &value
	}

	logger.Log.Infof("New value of metric %s (type: %s) = %s", inputData.ID, inputData.MType, newMetricvalue)
	res.WriteHeader(http.StatusOK)
	res.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(res)
	if err := encoder.Encode(resp); err != nil {
		res.WriteHeader(http.StatusBadRequest)
	}

}

func (h Handler) GetMetricJSONHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		logger.Log.Errorf("Error: %d", http.StatusMethodNotAllowed)
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	res.Header().Set("Content-Type", "application/json")

	var inputData models.Metrics
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&inputData); err != nil {
		logger.Log.Errorf("Cannot decode input body Error: %s", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	metric, err := h.service.GetMetric(inputData.MType, inputData.ID)
	if err != nil {
		logger.Log.Errorf("Metric %s not found. Error: %s", inputData.ID, err)
		res.WriteHeader(http.StatusNotFound)
		return
	}

	resp := models.Metrics{
		ID:    inputData.ID,
		MType: inputData.MType,
	}
	switch inputData.MType {
	case "gauge":
		value, err := strconv.ParseFloat(metric, 64)
		if err != nil {
			logger.Log.Debugf("Invalid value (%s) for type (gauge)", value)
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		resp.Value = &value
	case "counter":
		value, err := strconv.ParseInt(metric, 10, 64)
		if err != nil {
			logger.Log.Debugf("Invalid value (%s) for type (counter)", value)
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		resp.Delta = &value
	}

	res.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(res)
	if err := encoder.Encode(resp); err != nil {
		res.WriteHeader(http.StatusBadRequest)
	}

}
