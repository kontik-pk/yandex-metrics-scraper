package collector

import (
	"github.com/kontik-pk/yandex-metrics-scraper/internal/domain"
	"math/rand"
	"runtime"
)

func (c *collector) Collect(metrics runtime.MemStats) {
	c.storage.Metrics["Alloc"] = domain.Metric{Value: metrics.Alloc, MType: "gauge"}
	c.storage.Metrics["BuckHashSys"] = domain.Metric{Value: metrics.BuckHashSys, MType: "gauge"}
	c.storage.Metrics["Frees"] = domain.Metric{metrics.Frees, "gauge"}
	c.storage.Metrics["GCCPUFraction"] = domain.Metric{metrics.GCCPUFraction, "gauge"}
	c.storage.Metrics["GCSys"] = domain.Metric{metrics.GCSys, "gauge"}
	c.storage.Metrics["HeapAlloc"] = domain.Metric{metrics.HeapAlloc, "gauge"}
	c.storage.Metrics["HeapIdle"] = domain.Metric{metrics.HeapIdle, "gauge"}
	c.storage.Metrics["HeapInuse"] = domain.Metric{metrics.HeapInuse, "gauge"}
	c.storage.Metrics["HeapObjects"] = domain.Metric{metrics.HeapObjects, "gauge"}
	c.storage.Metrics["HeapReleased"] = domain.Metric{metrics.HeapReleased, "gauge"}
	c.storage.Metrics["HeapSys"] = domain.Metric{metrics.HeapSys, "gauge"}
	c.storage.Metrics["Lookups"] = domain.Metric{metrics.Lookups, "gauge"}
	c.storage.Metrics["MCacheInuse"] = domain.Metric{metrics.MCacheInuse, "gauge"}
	c.storage.Metrics["MCacheSys"] = domain.Metric{metrics.MCacheSys, "gauge"}
	c.storage.Metrics["MSpanInuse"] = domain.Metric{metrics.MSpanInuse, "gauge"}
	c.storage.Metrics["MSpanSys"] = domain.Metric{metrics.MSpanSys, "gauge"}
	c.storage.Metrics["Mallocs"] = domain.Metric{metrics.Mallocs, "gauge"}
	c.storage.Metrics["NextGC"] = domain.Metric{metrics.NextGC, "gauge"}
	c.storage.Metrics["NumForcedGC"] = domain.Metric{metrics.NumForcedGC, "gauge"}
	c.storage.Metrics["NumGC"] = domain.Metric{metrics.NumGC, "gauge"}
	c.storage.Metrics["OtherSys"] = domain.Metric{metrics.OtherSys, "gauge"}
	c.storage.Metrics["PauseTotalNs"] = domain.Metric{metrics.PauseTotalNs, "gauge"}
	c.storage.Metrics["StackInuse"] = domain.Metric{metrics.StackInuse, "gauge"}
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

func New(ms domain.MemStorage) *collector {
	return &collector{ms}
}

type collector struct {
	storage domain.MemStorage
}
