package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateMetricHandler(t *testing.T) {
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
			url := fmt.Sprintf("/update/%s/%s/%s", tt.args.metricType, tt.args.metricName, tt.args.metricValue)
			request := httptest.NewRequest(tt.args.httpMethod, url, nil)
			request.SetPathValue("metricType", tt.args.metricType)
			request.SetPathValue("metricName", tt.args.metricName)
			request.SetPathValue("metricValue", tt.args.metricValue)
			w := httptest.NewRecorder()
			UpdateMetricHandler(w, request)
			res := w.Result()

			if tt.want.code != res.StatusCode {
				t.Errorf("StatusCode (%d) are not %d", res.StatusCode, tt.want.code)
			}
			defer res.Body.Close()
		})
	}
}
