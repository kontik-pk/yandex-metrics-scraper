package main

import (
	"bytes"
	"compress/gzip"
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
	req := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip")
	for {
		for n, v := range collector.Collector.GetCounters() {
			jsonInput := fmt.Sprintf(`{"id":%q, "type":"counter", "delta": %s}`, n, v)
			if err := sendRequest(req, jsonInput, addr); err != nil {
				return fmt.Errorf("error while sending agent request for counter metric: %w", err)
			}
		}
		for n, v := range collector.Collector.GetGauges() {
			jsonInput := fmt.Sprintf(`{"id":%q, "type":"gauge", "value": %s}`, n, v)
			if err := sendRequest(req, jsonInput, addr); err != nil {
				return fmt.Errorf("error while sending agent request for gauge metric: %w", err)
			}
		}
		time.Sleep(time.Duration(reportTimeout) * time.Second)
	}
}

func sendRequest(req *resty.Request, jsonInput string, addr string) error {
	buf := bytes.NewBuffer(nil)
	zb := gzip.NewWriter(buf)
	if _, err := zb.Write([]byte(jsonInput)); err != nil {
		return err
	}
	if err := zb.Close(); err != nil {
		return fmt.Errorf("error while trying to close writer: %w", err)
	}

	err := retry.Do(
		func() error {
			var err error
			if _, err = req.SetBody(buf).Post(fmt.Sprintf("http://%s/update/", addr)); err != nil {
				return fmt.Errorf("error while trying to create post request: %w", err)
			}
			return nil
		},
		retry.Attempts(10),
		retry.OnRetry(func(n uint, err error) {
			log.Printf("Retrying request after error: %v", err)
		}),
	)
	if err != nil {
		return fmt.Errorf("error while trying to connect to server: %w", err)
	}
	// do something with the response
	return nil
}
