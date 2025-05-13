package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/sincin-v/collector/internal/server/handlers"
	zipMw "github.com/sincin-v/collector/internal/server/middlewares/compressing"
	logMw "github.com/sincin-v/collector/internal/server/middlewares/logging"
	"github.com/sincin-v/collector/internal/service"
)

func CreateRouter(service *service.MetricsService) (*chi.Mux, error) {

	h, err := handlers.New(service)
	if err != nil {
		return nil, err
	}
	router := chi.NewRouter()

	router.Use(logMw.LoggerMiddleware)
	router.Use(zipMw.CompressMiddleware)
	router.Post("/updates/", h.UpdateManyMetricsJSONHandler)
	router.Post("/update/", h.UpdateMetricJSONHandler)
	router.Post("/update/{metricType}/{metricName}/{metricValue}", h.UpdateMetricHandler)
	router.Post("/value/", h.GetMetricJSONHandler)
	router.Get("/value/{metricType}/{metricName}", h.GetMetricHandler)
	router.Get("/ping", h.Ping)
	router.Get("/", h.GetAllMetricsHandler)

	return router, nil
}
