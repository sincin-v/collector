package service

type metricStorage interface {
	CreateCounterMetric(string, int64)
	CreateGaugeMetric(string, float64)
	GetMetric(string, string) (string, error)
	GetAllCountersMetrics() map[string]int64
	GetAllGaugeMetrics() map[string]float64
}

type MetricsService struct {
	metricStorage metricStorage
}

func New(s metricStorage) MetricsService {
	return MetricsService{
		metricStorage: s,
	}
}

func (s MetricsService) CreateGaugeMetric(metricName string, value float64) {
	s.metricStorage.CreateGaugeMetric(metricName, value)
}

func (s MetricsService) CreateCounterMetric(metricName string, value int64) {
	s.metricStorage.CreateCounterMetric(metricName, value)
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
