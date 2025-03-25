package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type Metric interface {
	Set(v string, t string) error
	GetValueString() string
	GetType() string
}

type GaugeMetric struct {
	Name  string
	Type  string
	Value float64
}

func (g *GaugeMetric) Set(v string, t string) error {
	nv, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return err
	}
	g.Value = nv
	g.Type = t
	return nil
}

func (g *GaugeMetric) GetValueString() string {
	return fmt.Sprintf("%f", g.Value)
}

func (m GaugeMetric) GetType() string {
	return m.Type
}

type CounterMetric struct {
	Name  string
	Type  string
	Value int64
}

func (c *CounterMetric) Set(v string, t string) error {
	nv, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return err
	}
	c.Value += nv
	c.Type = t
	return nil
}

func (g *CounterMetric) GetValueString() string {
	return fmt.Sprintf("%d", g.Value)
}

func (m CounterMetric) GetType() string {
	return m.Type
}

func SetMetricValue(m Metric, v string, t string) error {
	err := m.Set(v, t)
	return err
}

type MetricStorage struct {
	metrics map[string]Metric
}

func (ms MetricStorage) CreateMetric(n string, m Metric) Metric {
	ms.metrics[n] = m
	return m
}

func (ms MetricStorage) GetMetrics(n string) Metric {
	m, ok := ms.metrics[n]
	if !ok {
		return nil
	}
	return m
}

var metricStorage = MetricStorage{metrics: make(map[string]Metric)}

func updateMetricHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		metricType := req.PathValue("metricType")
		metricName := req.PathValue("metricName")
		metricValue := req.PathValue("metricValue")
		log.Printf("Method: %s Url: %s, metricType: %s, metricName: %s, metricValue: %s", req.Method, req.URL.Path, metricType, metricName, metricValue)

		var metric = metricStorage.GetMetrics(metricName)
		if metric == nil {
			if metricType == `gauge` {
				metric = &GaugeMetric{Name: metricName}
			} else if metricType == `counter` {
				metric = &CounterMetric{Name: metricName}
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
		err := SetMetricValue(metric, metricValue, metricType)
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

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/update/{metricType}/{metricName}/{metricValue}", updateMetricHandler)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}

}
