package service

import (
	"testing"

	"github.com/sincin-v/collector/internal/storage"
)

func TestMetricsService_GetMetric(t *testing.T) {
	type fields struct {
		metricType  string
		metricName  string
		metricValue int64
	}
	type args struct {
		metricType string
		metricName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name:   "positive get counter metric",
			fields: fields{"counter", "testCounterMetric", 1},
			args:   args{"counter", "testCounterMetric"},
			want:   "1",
		},
		{
			name:    "negative get counter metric",
			fields:  fields{"counter", "testCounterMetric", 1},
			args:    args{"counter", "testInvalidCounterMetric"},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := storage.New()
			s := MetricsService{
				metricStorage: &st,
			}
			s.CreateCounterMetric(tt.fields.metricName, tt.fields.metricValue)
			got, err := s.GetMetric(tt.args.metricType, tt.args.metricName)
			if (err != nil) != tt.wantErr {
				t.Errorf("MetricsService.GetMetric() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MetricsService.GetMetric() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetricsService_GetAllMetrics(t *testing.T) {
	type fields struct {
		gaugeMetrics   map[string]float64
		counterMetrics map[string]int64
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]int64
		want1  map[string]float64
	}{
		{
			name: "positive test get all metrics",
			fields: fields{
				gaugeMetrics:   map[string]float64{"testGaugeMetric": 0.000001},
				counterMetrics: map[string]int64{"testCounterMetric": 1},
			},
			want1: map[string]float64{"testGaugeMetric": 0.000001},
			want:  map[string]int64{"testCounterMetric": 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := storage.New()
			s := MetricsService{
				metricStorage: &st,
			}
			for gaugeMetricName := range tt.fields.gaugeMetrics {
				s.CreateGaugeMetric(gaugeMetricName, tt.fields.gaugeMetrics[gaugeMetricName])
			}
			for counterMetricName := range tt.fields.counterMetrics {
				s.CreateCounterMetric(counterMetricName, tt.fields.counterMetrics[counterMetricName])
			}

			got, got1 := s.GetAllMetrics()

			for counterMetricName := range tt.want {
				if tt.want[counterMetricName] != got[counterMetricName] {
					t.Errorf("MetricsService.GetAllMetrics() Metric %s got = %v, want %v", counterMetricName, got[counterMetricName], tt.want[counterMetricName])
				}
			}
			for gaugeMetricName := range tt.want1 {
				if tt.want1[gaugeMetricName] != got1[gaugeMetricName] {
					t.Errorf("MetricsService.GetAllMetrics() Metric %s got = %v, want %v", gaugeMetricName, got1[gaugeMetricName], tt.want[gaugeMetricName])
				}
			}
		})
	}
}

func TestMetricsService_CreateGaugeMetric(t *testing.T) {
	type fields struct {
		metricStorage metricStorage
	}
	type args struct {
		metricName string
		value      float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "positive create gauge metric",
			args: args{"testGaugeMetric", 0.000001},
			want: "0.000001",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			st := storage.New()
			s := MetricsService{
				metricStorage: &st,
			}
			s.CreateGaugeMetric(tt.args.metricName, tt.args.value)
			got, _ := s.GetMetric("gauge", tt.args.metricName)
			if got != tt.want {
				t.Errorf("MetricsService.GetMetric() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetricsService_CreateCounterMetric(t *testing.T) {
	type fields struct {
		metricStorage metricStorage
	}
	type args struct {
		metricName string
		value      int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{

		{
			name: "positive create counter metric",
			args: args{"testCounterMetric", 11},
			want: "11",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := storage.New()
			s := MetricsService{
				metricStorage: &st,
			}
			s.CreateCounterMetric(tt.args.metricName, tt.args.value)
			got, _ := s.GetMetric("counter", tt.args.metricName)
			if got != tt.want {
				t.Errorf("MetricsService.GetMetric() = %v, want %v", got, tt.want)
			}

		})
	}
}
