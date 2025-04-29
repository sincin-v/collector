package router

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/sincin-v/collector/internal/server/clients/db"
	"github.com/sincin-v/collector/internal/server/handlers"
	zipMw "github.com/sincin-v/collector/internal/server/middlewares/compressing"
	logMw "github.com/sincin-v/collector/internal/server/middlewares/logging"
	"github.com/sincin-v/collector/internal/service"
)

func CreateRouter(service *service.MetricsService, databaseClient *db.DBClient) *chi.Mux {
	baseCtx := context.Background()
	h := handlers.New(service, databaseClient, baseCtx)
	router := chi.NewRouter()

	router.Use(logMw.LoggerMiddleware)
	router.Use(zipMw.CompressMiddleware)
	router.Post("/update/", h.UpdateMetricJSONHandler)
	router.Post("/update/{metricType}/{metricName}/{metricValue}", h.UpdateMetricHandler)
	router.Post("/value/", h.GetMetricJSONHandler)
	router.Get("/value/{metricType}/{metricName}", h.GetMetricHandler)
	router.Get("/ping", h.Ping)
	router.Get("/", h.GetAllMetricsHandler)

	return router
}
