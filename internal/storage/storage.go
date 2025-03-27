package storage

import (
	"fmt"
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

func (m *GaugeMetric) Set(v string, t string) error {
	nv, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return err
	}
	m.Value = nv
	m.Type = t
	return nil
}

func (m *GaugeMetric) GetValueString() string {
	return fmt.Sprintf("%f", m.Value)
}

func (m GaugeMetric) GetType() string {
	return m.Type
}

type CounterMetric struct {
	Name  string
	Type  string
	Value int64
}

func (m *CounterMetric) Set(v string, t string) error {
	nv, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return err
	}
	m.Value += nv
	m.Type = t
	return nil
}

func (m *CounterMetric) GetValueString() string {
	return fmt.Sprintf("%d", m.Value)
}

func (m CounterMetric) GetType() string {
	return m.Type
}

func SetMetricValue(m Metric, v string, t string) error {
	err := m.Set(v, t)
	return err
}

type MetricStorage struct {
	Metrics map[string]Metric
}

func (ms MetricStorage) CreateMetric(n string, m Metric) Metric {
	ms.Metrics[n] = m
	return m
}

func (ms MetricStorage) GetMetrics(n string) Metric {
	m, ok := ms.Metrics[n]
	if !ok {
		return nil
	}
	return m
}
