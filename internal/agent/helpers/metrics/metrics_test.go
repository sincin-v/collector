package metrics

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"context"
	"testing"
	"time"

	"github.com/sincin-v/collector/internal/agent/clients/rest"
	"github.com/sincin-v/collector/internal/agent/config"
	"github.com/sincin-v/collector/internal/service"
	"github.com/sincin-v/collector/internal/storage"
)

func TestCollector_CollectMetrics(t *testing.T) {
	tests := []struct {
		name      string
		want      string
		wantValue string
	}{
		{
			name:      "positive test collect metrics",
			want:      "PollCount",
			wantValue: "1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			st := storage.NewMemStorage()
			s := service.New(&st)

			hs := rest.New("localhost:8888", []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second})

			c := Collector{
				service:    s,
				httpClient: hs,
			}
			c.CollectMetrics(ctx)

			if gotValue, err := c.service.GetMetric(ctx, config.CounterMetricType, "PollCount"); gotValue != tt.wantValue || err != nil {
				t.Errorf("PollCount (%s) are not  eq %s", gotValue, tt.wantValue)

			}
		})
	}
}

func TestCollector_SendMetrics(t *testing.T) {
	type fields struct {
		countMetricName  string
		countMetricValue int64
		gaugeMetricName  string
		gaugeMetricValue float64
	}
	type args struct {
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name:   "positive test send metrics",
			fields: fields{"testCountMetric", 1, "testGaugeMetrics", 0.1},
		},
		{
			name:   "negative test send metrics",
			fields: fields{"testCountMetric", 1, "testGaugeMetrics", 0.1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				partsPath := strings.Split(r.URL.Path, "/")
				if partsPath[1] != "update" {
					t.Errorf("Error: There is no 'update' on url path")
				}
				w.WriteHeader(http.StatusOK)
			}))
			defer ts.Close()
			st := storage.NewMemStorage()
			ctx := context.Background()
			_ = st.UpdateCounterMetric(ctx, tt.fields.countMetricName, tt.fields.countMetricValue)
			_ = st.UpdateGaugeMetric(ctx, tt.fields.gaugeMetricName, tt.fields.gaugeMetricValue)
			s := service.New(&st)

			hs := rest.New(ts.URL, []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second})

			c := Collector{
				service:    s,
				httpClient: hs,
			}

			c.SendMetrics(ctx)
		})
	}
}
