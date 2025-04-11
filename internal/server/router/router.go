package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/sincin-v/collector/internal/server/handlers"
	mw "github.com/sincin-v/collector/internal/server/middlewares"
	"github.com/sincin-v/collector/internal/service"
	"github.com/sincin-v/collector/internal/storage"
)

func CreateRouter(storage *storage.MemStorage) *chi.Mux {
	service := service.New(storage)
	h := handlers.New(service)
	router := chi.NewRouter()

	router.Post("/update/{metricType}/{metricName}/{metricValue}", mw.WithLogger(h.UpdateMetricHandler))
	router.Get("/value/{metricType}/{metricName}", mw.WithLogger(h.GetMetricHandler))
	router.Get("/", mw.WithLogger(h.GetAllMetricsHandler))

	return router
}
