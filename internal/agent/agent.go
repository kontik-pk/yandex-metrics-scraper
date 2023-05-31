package agent

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/avast/retry-go"
	"github.com/go-resty/resty/v2"
	"github.com/kontik-pk/yandex-metrics-scraper/internal/collector"
	"github.com/kontik-pk/yandex-metrics-scraper/internal/flags"
	aggregator "github.com/kontik-pk/yandex-metrics-scraper/internal/metrics"
	"go.uber.org/zap"
	"log"
	"time"
)

func (a *Agent) CollectMetrics(ctx context.Context) (err error) {
	aggTicker := time.NewTicker(time.Duration(a.params.PollInterval) * time.Second)
	go func() {
		for {
			select {
			case <-ctx.Done():
				err = ctx.Err()
				return
			case <-aggTicker.C:
				a.aggregator.AggregateRuntimeMetrics()
			}
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				err = ctx.Err()
				return
			case <-aggTicker.C:
				a.aggregator.AggregateGopsutilMetrics()
			}
		}
	}()

	return err
}

func (a *Agent) SendMetrics(ctx context.Context) error {
	numRequests := make(chan struct{}, a.params.RateLimit)
	reportTicker := time.NewTicker(time.Duration(a.params.ReportInterval) * time.Second)
	client := resty.New()
	for {
		select {
		case <-ctx.Done():
			return nil
		// check if time to send metrics on server
		case <-reportTicker.C:
			select {
			case <-ctx.Done():
				return nil
			// check if the rate limit is exceeded
			case numRequests <- struct{}{}:
				a.log.Info("metrics sent on server")
				if err := a.sendMetrics(client); err != nil {
					return err
				}
			default:
				a.log.Info("rate limit is exceeded")
			}
		}
		// release the request pool for one
		<-numRequests
	}
}

func (a *Agent) sendMetrics(client *resty.Client) error {
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
		if a.params.Key != "" {
			//h := hmac.New(sha256.New, []byte(key))
			//h.Write(jsonInput)
			//dst := h.Sum(nil)
			dst := sha256.Sum256(jsonInput)
			req.SetHeader("HashSHA256", fmt.Sprintf("%x", dst))
		}
		if err := a.sendRequestsWithRetries(req, string(jsonInput)); err != nil {
			return fmt.Errorf("error while sending agent request for counter metric: %w", err)
		}
	}
	return nil
}

func (a *Agent) sendRequestsWithRetries(req *resty.Request, jsonInput string) error {
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
			if _, err := req.SetBody(buf).Post(fmt.Sprintf("http://%s/update/", a.params.FlagRunAddr)); err != nil {
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

func New(params *flags.Params, aggregator *aggregator.Aggregator, log zap.SugaredLogger) *Agent {
	return &Agent{
		params:     params,
		aggregator: aggregator,
		log:        log,
	}
}

type Agent struct {
	params     *flags.Params
	aggregator *aggregator.Aggregator
	log        zap.SugaredLogger
}
