package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/sincin-v/collector/internal/server/handlers"
	zipMw "github.com/sincin-v/collector/internal/server/middlewares/compressing"
	logMw "github.com/sincin-v/collector/internal/server/middlewares/logging"
	"github.com/sincin-v/collector/internal/service"
	"github.com/sincin-v/collector/internal/storage"
)

func CreateRouter(storage *storage.MemStorage) *chi.Mux {
	service := service.New(storage)
	h := handlers.New(service)
	router := chi.NewRouter()

	router.Use(logMw.LoggerMiddleware)
	router.Use(zipMw.CompressMiddleware)
	// router.Post("/update/", zipMw.CompressMiddleware(h.UpdateMetricJSONHandler))
	router.Post("/update/", h.UpdateMetricJSONHandler)
	router.Post("/update/{metricType}/{metricName}/{metricValue}", h.UpdateMetricHandler)
	// router.Post("/value/", zipMw.CompressMiddleware(h.GetMetricJSONHandler))
	router.Post("/value/", h.GetMetricJSONHandler)
	router.Get("/value/{metricType}/{metricName}", h.GetMetricHandler)
	router.Get("/", h.GetAllMetricsHandler)

	return router
}
