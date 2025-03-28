package main

import (
	"net/http"

	"github.com/sincin-v/collector/internal/server/handlers"
)

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/update/{metricType}/{metricName}/{metricValue}", handlers.UpdateMetricHandler)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}

}
