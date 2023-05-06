package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/kontik-pk/yandex-metrics-scraper/internal/compressor"
	log "github.com/kontik-pk/yandex-metrics-scraper/internal/logger"
	"github.com/kontik-pk/yandex-metrics-scraper/internal/router/handlers"
)

func New() *chi.Mux {
	r := chi.NewRouter()
	r.Use(log.RequestLogger)
	r.Use(compressor.Compress)
	r.Post("/update/", handlers.SaveMetricFromJSON)
	r.Post("/value/", handlers.GetMetricFromJSON)
	r.Post("/update/{type}/{name}/{value}", handlers.SaveMetric)
	r.Get("/value/{type}/{name}", handlers.GetMetric)
	r.Get("/", handlers.ShowMetrics)

	return r
}
