package rest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/sincin-v/collector/internal/models"
)

func TestHttpClient_SendPostRequest(t *testing.T) {
	type args struct {
		metricName  string
		metricType  string
		metricValue int64

		statusCode int
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "positive test send request",
			args: args{"TestMetric", "counter", 1, http.StatusOK},
			want: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.args.statusCode)
			}))
			defer ts.Close()
			h := HTTPClient{
				baseURL:        ts.URL,
				retryIntervals: []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second},
			}

			metricData := models.Metrics{
				ID: tt.args.metricName, MType: tt.args.metricType,
			}
			metricData.Delta = &tt.args.metricValue

			var body bytes.Buffer
			encoder := json.NewEncoder(&body)
			errEncode := encoder.Encode(metricData)
			if errEncode != nil {
				t.Errorf("Cannot encode test data")
			}
			got, err := h.SendPostRequest("/update/", body)
			if (err != nil) != tt.wantErr {
				t.Errorf("HttpClient.SendPostRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.StatusCode != tt.want {
				t.Errorf("HttpClient.SendPostRequest() = %v, want %v", got, tt.want)
			}
			defer func() {
				if errBodyClose := got.Body.Close(); errBodyClose != nil {
					t.Errorf("Close body error %s", errBodyClose)
				}
			}()
		})
	}
}
