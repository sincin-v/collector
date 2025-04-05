package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/sincin-v/collector/internal/common/service"
	"github.com/sincin-v/collector/internal/common/storage"
	"github.com/sincin-v/collector/internal/server/handlers"
)

func CreateRouter(storage *storage.MetricStorage) *chi.Mux {
	service := service.New(storage)
	h := handlers.New(service)
	router := chi.NewRouter()

	router.Post("/update/{metricType}/{metricName}/{metricValue}", h.UpdateMetricHandler)
	router.Get("/value/{metricType}/{metricName}", h.GetMetricHandler)
	router.Get("/", h.GetAllMetricsHandler)

	return router
}
