package main

import (
	"flag"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"

	"github.com/sincin-v/collector/internal/server/handlers"
)

func main() {

	hostStr := flag.String("a", "localhost:8080", "Listen host and port")

	flag.Parse()
	log.Printf("Start server work")

	router := chi.NewRouter()

	router.Post("/update/{metricType}/{metricName}/{metricValue}", handlers.UpdateMetricHandler)
	router.Get("/value/{metricType}/{metricName}", handlers.GetMetricHandler)
	router.Get("/", handlers.GetAllMetricsHandler)

	err := http.ListenAndServe(*hostStr, router)
	if err != nil {
		panic(err)
	}

}
