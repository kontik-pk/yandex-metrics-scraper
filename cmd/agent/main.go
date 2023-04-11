package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/kontik-pk/yandex-metrics-scraper/internal/collector"
	"golang.org/x/sync/errgroup"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"time"
)

var (
	flagRunAddr    string
	reportInterval int
	pollInterval   int
)

func parseFlags() {

	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&reportInterval, "r", 10, "report interval")
	flag.IntVar(&pollInterval, "p", 2, "poll interval")
	flag.Parse()

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		flagRunAddr = envRunAddr
	}
	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		reportIntervalEnv, err := strconv.Atoi(envReportInterval)
		if err == nil {
			reportInterval = reportIntervalEnv
		}
	}
	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		pollIntervalEnv, err := strconv.Atoi(envPollInterval)
		if err == nil {
			pollInterval = pollIntervalEnv
		}
	}
}

func main() {
	parseFlags()

	ctx := context.Background()

	errs, _ := errgroup.WithContext(ctx)
	errs.Go(func() error {
		if err := performCollect(time.Duration(pollInterval)); err != nil {
			panic(err)
		}
		return nil
	})

	client := resty.New()
	errs.Go(func() error {
		if err := send(client, reportInterval); err != nil {
			panic(err)
		}
		return nil
	})

	_ = errs.Wait()
}

func send(client *resty.Client, reportTimeout int) error {
	for {
		for n, v := range collector.Collector.GetCounters() {
			_, err := client.R().
				SetHeader("Content-Type", "text/plain").
				Post(fmt.Sprintf("http://%s/update/counter/%s/%s", flagRunAddr, n, v))
			if err != nil {
				return err
			}
		}
		for n, v := range collector.Collector.GetGauges() {
			_, err := client.R().
				SetHeader("Content-Type", "text/plain").
				Post(fmt.Sprintf("http://%s/update/gauge/%s/%s", flagRunAddr, n, v))
			if err != nil {
				return err
			}
		}
		time.Sleep(time.Second * time.Duration(reportTimeout))
	}
}

func performCollect(pollInterval time.Duration) error {
	for {
		metrics := runtime.MemStats{}
		runtime.ReadMemStats(&metrics)
		collector.Collector.Collect("Alloc", "gauge", strconv.FormatUint(metrics.Alloc, 10))
		collector.Collector.Collect("BuckHashSys", "gauge", strconv.FormatUint(metrics.BuckHashSys, 10))
		collector.Collector.Collect("Frees", "gauge", strconv.FormatUint(metrics.Frees, 10))
		collector.Collector.Collect("GCCPUFraction", "gauge", fmt.Sprintf("%.3f", metrics.GCCPUFraction))
		collector.Collector.Collect("GCSys", "gauge", strconv.FormatUint(metrics.GCSys, 10))
		collector.Collector.Collect("HeapAlloc", "gauge", strconv.FormatUint(metrics.HeapAlloc, 10))
		collector.Collector.Collect("HeapIdle", "gauge", strconv.FormatUint(metrics.HeapIdle, 10))
		collector.Collector.Collect("HeapInuse", "gauge", strconv.FormatUint(metrics.HeapInuse, 10))
		collector.Collector.Collect("HeapObjects", "gauge", strconv.FormatUint(metrics.HeapObjects, 10))
		collector.Collector.Collect("HeapReleased", "gauge", strconv.FormatUint(metrics.HeapReleased, 10))
		collector.Collector.Collect("HeapSys", "gauge", strconv.FormatUint(metrics.HeapSys, 10))
		collector.Collector.Collect("Lookups", "gauge", strconv.FormatUint(metrics.Lookups, 10))
		collector.Collector.Collect("MCacheInuse", "gauge", strconv.FormatUint(metrics.MCacheInuse, 10))
		collector.Collector.Collect("MCacheSys", "gauge", strconv.FormatUint(metrics.MCacheSys, 10))
		collector.Collector.Collect("MSpanInuse", "gauge", strconv.FormatUint(metrics.MSpanInuse, 10))
		collector.Collector.Collect("MSpanSys", "gauge", strconv.FormatUint(metrics.MSpanSys, 10))
		collector.Collector.Collect("Mallocs", "gauge", strconv.FormatUint(metrics.Mallocs, 10))
		collector.Collector.Collect("NextGC", "gauge", strconv.FormatUint(metrics.NextGC, 10))
		collector.Collector.Collect("NumForcedGC", "gauge", strconv.Itoa(int(metrics.NumForcedGC)))
		collector.Collector.Collect("NumGC", "gauge", strconv.FormatUint(uint64(metrics.NumGC), 10))
		collector.Collector.Collect("OtherSys", "gauge", strconv.Itoa(int(metrics.OtherSys)))
		collector.Collector.Collect("PauseTotalNs", "gauge", strconv.Itoa(int(metrics.PauseTotalNs)))
		collector.Collector.Collect("StackInuse", "gauge", strconv.Itoa(int(metrics.StackInuse)))
		collector.Collector.Collect("StackSys", "gauge", strconv.Itoa(int(metrics.StackSys)))
		collector.Collector.Collect("Sys", "gauge", strconv.Itoa(int(metrics.Sys)))
		collector.Collector.Collect("TotalAlloc", "gauge", strconv.Itoa(int(metrics.TotalAlloc)))
		collector.Collector.Collect("RandomValue", "gauge", strconv.Itoa(rand.Int()))

		cnt, _ := collector.Collector.GetMetric("PollCount", "counter")
		v, _ := strconv.Atoi(cnt)
		collector.Collector.Collect("PollCount", "counter", strconv.Itoa(v+1))

		time.Sleep(time.Second * pollInterval)
	}
}
