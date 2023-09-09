package router

import (
	"fmt"
	"github.com/kontik-pk/yandex-metrics-scraper/internal/server/middlewares/compressor"
	log "github.com/kontik-pk/yandex-metrics-scraper/internal/server/middlewares/logger"
	"github.com/kontik-pk/yandex-metrics-scraper/internal/server/router/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/kontik-pk/yandex-metrics-scraper/internal/flags"
)

func New(params flags.Params) (*chi.Mux, error) {
	handler, err := handlers.New(params.DatabaseAddress, params.Key, params.CryptoKeyPath, params.TrustedSubnet)
	if err != nil {
		return nil, fmt.Errorf("error while creating handler: %w", err)
	}

	r := chi.NewRouter()
	r.Use(log.RequestLogger)
	r.Use(compressor.Compress)
	r.Use(handler.CheckSubscription)
	r.Use(handler.CheckSubnet)
	r.Post("/update/", handler.SaveMetricFromJSON)
	r.Post("/value/", handler.GetMetricFromJSON)
	r.Post("/update/{type}/{name}/{value}", handler.SaveMetric)
	r.Get("/value/{type}/{name}", handler.GetMetric)
	r.Get("/", handler.ShowMetrics)
	r.Get("/ping", handler.Ping)
	r.Post("/updates/", handler.SaveListMetricsFromJSON)

	return r, nil
}
