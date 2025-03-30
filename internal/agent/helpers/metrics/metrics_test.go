package metrics

import (
	"fmt"
	"testing"
	"time"

	"github.com/sincin-v/collector/internal/storage"
)

func TestGetMetrics(t *testing.T) {
	type args struct {
		ch chan storage.MetricStorage
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "positive test get metrics",
			args: args{make(chan storage.MetricStorage)},
			want: "PollCount",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go GetMetrics(tt.args.ch)
			select {
			case result := <-tt.args.ch:
				_, ok := result.Metrics[tt.want]
				if !ok {
					t.Fatal(fmt.Printf("Test failed. Metric '%s' does not exist", tt.want))
				}
			case <-time.After(3 * time.Second):
				t.Fatal("Test timed out")
			}

		})
	}
}
