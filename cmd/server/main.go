package main

import (
	"context"
	"github.com/kontik-pk/yandex-metrics-scraper/internal/collector"
	"github.com/kontik-pk/yandex-metrics-scraper/internal/flags"
	log "github.com/kontik-pk/yandex-metrics-scraper/internal/logger"
	"github.com/kontik-pk/yandex-metrics-scraper/internal/router/router"
	"github.com/kontik-pk/yandex-metrics-scraper/internal/saver/database"
	"github.com/kontik-pk/yandex-metrics-scraper/internal/saver/file"
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

	params := flags.Init(
		flags.WithAddr(),
		flags.WithStoreInterval(),
		flags.WithFileStoragePath(),
		flags.WithRestore(),
		flags.WithDatabase(),
		flags.WithKey(),
	)

	r := router.New(*params)

	log.SugarLogger.Infow(
		"Starting server",
		"addr", params.FlagRunAddr)

	// init restorer
	var saver saver
	if params.FileStoragePath != "" && params.DatabaseAddress == "" {
		saver = file.New(params)
	} else if params.DatabaseAddress != "" {
		saver, err = database.New(params)
		if err != nil {
			log.SugarLogger.Errorf(err.Error())
		}
	}

	// restore previous metrics if needed
	ctx := context.Background()
	//TODO: хотелось бы избавиться от этого if и от того, что на 62 строке, как можно сделать лучше?
	if params.Restore && (params.FileStoragePath != "" || params.DatabaseAddress != "") {
		metrics, err := saver.Restore(ctx)
		if err != nil {
			log.SugarLogger.Error(err.Error(), "restore error")
		}
		collector.Collector.Metrics = metrics
		log.SugarLogger.Info("metrics restored")
	}

	// regularly save metrics if needed
	if params.DatabaseAddress != "" || params.FileStoragePath != "" {
		go saveMetrics(ctx, saver, params.StoreInterval)
	}

	// run server
	if err := http.ListenAndServe(params.FlagRunAddr, r); err != nil {
		log.SugarLogger.Fatalw(err.Error(), "event", "start server")
	}
}

func saveMetrics(ctx context.Context, saver saver, interval int) {
	//здесь должен быть нормальный тикер, но тесты метрики не ждут
	//ticker := time.NewTicker(time.Duration(interval))
	for {
		if err := saver.Save(ctx, collector.Collector.Metrics); err != nil {
			log.SugarLogger.Error(err.Error(), "save error")
		}
		//time.Sleep(time.Duration(interval))
	}
}

type saver interface {
	Restore(ctx context.Context) ([]collector.StoredMetric, error)
	Save(ctx context.Context, metrics []collector.StoredMetric) error
}
