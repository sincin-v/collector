package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"

	"github.com/sincin-v/collector/internal/server/handlers"
)

func main() {

	router := chi.NewRouter()

	router.Post("/update/{metricType}/{metricName}/{metricValue}", handlers.UpdateMetricHandler)
	router.Get("/value/{metricType}/{metricName}", handlers.GetMetricHandler)
	router.Get("/", handlers.GetAllMetricsHandler)

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		panic(err)
	}

}
