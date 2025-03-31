package rest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sincin-v/collector/internal/storage"
)

func TestSendMetricHelper(t *testing.T) {
	type args struct {
		channel chan storage.MetricStorage
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "positive test",
			args: args{make(chan storage.MetricStorage)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			}))

			go SendMetricHelper(tt.args.channel, ts.URL)

			tt.args.channel <- storage.MetricStorage{Metrics: map[string]storage.Metric{"TestMetric": &storage.CounterMetric{Name: "TestMetric", Type: "counter", Value: 1}}}
		})
	}
}
