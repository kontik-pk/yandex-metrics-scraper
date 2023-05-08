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
		flags.WithDatabase(),
	)

	r := router.New(*params)

	log.SugarLogger.Infow(
		"Starting server",
		"addr", params.FlagRunAddr,
	)

	// init restorer
	var saver saver
	if params.FileStoragePath != "" {
		saver = file.New(params)
	} else if params.DatabaseAddress != "" {
		saver, err = database.New(params)
		if err != nil {
			log.SugarLogger.Errorf(err.Error())
		}
	}

	// restore previous metrics if needed
	ctx := context.Background()
	if params.Restore {
		metrics, err := saver.Restore(ctx)
		if err != nil {
			log.SugarLogger.Error(err.Error(), "restore error")
		}
		collector.Collector.Metrics = metrics
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
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := saver.Save(ctx, collector.Collector.Metrics); err != nil {
				log.SugarLogger.Error(err.Error(), "save error")
			} else {
				log.SugarLogger.Info("successfully saved metrics")
			}
		}
	}
}

type saver interface {
	Restore(ctx context.Context) ([]collector.MetricJSON, error)
	Save(ctx context.Context, metrics []collector.MetricJSON) error
}
