package main

import (
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/kontik-pk/yandex-metrics-scraper/internal/collector"
	"github.com/kontik-pk/yandex-metrics-scraper/internal/domain"
	"golang.org/x/sync/errgroup"
	"runtime"
	"time"
)

func main() {
	metricsCollector := collector.New(&domain.MetricsStorage)

	ctx := context.Background()

	errs, _ := errgroup.WithContext(ctx)
	errs.Go(func() error {
		if err := performCollect(metricsCollector); err != nil {
			panic(err)
		}
		return nil
	})

	stick := time.NewTicker(time.Second * 10)
	client := resty.New()
	defer stick.Stop()
	errs.Go(func() error {
		if err := Send(client); err != nil {
			panic(err)
		}
		return nil
	})

	_ = errs.Wait()
}

type Icollector interface {
	Collect(metrics *runtime.MemStats)
}

func performCollect(metricsCollector Icollector) error {
	for {
		metrics := runtime.MemStats{}
		runtime.ReadMemStats(&metrics)
		metricsCollector.Collect(&metrics)
		time.Sleep(time.Second * 2)
	}
}

func Send(client *resty.Client) error {
	for {
		for n, i := range domain.MetricsStorage.Metrics {
			switch i.Value.(type) {
			case uint, uint64, int, int64:
				_, err := client.R().
					SetHeader("Content-Type", "text/plain").
					Post(fmt.Sprintf("http://localhost:8080/update/%s/%s/%d", i.MType, n, i.Value))
				if err != nil {
					return err
				}
			case float64:
				_, err := client.R().
					SetHeader("Content-Type", "text/plain").
					Post(fmt.Sprintf("http://localhost:8080/update/%s/%s/%f", i.MType, n, i.Value))
				if err != nil {
					return err
				}
			}
		}
		time.Sleep(time.Second * 10)
	}
}
