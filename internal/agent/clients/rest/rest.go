package rest

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/sincin-v/collector/internal/storage"
)

var reportInterval int64 = 10

func SendMetric(channel chan storage.MetricStorage) {
	for {
		time.Sleep(time.Duration(reportInterval) * time.Second)
		metricsStorage := <-channel
		log.Printf("Send metric")
		for metricName := range metricsStorage.Metrics {
			metric := metricsStorage.Metrics[metricName]
			metricValue := metric.GetValueString()
			metricType := metric.GetType()
			url := fmt.Sprintf("http://localhost:8080/update/%s/%s/%s", metricType, metricName, metricValue)
			log.Printf("Send request to url: %s", url)
			_, err := http.Post(url, "text/plain", nil)
			if err != nil {
				log.Fatalf("Error to send request %s", url)
				continue
			}
		}
	}

}
