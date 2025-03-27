package handlers

import (
	"log"
	"net/http"

	"github.com/sincin-v/collector/internal/storage"
)

var metricStorage = storage.MetricStorage{Metrics: make(map[string]storage.Metric)}

func UpdateMetricHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		metricType := req.PathValue("metricType")
		metricName := req.PathValue("metricName")
		metricValue := req.PathValue("metricValue")
		log.Printf("Method: %s Url: %s, metricType: %s, metricName: %s, metricValue: %s", req.Method, req.URL.Path, metricType, metricName, metricValue)

		var metric = metricStorage.GetMetrics(metricName)
		if metric == nil {
			if metricType == `gauge` {
				metric = &storage.GaugeMetric{Name: metricName}
			} else if metricType == `counter` {
				metric = &storage.CounterMetric{Name: metricName}
			} else {
				log.Printf("Invalid type of new metric (%s)", metricType)
				res.WriteHeader(http.StatusBadRequest)
				return
			}
		} else {
			if metric.GetType() != metricType {
				log.Printf("Invalid type of exist metric (%s)", metricType)
				res.WriteHeader(http.StatusBadRequest)
				return
			}
		}
		err := storage.SetMetricValue(metric, metricValue, metricType)
		if err != nil {
			log.Printf("Invalid value (%s) for type (%s)", metricValue, metricType)
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		metricStorage.CreateMetric(metricName, metric)
		log.Printf("New value of metric %s (type: %s) = %s", metricName, metricType, metric.GetValueString())
		res.WriteHeader(http.StatusOK)
	} else {
		log.Printf("Error: %d", http.StatusMethodNotAllowed)
	}
}
