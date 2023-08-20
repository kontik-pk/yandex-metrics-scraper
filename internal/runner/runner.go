package runner

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/kontik-pk/yandex-metrics-scraper/internal/collector"
	"github.com/kontik-pk/yandex-metrics-scraper/internal/flags"
	log "github.com/kontik-pk/yandex-metrics-scraper/internal/logger"
	"github.com/kontik-pk/yandex-metrics-scraper/internal/router/router"
	"github.com/kontik-pk/yandex-metrics-scraper/internal/saver/database"
	"github.com/kontik-pk/yandex-metrics-scraper/internal/saver/file"
	"go.uber.org/zap"
	"net/http"
	_ "net/http/pprof"
	"time"
)

const (
	pprofAddr    string = ":90"
	buildVersion string = ""
	buildDate    string = ""
	buildCommit  string = ""
)

type Runner struct {
	saver           saver
	metricsInterval time.Duration
	router          *chi.Mux
	isRestore       bool
	storeInterval   int
	runAddress      string
	tlsKey          string
}

func New(params *flags.Params) *Runner {
	// init restorer
	saver, err := initSaver(params)
	if err != nil {
		log.SugarLogger.Fatalw(err.Error(), "error", "init metrics saver")
	}
	r, err := router.New(*params)
	if err != nil {
		log.SugarLogger.Fatalw(err.Error(), "error", "creating router")
	}
	return &Runner{
		saver:           saver,
		metricsInterval: time.Duration(params.StoreInterval),
		router:          r,
		isRestore:       params.Restore,
		storeInterval:   params.StoreInterval,
		runAddress:      params.FlagRunAddr,
		tlsKey:          params.CryptoKeyPath,
	}
}

func (r *Runner) Run(ctx context.Context) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		fmt.Println("error while creating logger, exit")
		return
	}
	defer logger.Sync()
	log.SugarLogger = *logger.Sugar()

	log.SugarLogger.Info("Build version: %s\nBuild date: %s\nBuild commit: %s\n", buildVersion, buildDate, buildCommit)

	// restore previous metrics if needed
	if r.isRestore {
		metrics, err := r.saver.Restore(ctx)
		if err != nil {
			log.SugarLogger.Error(err.Error(), "restore error")
		}
		collector.Collector.Metrics = metrics
		log.SugarLogger.Info("metrics restored")
	}

	// regularly save metrics
	go r.saveMetrics(ctx, r.storeInterval)

	// pprof
	go func() {
		if err := http.ListenAndServe(pprofAddr, nil); err != nil {
			log.SugarLogger.Fatalw(err.Error(), "pprof", "start pprof server")
		}
	}()

	// run server
	log.SugarLogger.Infow(
		"Starting server",
		"addr", r.runAddress,
	)
	if err := http.ListenAndServe(r.runAddress, r.router); err != nil {
		log.SugarLogger.Fatalw(err.Error(), "event", "start server")
	}
}

func (r *Runner) runServer() {
	if err := http.ListenAndServe(r.runAddress, r.router); err != nil {
		log.SugarLogger.Fatalw(err.Error(), "event", "start server")
	}
}

func (r *Runner) saveMetrics(ctx context.Context, interval int) {
	ticker := time.NewTicker(time.Duration(interval))
	if err := r.saver.Save(ctx, collector.Collector.Metrics); err != nil {
		log.SugarLogger.Error(err.Error(), "save error")
	}
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := r.saver.Save(ctx, collector.Collector.Metrics); err != nil {
				log.SugarLogger.Error(err.Error(), "save error")
			}
		}
	}
}

func initSaver(params *flags.Params) (saver, error) {
	if params.DatabaseAddress != "" {
		db, err := sql.Open("pgx", params.DatabaseAddress)
		if err != nil {
			return nil, err
		}
		return database.New(db)
	} else if params.FileStoragePath != "" {
		return file.New(params.FileStoragePath), nil
	}
	return nil, fmt.Errorf("neither file path nor database address was specified")
}

type saver interface {
	Restore(ctx context.Context) ([]collector.StoredMetric, error)
	Save(ctx context.Context, metrics []collector.StoredMetric) error
}
