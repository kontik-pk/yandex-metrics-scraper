package main

import (
	"github.com/kontik-pk/yandex-metrics-scraper/internal/collector"
	"github.com/kontik-pk/yandex-metrics-scraper/internal/flags"
	log "github.com/kontik-pk/yandex-metrics-scraper/internal/logger"
	"github.com/kontik-pk/yandex-metrics-scraper/internal/router/router"
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

	r := router.New()

	log.SugarLogger.Infow(
		"Starting server",
		"addr", params.FlagRunAddr,
	)
	// restore previous metrics if needed
	if params.Restore {
		if err := collector.Collector.Restore(params.FileStoragePath); err != nil {
			log.SugarLogger.Error(err.Error(), "restore error")
		}
	}

	// regularly save metrics if needed
	if params.FileStoragePath != "" {
		go saveMetrics(params.FileStoragePath, params.StoreInterval)
	}

	// run server
	if err := http.ListenAndServe(params.FlagRunAddr, r); err != nil {
		log.SugarLogger.Fatalw(err.Error(), "event", "start server")
	}
}

func saveMetrics(path string, interval int) {
	for {
		if err := collector.Collector.Save(path); err != nil {
			log.SugarLogger.Error(err.Error(), "save error")
		} else {
			log.SugarLogger.Info("successfully saved metrics")
		}
		time.Sleep(time.Duration(interval) * time.Second)
	}
}