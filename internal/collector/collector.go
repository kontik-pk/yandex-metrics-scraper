package collector

import (
	"github.com/kontik-pk/yandex-metrics-scraper/internal/domain"
	"math/rand"
	"runtime"
)

func (c *collector) Collect(metrics *runtime.MemStats) {
	c.storage.Metrics["Alloc"] = domain.Metric{Value: metrics.Alloc, MType: "gauge"}
	c.storage.Metrics["BuckHashSys"] = domain.Metric{Value: metrics.BuckHashSys, MType: "gauge"}
	c.storage.Metrics["Frees"] = domain.Metric{Value: metrics.Frees, MType: "gauge"}
	c.storage.Metrics["GCCPUFraction"] = domain.Metric{Value: metrics.GCCPUFraction, MType: "gauge"}
	c.storage.Metrics["GCSys"] = domain.Metric{Value: metrics.GCSys, MType: "gauge"}
	c.storage.Metrics["HeapAlloc"] = domain.Metric{Value: metrics.HeapAlloc, MType: "gauge"}
	c.storage.Metrics["HeapIdle"] = domain.Metric{Value: metrics.HeapIdle, MType: "gauge"}
	c.storage.Metrics["HeapInuse"] = domain.Metric{Value: metrics.HeapInuse, MType: "gauge"}
	c.storage.Metrics["HeapObjects"] = domain.Metric{Value: metrics.HeapObjects, MType: "gauge"}
	c.storage.Metrics["HeapReleased"] = domain.Metric{Value: metrics.HeapReleased, MType: "gauge"}
	c.storage.Metrics["HeapSys"] = domain.Metric{Value: metrics.HeapSys, MType: "gauge"}
	c.storage.Metrics["Lookups"] = domain.Metric{Value: metrics.Lookups, MType: "gauge"}
	c.storage.Metrics["MCacheInuse"] = domain.Metric{Value: metrics.MCacheInuse, MType: "gauge"}
	c.storage.Metrics["MCacheSys"] = domain.Metric{Value: metrics.MCacheSys, MType: "gauge"}
	c.storage.Metrics["MSpanInuse"] = domain.Metric{Value: metrics.MSpanInuse, MType: "gauge"}
	c.storage.Metrics["MSpanSys"] = domain.Metric{Value: metrics.MSpanSys, MType: "gauge"}
	c.storage.Metrics["Mallocs"] = domain.Metric{Value: metrics.Mallocs, MType: "gauge"}
	c.storage.Metrics["NextGC"] = domain.Metric{Value: metrics.NextGC, MType: "gauge"}
	c.storage.Metrics["NumForcedGC"] = domain.Metric{Value: metrics.NumForcedGC, MType: "gauge"}
	c.storage.Metrics["NumGC"] = domain.Metric{Value: metrics.NumGC, MType: "gauge"}
	c.storage.Metrics["OtherSys"] = domain.Metric{Value: metrics.OtherSys, MType: "gauge"}
	c.storage.Metrics["PauseTotalNs"] = domain.Metric{Value: metrics.PauseTotalNs, MType: "gauge"}
	c.storage.Metrics["StackInuse"] = domain.Metric{Value: metrics.StackInuse, MType: "gauge"}
	c.storage.Metrics["StackSys"] = domain.Metric{Value: metrics.StackSys, MType: "gauge"}
	c.storage.Metrics["Sys"] = domain.Metric{Value: metrics.Sys, MType: "gauge"}
	c.storage.Metrics["TotalAlloc"] = domain.Metric{Value: metrics.TotalAlloc, MType: "gauge"}
	c.storage.Metrics["RandomValue"] = domain.Metric{Value: rand.Int(), MType: "gauge"}

	var cnt int64
	if c.storage.Metrics["PollCount"].Value != nil {
		cnt = c.storage.Metrics["PollCount"].Value.(int64) + 1
	}
	c.storage.Metrics["PollCount"] = domain.Metric{Value: cnt, MType: "counter"}
}

func New(ms *domain.MemStorage) *collector {
	return &collector{ms}
}

type collector struct {
	storage *domain.MemStorage
}
