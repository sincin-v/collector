package rest

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"strings"

	"github.com/sincin-v/collector/internal/storage"
)

func SendMetricHelper(channel chan storage.MetricStorage, baseURL string) {
	metricsStorage := <-channel
	log.Printf("Send metric")
	for metricName := range metricsStorage.Metrics {
		metric := metricsStorage.Metrics[metricName]
		metricValue := metric.GetValueString()
		metricType := metric.GetType()
		var url = fmt.Sprintf("%s/update/%s/%s/%s", baseURL, metricType, metricName, metricValue)
		if !strings.HasPrefix(url, "http") {
			url = fmt.Sprintf("http://%s", url)
		}
		log.Printf("Send request to url: %s", url)
		resp, err := http.Post(url, "text/plain", nil)
		if err != nil {
			log.Fatalf("Error to send request %s", url)
			continue
		}
		defer resp.Body.Close()
		log.Printf("Finish send request")
	}
}

func SendMetric(channel chan storage.MetricStorage, baseURL string, reportInterval int) {
	for {
		time.Sleep(time.Duration(reportInterval) * time.Second)
		SendMetricHelper(channel, baseURL)
	}
}
