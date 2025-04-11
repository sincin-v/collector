package storage

import (
	"reflect"
	"testing"
)

func TestMetricStorage_CreateGaugeMetric(t *testing.T) {
	type fields struct {
		gauge   map[string]float64
		counter map[string]int64
	}
	type args struct {
		n string
		v float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name:   "positive test set gauge metric",
			fields: fields{gauge: make(map[string]float64)},
			args:   args{"testMetric", 0.001},
			want:   "0.001",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := MemStorage{
				gauge:   tt.fields.gauge,
				counter: tt.fields.counter,
			}
			ms.CreateGaugeMetric(tt.args.n, tt.args.v)
			if got, _ := ms.GetMetric("gauge", tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MetricStorage.CreateGaugeMetric() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetricStorage_CreateCounterMetric(t *testing.T) {
	type fields struct {
		gauge   map[string]float64
		counter map[string]int64
	}
	type args struct {
		n string
		v int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name:   "positive test set counter metric",
			fields: fields{counter: make(map[string]int64)},
			args:   args{"testMetric", 1},
			want:   "1",
		},
		{
			name:   "positive test set exist counter metric",
			fields: fields{counter: map[string]int64{"testMetric": 1}},
			args:   args{"testMetric", 1},
			want:   "2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := MemStorage{
				gauge:   tt.fields.gauge,
				counter: tt.fields.counter,
			}
			ms.CreateCounterMetric(tt.args.n, tt.args.v)
			if got, _ := ms.GetMetric("counter", tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MetricStorage.CreateCounterMetric() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetricStorage_GetMetric(t *testing.T) {
	type fields struct {
		gauge   map[string]float64
		counter map[string]int64
	}
	type args struct {
		t string
		n string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name:   "positive test get counter metric",
			fields: fields{counter: map[string]int64{"testMetric": 1}},
			args:   args{"counter", "testMetric"},
			want:   "1",
		},
		{
			name:   "positive test get gauge metric",
			fields: fields{gauge: map[string]float64{"testMetric": 0.1}},
			args:   args{"gauge", "testMetric"},
			want:   "0.1",
		},
		{
			name:    "positive test get not exist gauge metric",
			fields:  fields{gauge: make(map[string]float64)},
			args:    args{"gauge", "testMetric"},
			want:    "",
			wantErr: true,
		},
		{
			name:    "positive test get not exist counter metric",
			fields:  fields{counter: make(map[string]int64)},
			args:    args{"counter", "testMetric"},
			want:    "",
			wantErr: true,
		},
		{
			name:    "positive test get not avalible metric type",
			fields:  fields{counter: make(map[string]int64)},
			args:    args{"histogram", "testMetric"},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := MemStorage{
				gauge:   tt.fields.gauge,
				counter: tt.fields.counter,
			}
			got, err := ms.GetMetric(tt.args.t, tt.args.n)
			if (err != nil) != tt.wantErr {
				t.Errorf("MetricStorage.GetMetric() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MetricStorage.GetMetric() = %v, want %v", got, tt.want)
			}
		})
	}
}
