package collector

import "time"

type MetricsService interface {
	FlushAllMetrics(string) error
	DumpAllMetrics(string) error
}

type MetricCollector struct {
	service MetricsService
	path    string
}

func New(s MetricsService, path string) MetricCollector {
	return MetricCollector{
		service: s,
		path:    path,
	}
}

func (mc MetricCollector) SaveMetrics(interval int) error {
	for {
		time.Sleep(time.Duration(interval) * time.Second)
		err := mc.service.FlushAllMetrics(mc.path)
		if err != nil {
			return err
		}

	}
}

func (mc MetricCollector) RestoreMetrics() error {
	return mc.service.DumpAllMetrics(mc.path)
}
