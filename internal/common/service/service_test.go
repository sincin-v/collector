package service

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/sincin-v/collector/internal/common/storage"
)

func TestMetricsService_CreateMetric(t *testing.T) {

	type args struct {
		metricType  string
		metricName  string
		metricValue string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "positive create gauge metric",
			args: args{"gauge", "testGaugeMetric", "0.000001"},
			want: "0.000001",
		},
		{
			name:    "negative create gauge metric with invalid value",
			args:    args{"gauge", "testGaugeMetric", "invalidValue"},
			want:    "",
			wantErr: true,
		},
		{
			name: "positive create counter metric",
			args: args{"counter", "testCounterMetric", "11"},
			want: "11",
		},
		{
			name:    "negative create counter metric with invalid value",
			args:    args{"counter", "testCounterMetric", "invalidValue"},
			want:    "",
			wantErr: true,
		},
		{
			name:    "negative create metric with invalid type",
			args:    args{"histogram", "testhistogramMetric", "invalidValue"},
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
			got, err := s.CreateMetric(tt.args.metricType, tt.args.metricName, tt.args.metricValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("MetricsService.CreateMetric() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MetricsService.CreateMetric() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetricsService_GetMetric(t *testing.T) {
	type fields struct {
		metricType  string
		metricName  string
		metricValue string
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
			fields: fields{"counter", "testCounterMetric", "1"},
			args:   args{"counter", "testCounterMetric"},
			want:   "1",
		},
		{
			name:    "negative get counter metric",
			fields:  fields{"counter", "testCounterMetric", "1"},
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
			s.CreateMetric(tt.fields.metricType, tt.fields.metricName, tt.fields.metricValue)
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
				metricValue := fmt.Sprintf("%f", tt.fields.gaugeMetrics[gaugeMetricName])
				s.CreateMetric("gauge", gaugeMetricName, metricValue)
			}
			for counterMetricName := range tt.fields.counterMetrics {
				metricValue := fmt.Sprintf("%d", tt.fields.counterMetrics[counterMetricName])
				s.CreateMetric("counter", counterMetricName, metricValue)
			}

			got, got1 := s.GetAllMetrics()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MetricsService.GetAllMetrics() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("MetricsService.GetAllMetrics() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
