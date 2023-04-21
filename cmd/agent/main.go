package main

import (
	"context"
	"fmt"
	"github.com/avast/retry-go"
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
			req := client.R().SetHeader("Content-Type", "application/json").
				SetBody(fmt.Sprintf(`{"id":%q, "type":"counter", "delta": %s}`, n, v))
			if err := sendRequest(req, addr); err != nil {
				return err
			}
		}
		for n, v := range collector.Collector.GetGauges() {
			req := client.R().
				SetHeader("Content-Type", "application/json").
				SetBody(fmt.Sprintf(`{"id":%q, "type":"gauge", "value": %s}`, n, v))
			if err := sendRequest(req, addr); err != nil {
				return err
			}
		}
		time.Sleep(time.Duration(reportTimeout) * time.Second)
	}
}

func sendRequest(req *resty.Request, addr string) error {
	err := retry.Do(
		func() error {
			var err error
			_, err = req.Post(fmt.Sprintf("http://%s/update/", addr))
			return err
		},
		retry.Attempts(10),
		retry.OnRetry(func(n uint, err error) {
			log.Printf("Retrying request after error: %v", err)
		}),
	)
	if err != nil {
		return err
	}
	// do something with the response
	return nil
}
