package rest

import (
	"net/http"
	"net/http/httptest"
	"strings"

	"testing"
)

func TestHttpClient_SendPostRequest(t *testing.T) {
	type fields struct {
		baseURL string
	}
	type args struct {
		url        string
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
			args: args{"/update/counter/TestMetric/1", http.StatusOK},
			want: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				partsPath := strings.Split(r.URL.Path, "/")
				if partsPath[1] != "update" {
					t.Errorf("Error: There is no 'update' on url path")
				}
				inputMetricType := partsPath[2]
				inputMetricName := partsPath[3]
				inputMetricValue := partsPath[4]
				if inputMetricType == "" || inputMetricName == "" || inputMetricValue == "" {
					t.Errorf("ERROR: The url path invalid %s", r.URL.Path)
				}
				w.WriteHeader(tt.args.statusCode)
			}))

			h := HttpClient{
				baseURL: ts.URL,
			}
			got, err := h.SendPostRequest(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("HttpClient.SendPostRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.StatusCode != tt.want {
				t.Errorf("HttpClient.SendPostRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}
