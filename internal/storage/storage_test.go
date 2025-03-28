package storage

import (
	"reflect"
	"testing"
)

func TestMetricStorage_GetMetrics(t *testing.T) {
	type fields struct {
		Metrics map[string]Metric
	}

	type args struct {
		n string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Metric
	}{
		{
			name:   "positive test get",
			fields: fields{Metrics: map[string]Metric{"testMetric": &GaugeMetric{Name: "testMetric", Type: "Gauge", Value: 0.001}}},
			args:   args{"testMetric"},
			want:   &GaugeMetric{Name: "testMetric", Type: "Gauge", Value: 0.001},
		},
		{
			name:   "positive test get nil field",
			fields: fields{Metrics: map[string]Metric{"testMetric": &GaugeMetric{Name: "testMetric", Type: "Gauge", Value: 0.001}}},
			args:   args{"testNotExistMetric"},
			want:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := MetricStorage{
				Metrics: tt.fields.Metrics,
			}
			if got := ms.GetMetrics(tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MetricStorage.GetMetrics() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetricStorage_CreateMetric(t *testing.T) {
	type fields struct {
		Metrics map[string]Metric
	}
	type args struct {
		n string
		m Metric
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Metric
	}{
		{
			name:   "positive test create",
			fields: fields{make(map[string]Metric)},
			args:   args{"testMetric", &GaugeMetric{Name: "testMetric", Type: "Gauge", Value: 0.001}},
			want:   &GaugeMetric{Name: "testMetric", Type: "Gauge", Value: 0.001},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := MetricStorage{
				Metrics: tt.fields.Metrics,
			}
			if got := ms.CreateMetric(tt.args.n, tt.args.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MetricStorage.CreateMetric() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetMetricValue(t *testing.T) {
	type args struct {
		m Metric
		v string
		t string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "positive test set metric value",
			args:    args{m: &GaugeMetric{Name: "testMetric", Type: "Gauge", Value: 0.001}, v: "1.1", t: "gauge"},
			wantErr: false,
		},
		{
			name:    "positive test set metric invalid value",
			args:    args{m: &GaugeMetric{Name: "testMetric", Type: "Gauge", Value: 0.001}, v: "invalidValue", t: "gauge"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SetMetricValue(tt.args.m, tt.args.v, tt.args.t); (err != nil) != tt.wantErr {
				t.Errorf("SetMetricValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCounterMetric_GetType(t *testing.T) {
	type fields struct {
		Name  string
		Type  string
		Value int64
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{name: "positive test get type of counter metric",
			fields: fields{Name: "testMetric", Type: "counter", Value: 1},
			want:   "counter"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := CounterMetric{
				Name:  tt.fields.Name,
				Type:  tt.fields.Type,
				Value: tt.fields.Value,
			}
			if got := m.GetType(); got != tt.want {
				t.Errorf("CounterMetric.GetType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCounterMetric_GetValueString(t *testing.T) {
	type fields struct {
		Name  string
		Type  string
		Value int64
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{name: "positive test get value of counter metric",
			fields: fields{Name: "testMetric", Type: "counter", Value: 1},
			want:   "1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &CounterMetric{
				Name:  tt.fields.Name,
				Type:  tt.fields.Type,
				Value: tt.fields.Value,
			}
			if got := m.GetValueString(); got != tt.want {
				t.Errorf("CounterMetric.GetValueString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCounterMetric_Set(t *testing.T) {
	type fields struct {
		Name  string
		Type  string
		Value int64
	}
	type args struct {
		v string
		t string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "positive test set counter value",
			fields:  fields{Name: "testMetric", Type: "counter", Value: 1},
			args:    args{v: "2", t: "counter"},
			wantErr: false,
		},
		{
			name:    "positive test set counter invalid value",
			fields:  fields{Name: "testMetric", Type: "counter", Value: 1},
			args:    args{v: "invalidValue", t: "counter"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &CounterMetric{
				Name:  tt.fields.Name,
				Type:  tt.fields.Type,
				Value: tt.fields.Value,
			}
			if err := m.Set(tt.args.v, tt.args.t); (err != nil) != tt.wantErr {
				t.Errorf("CounterMetric.Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGaugeMetric_GetType(t *testing.T) {
	type fields struct {
		Name  string
		Type  string
		Value float64
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{name: "positive test get type of gauge metric",
			fields: fields{Name: "testMetric", Type: "gauge", Value: 0.01},
			want:   "gauge"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := GaugeMetric{
				Name:  tt.fields.Name,
				Type:  tt.fields.Type,
				Value: tt.fields.Value,
			}
			if got := m.GetType(); got != tt.want {
				t.Errorf("GaugeMetric.GetType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGaugeMetric_GetValueString(t *testing.T) {
	type fields struct {
		Name  string
		Type  string
		Value float64
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{name: "positive test get value of gauge metric",
			fields: fields{Name: "testMetric", Type: "gauge", Value: 0.01},
			want:   "0.01"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &GaugeMetric{
				Name:  tt.fields.Name,
				Type:  tt.fields.Type,
				Value: tt.fields.Value,
			}
			if got := m.GetValueString(); got != tt.want {
				t.Errorf("GaugeMetric.GetValueString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGaugeMetric_Set(t *testing.T) {
	type fields struct {
		Name  string
		Type  string
		Value float64
	}
	type args struct {
		v string
		t string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "positive test set gauge value",
			fields:  fields{Name: "testMetric", Type: "gauge", Value: 0.01},
			args:    args{v: "2", t: "gauge"},
			wantErr: false,
		},
		{
			name:    "positive test set gauge invalid value",
			fields:  fields{Name: "testMetric", Type: "gauge", Value: 0.01},
			args:    args{v: "invalidValue", t: "gauge"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &GaugeMetric{
				Name:  tt.fields.Name,
				Type:  tt.fields.Type,
				Value: tt.fields.Value,
			}
			if err := m.Set(tt.args.v, tt.args.t); (err != nil) != tt.wantErr {
				t.Errorf("GaugeMetric.Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
