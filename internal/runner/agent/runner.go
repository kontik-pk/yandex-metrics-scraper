package agent

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	agent2 "github.com/kontik-pk/yandex-metrics-scraper/internal/agent"
	"github.com/kontik-pk/yandex-metrics-scraper/internal/collector"
	"github.com/kontik-pk/yandex-metrics-scraper/internal/flags"
	aggregator "github.com/kontik-pk/yandex-metrics-scraper/internal/metrics"
	"go.uber.org/zap"
)

type Runner struct {
	params  *flags.Params
	logger  *zap.SugaredLogger
	signals chan os.Signal
}

func New(params *flags.Params) *Runner {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	logger, err := zap.NewDevelopment()
	if err != nil {
		fmt.Println("error while creating logger, exit")
		return nil
	}
	defer logger.Sync()

	return &Runner{
		params:  params,
		signals: sigs,
		logger:  logger.Sugar(),
	}
}

func (r *Runner) Run(ctx context.Context) {
	runCtx, cancel := context.WithCancel(ctx)

	var wg sync.WaitGroup

	// init agent
	agent, err := agent2.New(r.params, aggregator.New(collector.Collector()), r.logger)
	if err != nil {
		r.logger.Fatalw(err.Error(), "error", "creating agent")
	}

	// collect all necessary metrics
	wg.Add(1)
	go func() {
		agent.CollectMetrics(runCtx)
		wg.Done()
	}()

	// send metrics on server by timer internally
	wg.Add(1)
	go func() {
		if err = agent.SendMetricsLoop(runCtx); err != nil {
			r.logger.Errorf("send metrics loop exited with error: %s", err.Error())
			wg.Done()
			cancel()
		}
		wg.Done()
	}()

	// catch signals
	wg.Add(1)
	go func() {
		sig := <-r.signals
		r.logger.Info(fmt.Sprintf("got signal: %s", sig.String()))
		if err = agent.SendMetrics(runCtx); err != nil {
			r.logger.Errorf("send metrics after signal %q exited with error: %s", sig.String(), err.Error())
		} else {
			r.logger.Infof("metrics successfully sent after signal %q", sig.String())
		}
		cancel()
		wg.Done()
	}()

	// wait for all goroutines to complete
	wg.Wait()
}
