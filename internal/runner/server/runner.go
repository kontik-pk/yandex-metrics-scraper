package server

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/kontik-pk/yandex-metrics-scraper/internal/collector"
	"github.com/kontik-pk/yandex-metrics-scraper/internal/flags"
	log "github.com/kontik-pk/yandex-metrics-scraper/internal/logger"
	"github.com/kontik-pk/yandex-metrics-scraper/internal/router/router"
	"github.com/kontik-pk/yandex-metrics-scraper/internal/saver/database"
	"github.com/kontik-pk/yandex-metrics-scraper/internal/saver/file"
	"go.uber.org/zap"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
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
	isRestore       bool
	storeInterval   int
	tlsKey          string
	appSrv          server
	pprofSrv        server
	signals         chan os.Signal
}

func New(params *flags.Params) *Runner {
	// init saver (file or db)
	saver, err := initSaver(params)
	if err != nil {
		log.SugarLogger.Fatalw(err.Error(), "error", "init metrics saver")
	}
	// init router
	r, err := router.New(*params)
	if err != nil {
		log.SugarLogger.Fatalw(err.Error(), "error", "creating router")
	}
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	return &Runner{
		saver:           saver,
		metricsInterval: time.Duration(params.StoreInterval),
		isRestore:       params.Restore,
		storeInterval:   params.StoreInterval,
		tlsKey:          params.CryptoKeyPath,
		appSrv: &http.Server{
			Addr:    params.FlagRunAddr,
			Handler: r,
		},
		pprofSrv: &http.Server{
			Addr:    pprofAddr,
			Handler: nil,
		},
		signals: sigs,
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
		if err := r.pprofSrv.ListenAndServe(); err != nil {
			log.SugarLogger.Fatalw(err.Error(), "pprof", "start pprof server")
		}
	}()

	// run server
	log.SugarLogger.Info(
		"Starting server",
	)
	go func() {
		if err := r.appSrv.ListenAndServe(); err != nil {
			log.SugarLogger.Fatalw(err.Error(), "event", "start server")
		}
	}()

	sig := <-r.signals
	log.SugarLogger.Info(fmt.Sprintf("got signal: %s", sig.String()))
	if err = r.saver.Save(ctx, collector.Collector.Metrics); err != nil {
		log.SugarLogger.Error(err.Error(), "save error")
	} else {
		log.SugarLogger.Info("metrics was successfully saved")
	}
	if err = r.appSrv.Shutdown(ctx); err != nil {
		log.SugarLogger.Error(fmt.Sprintf("error while server shutdown: %s", err.Error()), "server shutdown error")
	}
}

func (r *Runner) saveMetrics(ctx context.Context, interval int) {
	ticker := time.NewTicker(time.Duration(interval))
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

//go:generate mockery --inpackage --disable-version-string --filename saver_mock.go --name saver
type saver interface {
	Restore(ctx context.Context) ([]collector.StoredMetric, error)
	Save(ctx context.Context, metrics []collector.StoredMetric) error
}

//go:generate mockery --inpackage --disable-version-string --filename server_mock.go --name server
type server interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}
