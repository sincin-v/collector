package rest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sincin-v/collector/internal/models"
)

func TestHttpClient_SendPostRequest(t *testing.T) {
	type fields struct {
		baseURL string
	}
	type args struct {
		metricName  string
		metricType  string
		metricValue int64

		statusCode int
	}
	tests := []struct {
		name    string
		fields  fields
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
				// partsPath := strings.Split(r.URL.Path, "/")
				// if partsPath[1] != "update" {
				// 	t.Errorf("Error: There is no 'update' on url path")
				// }
				// inputMetricType := partsPath[2]
				// inputMetricName := partsPath[3]
				// inputMetricValue := partsPath[4]
				// if inputMetricType == "" || inputMetricName == "" || inputMetricValue == "" {
				// 	t.Errorf("ERROR: The url path invalid %s", r.URL.Path)
				// }
				w.WriteHeader(tt.args.statusCode)
			}))
			defer ts.Close()
			h := HTTPClient{
				baseURL: ts.URL,
			}

			metricData := models.Metrics{
				ID: tt.args.metricName, MType: tt.args.metricType,
			}
			metricData.Delta = &tt.args.metricValue

			var body bytes.Buffer
			encoder := json.NewEncoder(&body)
			encoder.Encode(metricData)
			got, err := h.SendPostRequest("/update/", body)
			if (err != nil) != tt.wantErr {
				t.Errorf("HttpClient.SendPostRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.StatusCode != tt.want {
				t.Errorf("HttpClient.SendPostRequest() = %v, want %v", got, tt.want)
			}
			defer got.Body.Close()
		})
	}
}
