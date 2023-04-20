package main

import (
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/kontik-pk/yandex-metrics-scraper/internal/collector"
	"github.com/kontik-pk/yandex-metrics-scraper/internal/flags"
	aggregator "github.com/kontik-pk/yandex-metrics-scraper/internal/metrics"
	"golang.org/x/sync/errgroup"
	"log"
	"time"
)

func main() {
	params := flags.Init(flags.WithPollInterval(), flags.WithReportInterval(), flags.WithAddr())

	ctx := context.Background()

	errs, _ := errgroup.WithContext(ctx)
	errs.Go(func() error {
		agg := aggregator.New(&collector.Collector)
		for {
			agg.Aggregate()
			time.Sleep(time.Duration(params.PollInterval) * time.Second)
		}
	})

	client := resty.New()
	errs.Go(func() error {
		if err := send(client, params.ReportInterval, params.FlagRunAddr); err != nil {
			log.Fatalln(err)
		}
		return nil
	})

	_ = errs.Wait()
}

func send(client *resty.Client, reportTimeout int, addr string) error {
	for {
		for n, v := range collector.Collector.GetCounters() {
			if _, err := client.R().
				SetHeader("Content-Type", "text/plain").
				Post(fmt.Sprintf("http://%s/update/counter/%s/%s", addr, n, v)); err != nil {
				return err
			}
		}
		for n, v := range collector.Collector.GetGauges() {
			if _, err := client.R().
				SetHeader("Content-Type", "text/plain").
				Post(fmt.Sprintf("http://%s/update/gauge/%s/%s", addr, n, v)); err != nil {
				return err
			}
		}
		time.Sleep(time.Duration(reportTimeout) * time.Second)
	}
}
