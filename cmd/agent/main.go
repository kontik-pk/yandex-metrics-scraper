package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
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

	aggTicker := time.NewTicker(time.Duration(params.PollInterval) * time.Second)
	errs, ctx := errgroup.WithContext(context.Background())
	errs.Go(func() error {
		agg := aggregator.New(&collector.Collector)
		for {
			select {
			case <-ctx.Done():
				return nil
			case <-aggTicker.C:
				agg.Aggregate()
			}
		}
	})

	reportTicker := time.NewTicker(time.Duration(params.ReportInterval) * time.Second)
	client := resty.New()
	errs.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return nil
			case <-reportTicker.C:
				if err := sendMetrics(client, params.FlagRunAddr); err != nil {
					return err
				}
			}
		}
	})

	_ = errs.Wait()
}

func sendMetrics(client *resty.Client, addr string) error {
	req := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept-Encoding", "gzip").
		SetHeader("Content-Encoding", "gzip")

	for _, v := range collector.Collector.Metrics {
		jsonInput, _ := json.Marshal(collector.MetricRequest{
			ID:    v.ID,
			MType: v.MType,
			Delta: v.CounterValue,
			Value: v.GaugeValue,
		})
		if err := sendRequestsWithRetries(req, string(jsonInput), addr); err != nil {
			return fmt.Errorf("error while sending agent request for counter metric: %w", err)
		}
	}
	return nil
}

func sendRequestsWithRetries(req *resty.Request, jsonInput string, addr string) error {
	buf := bytes.NewBuffer(nil)
	zb := gzip.NewWriter(buf)
	if _, err := zb.Write([]byte(jsonInput)); err != nil {
		return fmt.Errorf("error while write json input: %w", err)
	}
	if err := zb.Close(); err != nil {
		return fmt.Errorf("error while trying to close writer: %w", err)
	}

	if err := retry.Do(
		func() error {
			if _, err := req.SetBody(buf).Post(fmt.Sprintf("http://%s/update/", addr)); err != nil {
				return fmt.Errorf("error while trying to create post request: %w", err)
			}
			return nil
		},
		retry.Attempts(10),
		retry.OnRetry(func(n uint, err error) {
			log.Printf("Retrying request after error: %v", err)
		}),
	); err != nil {
		return fmt.Errorf("error while trying to connect to server: %w", err)
	}
	return nil
}
