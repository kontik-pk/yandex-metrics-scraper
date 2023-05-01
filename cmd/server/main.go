package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/kontik-pk/yandex-metrics-scraper/internal/collector"
	"github.com/kontik-pk/yandex-metrics-scraper/internal/compressor"
	"github.com/kontik-pk/yandex-metrics-scraper/internal/flags"
	"github.com/kontik-pk/yandex-metrics-scraper/internal/handlers"
	log "github.com/kontik-pk/yandex-metrics-scraper/internal/logger"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	log.SugarLogger = *logger.Sugar()

	params := flags.Init(
		flags.WithAddr(),
		flags.WithStoreInterval(),
		flags.WithFileStoragePath(),
		flags.WithRestore(),
	)
	r := chi.NewRouter()
	r.Use(log.RequestLogger)
	r.Use(compressor.Compress)
	r.Post("/update/", handlers.SaveMetricFromJSON)
	r.Post("/value/", handlers.GetMetricFromJSON)
	r.Post("/update/{type}/{name}/{value}", handlers.SaveMetric)
	r.Get("/value/{type}/{name}", handlers.GetMetric)
	r.Get("/", handlers.ShowMetrics)

	log.SugarLogger.Infow(
		"Starting server",
		"addr", params.FlagRunAddr,
	)
	if params.Restore {
		if err := collector.Collector.Restore(params.FileStoragePath); err != nil {
			log.SugarLogger.Error(err.Error(), "restore error")
		}
	}

	go func() {
		if params.FileStoragePath != "" {
			for {
				err = collector.Collector.Save(params.FileStoragePath)
				if err != nil {
					log.SugarLogger.Error(err.Error(), "save error")
				} else {
					log.SugarLogger.Info("successfully saved metrics")
				}
				time.Sleep(time.Duration(params.StoreInterval) * time.Second)
			}
		}
	}()
	if err := http.ListenAndServe(params.FlagRunAddr, r); err != nil {
		log.SugarLogger.Fatalw(err.Error(), "event", "start server")
	}
}
