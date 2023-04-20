package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/kontik-pk/yandex-metrics-scraper/internal/flags"
	"github.com/kontik-pk/yandex-metrics-scraper/internal/handlers"
	log "github.com/kontik-pk/yandex-metrics-scraper/internal/logger"
	"go.uber.org/zap"
	"net/http"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	log.SugarLogger = *logger.Sugar()

	params := flags.Init(flags.WithAddr())
	r := chi.NewRouter()
	r.Use(log.RequestLogger)
	r.Post("/update/{type}/{name}/{value}", handlers.SaveMetric)
	r.Get("/value/{type}/{name}", handlers.GetMetric)
	r.Get("/", handlers.ShowMetrics)

	log.SugarLogger.Infow(
		"Starting server",
		"addr", params.FlagRunAddr,
	)
	if err := http.ListenAndServe(params.FlagRunAddr, r); err != nil {
		// записываем в лог ошибку, если сервер не запустился
		log.SugarLogger.Fatalw(err.Error(), "event", "start server")
	}
}
