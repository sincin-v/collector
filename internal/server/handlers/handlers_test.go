package handlers

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
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
		code        int
		metricValue string
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
				code:        http.StatusOK,
				metricValue: "1",
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
				code:        http.StatusOK,
				metricValue: "1.0",
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
				code:        http.StatusBadRequest,
				metricValue: "1.0",
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
				code:        http.StatusBadRequest,
				metricValue: "1.0",
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
				code:        http.StatusBadRequest,
				metricValue: "1.0",
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
				code:        http.StatusMethodNotAllowed,
				metricValue: "1.0",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := logger.Initialize("INFO"); err != nil {
				t.Errorf("Could not initialize logger Error: %s", err)
			}
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			s := storage.NewMockMetricStorage(ctrl)
			ctxt := context.Background()

			switch tt.args.metricType {
			case config.CounterMetricType:
				metricValue, errParseInt := strconv.ParseInt(tt.args.metricValue, 10, 64)
				if errParseInt == nil {
					s.EXPECT().UpdateCounterMetric(ctxt, tt.args.metricName, metricValue).Return(nil)
					s.EXPECT().GetMetric(ctxt, tt.args.metricType, tt.args.metricName).Return(tt.want.metricValue, nil)
				}
			case config.GaugeMetricType:
				metricValue, errParseFloat := strconv.ParseFloat(tt.args.metricValue, 64)
				if errParseFloat == nil {
					s.EXPECT().UpdateGaugeMetric(ctxt, tt.args.metricName, metricValue).Return(nil)
					s.EXPECT().GetMetric(ctxt, tt.args.metricType, tt.args.metricName).Return(tt.want.metricValue, nil)
				}
			}

			h := &Handler{
				service: service.New(s),
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
		code        int
		metricValue string
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
				code:        http.StatusOK,
				metricValue: "1",
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
				code:        http.StatusMethodNotAllowed,
				metricValue: "1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := logger.Initialize("INFO"); err != nil {
				t.Errorf("Could not initialize logger Error: %s", err)
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			storage := storage.NewMockMetricStorage(ctrl)
			ctxt := context.Background()
			if tt.want.code == http.StatusOK {
				storage.EXPECT().GetMetric(ctxt, tt.args.metricType, tt.args.metricName).Return(tt.want.metricValue, nil)
			}
			service := service.New(storage)
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
			fields: fields{"testCounterMetric", 1, "testGaugeMetric", 1.1},
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
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			storage := storage.NewMockMetricStorage(ctrl)
			if tt.want.code == http.StatusOK {
				storage.EXPECT().GetAllCountersMetrics(gomock.Any()).Return(map[string]int64{tt.fields.counterMetricName: tt.fields.counterMetricValue})
				storage.EXPECT().GetAllGaugeMetrics(gomock.Any()).Return(map[string]float64{tt.fields.gaugeMetricName: tt.fields.gaugeMetricValue})
			}

			service := service.New(storage)
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
