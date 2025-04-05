package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sincin-v/collector/internal/common/service"
	"github.com/sincin-v/collector/internal/common/storage"
)

func TestHandler_UpdateMetricHandler(t *testing.T) {
	type args struct {
		metricType  string
		metricName  string
		metricValue string
		httpMethod  string
	}
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "positive test update counter metric handler",
			args: args{
				metricType:  "counter",
				metricName:  "testCounterMetric",
				metricValue: "1",
				httpMethod:  http.MethodPost,
			},
			want: want{
				code: 200,
			},
		},
		{
			name: "positive test update gauge metric handler",
			args: args{
				metricType:  "gauge",
				metricName:  "testGaugeMetric",
				metricValue: "1.0",
				httpMethod:  http.MethodPost,
			},
			want: want{
				code: 200,
			},
		},
		{
			name: "negative test update gauge metric handler with invalid type",
			args: args{
				metricType:  "counter",
				metricName:  "testGaugeMetric",
				metricValue: "1.0",
				httpMethod:  http.MethodPost,
			},
			want: want{
				code: 400,
			},
		},
		{
			name: "negative test update metric handler with invalid value",
			args: args{
				metricType:  "gauge",
				metricName:  "testGaugeMetric",
				metricValue: "invalidValue",
				httpMethod:  http.MethodPost,
			},
			want: want{
				code: 400,
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
				code: 400,
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
				code: 405,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := storage.New()
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
			defer res.Body.Close()

		})
	}
}

func TestHandler_GetMetricHandler(t *testing.T) {
	type fields struct {
		metricType  string
		metricName  string
		metricValue string
	}
	type args struct {
		metricType string
		metricName string
		httpMethod string
	}
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name:   "positive test get counter metric handler",
			fields: fields{"counter", "testCounterMetric", "1"},
			args: args{
				metricType: "counter",
				metricName: "testCounterMetric",
				httpMethod: http.MethodGet,
			},
			want: want{
				code: 200,
			},
		},
		{
			name:   "negative test invalid method",
			fields: fields{"counter", "testCounterMetric", "1"},
			args: args{
				metricType: "histogram",
				metricName: "testHistogramMetric",
				httpMethod: http.MethodPost,
			},
			want: want{
				code: 405,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := storage.New()
			service := service.New(&storage)
			service.CreateMetric(tt.fields.metricType, tt.fields.metricName, tt.fields.metricValue)
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
			defer res.Body.Close()

		})
	}
}

func TestHandler_GetAllMetricsHandler(t *testing.T) {
	type fields struct {
		metricType  string
		metricName  string
		metricValue string
	}
	type args struct {
		httpMethod string
	}
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name:   "positive test get all counter metric handler",
			fields: fields{"counter", "testCounterMetric", "1"},
			args: args{
				httpMethod: http.MethodGet,
			},
			want: want{
				code: 200,
			},
		},
		{
			name:   "positive test get all gauge metric handler",
			fields: fields{"gauge", "testCounterMetric", "1.0"},
			args: args{
				httpMethod: http.MethodGet,
			},
			want: want{
				code: 200,
			},
		},
		{
			name:   "negative test invalid method",
			fields: fields{"counter", "testCounterMetric", "1"},
			args: args{
				httpMethod: http.MethodPost,
			},
			want: want{
				code: 405,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := storage.New()
			service := service.New(&storage)
			service.CreateMetric(tt.fields.metricType, tt.fields.metricName, tt.fields.metricValue)
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
			defer res.Body.Close()
		})
	}
}
