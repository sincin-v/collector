package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/sincin-v/collector/internal/logger"
	"github.com/sincin-v/collector/internal/models"
	"github.com/sincin-v/collector/internal/server/config"
)

type MetricsService interface {
	UpdateCounterMetric(context.Context, string, int64) error
	UpdateGaugeMetric(context.Context, string, float64) error
	GetMetric(context.Context, string, string) (string, error)
	GetAllMetrics(context.Context) (map[string]int64, map[string]float64)
	UpdateMetricsByBatch(context.Context, []models.Metrics) error
	HealthCheck(context.Context) error
}

type Handler struct {
	service  MetricsService
}

func New(s MetricsService) (*Handler, error) {
	return &Handler{
		service:  s,
	}, nil
}

func (h Handler) UpdateMetricHandler(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	if req.Method != http.MethodPost {
		logger.Log.Errorf("[Handler] Error: %d", http.StatusMethodNotAllowed)
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	metricType := req.PathValue("metricType")
	metricName := req.PathValue("metricName")
	metricValue := req.PathValue("metricValue")
	logger.Log.Infof("[Handler] Method: %s Url: %s, metricType: %s, metricName: %s, metricValue: %s", req.Method, req.URL.Path, metricType, metricName, metricValue)

	var err error

	switch metricType {
	case config.GaugeMetricType:
		value, errParseFloat := strconv.ParseFloat(metricValue, 64)
		if errParseFloat != nil {
			logger.Log.Errorf("[Handler] Invalid value (%s) for type (%s)", metricValue, metricType)
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		if errUpGaugeMetric := h.service.UpdateGaugeMetric(ctx, metricName, value); errUpGaugeMetric != nil {
			err = errors.Join(err, fmt.Errorf("cannot update metric %s, error: %s", metricName, errUpGaugeMetric))
		}
	case config.CounterMetricType:
		value, errParseInt := strconv.ParseInt(metricValue, 10, 64)
		if errParseInt != nil {
			logger.Log.Errorf("[Handler] Invalid value (%s) for type (%s)", metricValue, metricType)
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		if errUpCouterMetric := h.service.UpdateCounterMetric(ctx, metricName, value); errUpCouterMetric != nil {
			err = errors.Join(err, fmt.Errorf("cannot update metric %s, error: %s", metricName, errUpCouterMetric))
		}
	default:
		logger.Log.Errorf("[Handler] Invalid type of new metric (%s)", metricType)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	newMetricvalue, errGetMetric := h.service.GetMetric(ctx, metricType, metricName)
	if errGetMetric != nil {
		err = errors.Join(err, fmt.Errorf("cannot get new value for %s", metricName))
	}
	if err != nil {
		logger.Log.Errorf("[Handler] Cannot set new value (%s) for metric '%s' Error: %s", metricValue, metricName, err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	logger.Log.Infof("[Handler] New value of metric %s (type: %s) = %s", metricName, metricType, newMetricvalue)
	res.WriteHeader(http.StatusOK)
}

func (h Handler) GetMetricHandler(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	if req.Method != http.MethodGet {
		logger.Log.Errorf("[Handler] Error: %d", http.StatusMethodNotAllowed)
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	metricType := req.PathValue("metricType")
	metricName := req.PathValue("metricName")
	metric, err := h.service.GetMetric(ctx, metricType, metricName)
	if err != nil {
		logger.Log.Errorf("[Handler] Metric %s not found. Error: %s", metricName, err)
		res.WriteHeader(http.StatusNotFound)
		return
	}
	res.WriteHeader(http.StatusOK)
	if _, err := io.WriteString(res, metric); err != nil {
		logger.Log.Errorf("[Handler] Could not return metric %s data Error: %s", metricName, err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h Handler) GetAllMetricsHandler(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	ctxTimeout, cancel := context.WithTimeout(ctx, config.OperationTimeout)
	defer cancel()
	res.Header().Set("Content-Type", "text/html")
	if req.Method != http.MethodGet {
		logger.Log.Errorf("[Handler] Error: %d", http.StatusMethodNotAllowed)
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	counterMetric, gaugeMetrics := h.service.GetAllMetrics(ctxTimeout)
	for metricName := range counterMetric {
		metricValue, err := h.service.GetMetric(ctxTimeout, config.CounterMetricType, metricName)
		if err != nil {
			logger.Log.Errorf("[Handler] Cannot get value of metric '%s' . Error: %s", metricName, err)
			continue
		}
		if _, err := io.WriteString(res, fmt.Sprintf("%s = %s\n", metricName, metricValue)); err != nil {
			logger.Log.Errorf("[Handler] Could not return metrics data Error: %s", err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	for metricName := range gaugeMetrics {
		metricValue, err := h.service.GetMetric(ctxTimeout, config.GaugeMetricType, metricName)
		if err != nil {
			logger.Log.Errorf("[Handler] Cannot get value of metric '%s' . Error: %s", metricName, err)
			continue
		}
		if _, err := io.WriteString(res, fmt.Sprintf("%s = %s\n", metricName, metricValue)); err != nil {
			logger.Log.Errorf("Could not return metrics data Error: %s", err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		res.WriteHeader(http.StatusOK)
	}
}

func (h Handler) UpdateMetricJSONHandler(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	ctxTimeout, cancel := context.WithTimeout(ctx, config.OperationTimeout)
	defer cancel()
	if req.Method != http.MethodPost {
		logger.Log.Errorf("Error: %d", http.StatusMethodNotAllowed)
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var inputData models.Metrics
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&inputData); err != nil {
		logger.Log.Errorf("[Handler] Cannot decode input body Error: %s", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	resp := models.Metrics{
		ID:    inputData.ID,
		MType: inputData.MType,
	}
	var err error

	switch inputData.MType {
	case config.GaugeMetricType:
		if errUpGaugeMetric := h.service.UpdateGaugeMetric(ctxTimeout, inputData.ID, *inputData.Value); errUpGaugeMetric != nil {
			err = errors.Join(err, fmt.Errorf("cannot update metric %s, error: %S", inputData.ID, errUpGaugeMetric))
		}
	case config.CounterMetricType:

		if errUpCounterMetric := h.service.UpdateCounterMetric(ctxTimeout, inputData.ID, *inputData.Delta); errUpCounterMetric != nil {
			err = errors.Join(err, fmt.Errorf("cannot update metric %s, error: %S", inputData.ID, errUpCounterMetric))
		}
	default:
		logger.Log.Infof("[Handler] Invalid type of new metric (%s)", inputData.MType)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	newMetricvalue, errGetMetric := h.service.GetMetric(ctxTimeout, inputData.MType, inputData.ID)
	if errGetMetric != nil {
		err = errors.Join(err, fmt.Errorf("cannot get new value for %s", inputData.ID))
	}
	if err != nil {
		logger.Log.Errorf("[Handler] Cannot set new value for metric '%s' Error: %s", inputData.ID, err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	switch inputData.MType {
	case config.GaugeMetricType:
		value, err := strconv.ParseFloat(newMetricvalue, 64)
		if err != nil {
			logger.Log.Debugf("[Handler] Invalid value (%s) for type (gauge)", value)
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		resp.Value = &value
	case config.CounterMetricType:
		value, err := strconv.ParseInt(newMetricvalue, 10, 64)
		if err != nil {
			logger.Log.Debugf("[Handler] Invalid value (%s) for type (counter)", value)
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		resp.Delta = &value
	}

	logger.Log.Infof("[Handler] New value of metric %s (type: %s) = %s", inputData.ID, inputData.MType, newMetricvalue)
	res.WriteHeader(http.StatusOK)
	res.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(res)
	if err := encoder.Encode(resp); err != nil {
		res.WriteHeader(http.StatusBadRequest)
	}
}

func (h Handler) UpdateManyMetricsJSONHandler(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	ctxTimeout, cancel := context.WithTimeout(ctx, config.OperationTimeout)
	defer cancel()
	if req.Method != http.MethodPost {
		logger.Log.Errorf("[Handler] Error: %d", http.StatusMethodNotAllowed)
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	res.Header().Set("Content-Type", "application/json")

	var inputData []models.Metrics
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&inputData); err != nil {
		logger.Log.Errorf("[Handler] Cannot decode input body Error: %s", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	err := h.service.UpdateMetricsByBatch(ctxTimeout, inputData)
	if err != nil {
		logger.Log.Error("[Handler] Cannot update metrics in db. Error: %s", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	logger.Log.Debug("[Handler] Metrics has been updated")
	res.WriteHeader(http.StatusOK)
	res.Header().Set("Content-Type", "application/json")
}

func (h Handler) GetMetricJSONHandler(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	ctxTimeout, cancel := context.WithTimeout(ctx, config.OperationTimeout)
	defer cancel()
	if req.Method != http.MethodPost {
		logger.Log.Errorf("[Handler] Error: %d", http.StatusMethodNotAllowed)
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	res.Header().Set("Content-Type", "application/json")

	var inputData models.Metrics
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&inputData); err != nil {
		logger.Log.Errorf("[Handler] Cannot decode input body Error: %s", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	metric, err := h.service.GetMetric(ctxTimeout, inputData.MType, inputData.ID)
	if err != nil {
		logger.Log.Errorf("[Handler] Metric %s not found. Error: %s", inputData.ID, err)
		res.WriteHeader(http.StatusNotFound)
		return
	}

	resp := models.Metrics{
		ID:    inputData.ID,
		MType: inputData.MType,
	}
	switch inputData.MType {
	case config.GaugeMetricType:
		value, err := strconv.ParseFloat(metric, 64)
		if err != nil {
			logger.Log.Debugf("[Handler] Invalid value (%s) for type (gauge)", value)
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		resp.Value = &value
	case config.CounterMetricType:
		value, err := strconv.ParseInt(metric, 10, 64)
		if err != nil {
			logger.Log.Debugf("[Handler] Invalid value (%s) for type (counter)", value)
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

func (h Handler) Ping(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	ctxTimeout, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	err := h.service.HealthCheck(ctxTimeout)
	if err != nil {
		logger.Log.Error("[Handler] Ping Error: %s", err)
	}
}
