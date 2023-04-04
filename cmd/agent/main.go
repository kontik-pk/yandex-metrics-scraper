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

var m = domain.MemStorage{Metrics: map[string]domain.Metric{}}

func main() {
	metricsCollector := collector.New(&m)

	ctx := context.Background()

	mtick := time.NewTicker(time.Second * 2)
	defer mtick.Stop()

	errs, ctx := errgroup.WithContext(ctx)
	errs.Go(func() error {
		if err := performCollect(ctx, mtick, metricsCollector); err != nil {
			panic(err)
		}
		return nil
	})

	stick := time.NewTicker(time.Second * 10)
	client := resty.New()
	defer stick.Stop()
	errs.Go(func() error {
		if err := Send(ctx, stick, client); err != nil {
			panic(err)
		}
		return nil
	})

	_ = errs.Wait()
}

type Icollector interface {
	Collect(metrics *runtime.MemStats)
}

func performCollect(ctx context.Context, ticker *time.Ticker, metricsCollector Icollector) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			metrics := runtime.MemStats{}
			runtime.ReadMemStats(&metrics)
			metricsCollector.Collect(&metrics)
		}
	}
}

func Send(ctx context.Context, ticker *time.Ticker, client *resty.Client) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			for n, i := range m.Metrics {
				switch i.Value.(type) {
				case uint, uint64:
					resp, err := client.R().
						SetHeader("Content-Type", "text/plain").
						Post(fmt.Sprintf("http://localhost:8080/update/%s/%s/%d", i.MType, n, i.Value))
					if err != nil {
						return err
					}
					fmt.Println(resp.Status())
				case float64:
					resp, err := client.R().
						SetHeader("Content-Type", "text/plain").
						Post(fmt.Sprintf("http://localhost:8080/update/%s/%s/%f", i.MType, n, i.Value))
					if err != nil {
						return err
					}
					fmt.Println(resp.Status())
				}
			}
		}
	}
}
