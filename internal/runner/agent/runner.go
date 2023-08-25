package agent

import (
	"context"
	"fmt"
	agent2 "github.com/kontik-pk/yandex-metrics-scraper/internal/agent"
	"github.com/kontik-pk/yandex-metrics-scraper/internal/collector"
	"github.com/kontik-pk/yandex-metrics-scraper/internal/flags"
	log "github.com/kontik-pk/yandex-metrics-scraper/internal/logger"
	aggregator "github.com/kontik-pk/yandex-metrics-scraper/internal/metrics"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type Runner struct {
	params *flags.Params
}

func New(params *flags.Params) *Runner {
	return &Runner{
		params: params,
	}
}

func (r *Runner) Run(ctx context.Context) {
	errs, ctx := errgroup.WithContext(ctx)

	logger, err := zap.NewDevelopment()
	if err != nil {
		fmt.Println("error while creating logger, exit")
		return
	}
	defer logger.Sync()
	log.SugarLogger = *logger.Sugar()

	agent, err := agent2.New(r.params, aggregator.New(&collector.Collector), log.SugarLogger)
	if err != nil {
		log.SugarLogger.Fatalw(err.Error(), "error", "creating agent")
	}
	errs.Go(func() error {
		return agent.CollectMetrics(ctx)
	})
	errs.Go(func() error {
		return agent.SendMetrics(ctx)
	})
	if err = errs.Wait(); err != nil {
		log.SugarLogger.Errorf(fmt.Sprintf("error while running agent: %s", err.Error()))
	}
}
