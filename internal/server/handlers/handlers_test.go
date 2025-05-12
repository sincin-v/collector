package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"context"

	"github.com/sincin-v/collector/internal/logger"
	"github.com/sincin-v/collector/internal/server/config"
	"github.com/sincin-v/collector/internal/service"
	"github.com/sincin-v/collector/internal/storage"
)

func TestHandler_UpdateMetricHandler(t *testing.T) {
	type args struct {
		metricType  string
		metricName  string
		metricValue string
		httpMethod  string
	}
	type want struct {
		code int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "positive test update counter metric handler",
			args: args{
				metricType:  config.CounterMetricType,
				metricName:  "testCounterMetric",
				metricValue: "1",
				httpMethod:  http.MethodPost,
			},
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name: "positive test update gauge metric handler",
			args: args{
				metricType:  config.GaugeMetricType,
				metricName:  "testGaugeMetric",
				metricValue: "1.0",
				httpMethod:  http.MethodPost,
			},
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name: "negative test update gauge metric handler with invalid type",
			args: args{
				metricType:  config.CounterMetricType,
				metricName:  "testGaugeMetric",
				metricValue: "1.0",
				httpMethod:  http.MethodPost,
			},
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "negative test update metric handler with invalid value",
			args: args{
				metricType:  config.GaugeMetricType,
				metricName:  "testGaugeMetric",
				metricValue: "invalidValue",
				httpMethod:  http.MethodPost,
			},
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "negative test update new metric handler with invalid type",
			args: args{
				metricType:  "histogram",
				metricName:  "testHistogramMetric",
				metricValue: "1.0",
				httpMethod:  http.MethodPost,
			},
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "negative test invalid method type",
			args: args{
				metricType:  "histogram",
				metricName:  "testHistogramMetric",
				metricValue: "1.0",
				httpMethod:  http.MethodGet,
			},
			want: want{
				code: http.StatusMethodNotAllowed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := logger.Initialize("INFO"); err != nil {
				t.Errorf("Could not initialize logger Error: %s", err)
			}
			s := storage.NewMemStorage()
			h := &Handler{
				service: service.New(&s),
			}

			url := fmt.Sprintf("/update/%s/%s/%s", tt.args.metricType, tt.args.metricName, tt.args.metricValue)
			request := httptest.NewRequest(tt.args.httpMethod, url, nil)
			request.SetPathValue("metricType", tt.args.metricType)
			request.SetPathValue("metricName", tt.args.metricName)
			request.SetPathValue("metricValue", tt.args.metricValue)
			w := httptest.NewRecorder()
			h.UpdateMetricHandler(w, request)
			res := w.Result()

			if tt.want.code != res.StatusCode {
				t.Errorf("StatusCode (%d) are not %d", res.StatusCode, tt.want.code)
			}
			defer func() {
				if errBodyClose := res.Body.Close(); errBodyClose != nil {
					t.Errorf("Close body error %s", errBodyClose)
				}
			}()

		})
	}
}

func TestHandler_GetMetricHandler(t *testing.T) {
	type fields struct {
		metricType  string
		metricName  string
		metricValue int64
	}
	type args struct {
		metricType string
		metricName string
		httpMethod string
	}
	type want struct {
		code int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name:   "positive test get counter metric handler",
			fields: fields{config.CounterMetricType, "testCounterMetric", 1},
			args: args{
				metricType: config.CounterMetricType,
				metricName: "testCounterMetric",
				httpMethod: http.MethodGet,
			},
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name:   "negative test invalid method",
			fields: fields{config.CounterMetricType, "testCounterMetric", 1},
			args: args{
				metricType: "histogram",
				metricName: "testHistogramMetric",
				httpMethod: http.MethodPost,
			},
			want: want{
				code: http.StatusMethodNotAllowed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := logger.Initialize("INFO"); err != nil {
				t.Errorf("Could not initialize logger Error: %s", err)
			}
			storage := storage.NewMemStorage()
			service := service.New(&storage)
			ctx := context.Background()
			_ = service.UpdateCounterMetric(ctx, tt.fields.metricName, tt.fields.metricValue)
			h := &Handler{
				service: service,
			}

			url := fmt.Sprintf("/update/%s/%s", tt.args.metricType, tt.args.metricName)
			request := httptest.NewRequest(tt.args.httpMethod, url, nil)
			request.SetPathValue("metricType", tt.args.metricType)
			request.SetPathValue("metricName", tt.args.metricName)
			w := httptest.NewRecorder()
			h.GetMetricHandler(w, request)
			res := w.Result()

			if tt.want.code != res.StatusCode {
				t.Errorf("StatusCode (%d) are not %d", res.StatusCode, tt.want.code)
			}
			defer func() {
				if errBodyClose := res.Body.Close(); errBodyClose != nil {
					t.Errorf("Close body error %s", errBodyClose)
				}
			}()

		})
	}
}

func TestHandler_GetAllMetricsHandler(t *testing.T) {
	type fields struct {
		counterMetricName  string
		counterMetricValue int64
		gaugeMetricName    string
		gaugeMetricValue   float64
	}
	type args struct {
		httpMethod string
	}
	type want struct {
		code int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name:   "positive test get all counter metric handler",
			fields: fields{"testCounterMetric", 1, "testGaugeMetric", 1.0},
			args: args{
				httpMethod: http.MethodGet,
			},
			want: want{
				code: http.StatusOK,
			},
		},

		{
			name:   "negative test invalid method",
			fields: fields{"testCounterMetric", 1, "testGaugeMetric", 1.0},
			args: args{
				httpMethod: http.MethodPost,
			},
			want: want{
				code: http.StatusMethodNotAllowed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := logger.Initialize("INFO"); err != nil {
				t.Errorf("Could not initialize logger Error: %s", err)
			}
			storage := storage.NewMemStorage()
			service := service.New(&storage)
			ctx := context.Background()
			_ = service.UpdateCounterMetric(ctx, tt.fields.counterMetricName, tt.fields.counterMetricValue)
			_ = service.UpdateGaugeMetric(ctx, tt.fields.gaugeMetricName, tt.fields.gaugeMetricValue)
			h := &Handler{
				service: service,
			}

			url := "/"
			request := httptest.NewRequest(tt.args.httpMethod, url, nil)
			w := httptest.NewRecorder()
			h.GetAllMetricsHandler(w, request)
			res := w.Result()

			if tt.want.code != res.StatusCode {
				t.Errorf("StatusCode (%d) are not %d", res.StatusCode, tt.want.code)
			}
			defer func() {
				if errBodyClose := res.Body.Close(); errBodyClose != nil {
					t.Errorf("Close body error %s", errBodyClose)
				}
			}()
		})
	}
}
