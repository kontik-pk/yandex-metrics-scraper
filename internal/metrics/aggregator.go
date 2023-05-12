package aggregator

import (
	"github.com/kontik-pk/yandex-metrics-scraper/internal/collector"
	"math/rand"
	"runtime"
)

func (a *Aggregator) Aggregate() {
	metrics := runtime.MemStats{}
	runtime.ReadMemStats(&metrics)

	a.metricsCollector.Collect(collector.MetricJSON{ID: "Alloc", MType: "gauge", Value: PtrFloat64(float64(metrics.Alloc))})
	a.metricsCollector.Collect(collector.MetricJSON{ID: "BuckHashSys", MType: "gauge", Value: PtrFloat64(float64(metrics.BuckHashSys))})
	a.metricsCollector.Collect(collector.MetricJSON{ID: "Frees", MType: "gauge", Value: PtrFloat64(float64(metrics.Frees))})
	a.metricsCollector.Collect(collector.MetricJSON{ID: "GCCPUFraction", MType: "gauge", Value: &metrics.GCCPUFraction})
	a.metricsCollector.Collect(collector.MetricJSON{ID: "GCSys", MType: "gauge", Value: PtrFloat64(float64(metrics.GCSys))})
	a.metricsCollector.Collect(collector.MetricJSON{ID: "HeapAlloc", MType: "gauge", Value: PtrFloat64(float64(metrics.HeapAlloc))})
	a.metricsCollector.Collect(collector.MetricJSON{ID: "HeapIdle", MType: "gauge", Value: PtrFloat64(float64(metrics.HeapIdle))})
	a.metricsCollector.Collect(collector.MetricJSON{ID: "HeapInuse", MType: "gauge", Value: PtrFloat64(float64(metrics.HeapInuse))})
	a.metricsCollector.Collect(collector.MetricJSON{ID: "HeapObjects", MType: "gauge", Value: PtrFloat64(float64(metrics.HeapObjects))})
	a.metricsCollector.Collect(collector.MetricJSON{ID: "HeapReleased", MType: "gauge", Value: PtrFloat64(float64(metrics.HeapReleased))})
	a.metricsCollector.Collect(collector.MetricJSON{ID: "HeapSys", MType: "gauge", Value: PtrFloat64(float64(metrics.HeapSys))})
	a.metricsCollector.Collect(collector.MetricJSON{ID: "Lookups", MType: "gauge", Value: PtrFloat64(float64(metrics.Lookups))})
	a.metricsCollector.Collect(collector.MetricJSON{ID: "MCacheInuse", MType: "gauge", Value: PtrFloat64(float64(metrics.MCacheInuse))})
	a.metricsCollector.Collect(collector.MetricJSON{ID: "MCacheSys", MType: "gauge", Value: PtrFloat64(float64(metrics.MCacheSys))})
	a.metricsCollector.Collect(collector.MetricJSON{ID: "MSpanInuse", MType: "gauge", Value: PtrFloat64(float64(metrics.MSpanInuse))})
	a.metricsCollector.Collect(collector.MetricJSON{ID: "MSpanSys", MType: "gauge", Value: PtrFloat64(float64(metrics.MSpanSys))})
	a.metricsCollector.Collect(collector.MetricJSON{ID: "Mallocs", MType: "gauge", Value: PtrFloat64(float64(metrics.Mallocs))})
	a.metricsCollector.Collect(collector.MetricJSON{ID: "NextGC", MType: "gauge", Value: PtrFloat64(float64(metrics.NextGC))})
	a.metricsCollector.Collect(collector.MetricJSON{ID: "NumForcedGC", MType: "gauge", Value: PtrFloat64(float64(metrics.NumForcedGC))})
	a.metricsCollector.Collect(collector.MetricJSON{ID: "NumGC", MType: "gauge", Value: PtrFloat64(float64(metrics.NumGC))})
	a.metricsCollector.Collect(collector.MetricJSON{ID: "OtherSys", MType: "gauge", Value: PtrFloat64(float64(metrics.OtherSys))})
	a.metricsCollector.Collect(collector.MetricJSON{ID: "PauseTotalNs", MType: "gauge", Value: PtrFloat64(float64(metrics.PauseTotalNs))})
	a.metricsCollector.Collect(collector.MetricJSON{ID: "StackInuse", MType: "gauge", Value: PtrFloat64(float64(metrics.StackInuse))})
	a.metricsCollector.Collect(collector.MetricJSON{ID: "StackSys", MType: "gauge", Value: PtrFloat64(float64(metrics.StackSys))})
	a.metricsCollector.Collect(collector.MetricJSON{ID: "Sys", MType: "gauge", Value: PtrFloat64(float64(metrics.Sys))})
	a.metricsCollector.Collect(collector.MetricJSON{ID: "TotalAlloc", MType: "gauge", Value: PtrFloat64(float64(metrics.TotalAlloc))})
	a.metricsCollector.Collect(collector.MetricJSON{ID: "RandomValue", MType: "gauge", Value: PtrFloat64(float64(rand.Int()))})
	a.metricsCollector.Collect(collector.MetricJSON{ID: "LastGC", MType: "gauge", Value: PtrFloat64(float64(metrics.LastGC))})

	cnt, _ := collector.Collector.GetMetric("PollCount")
	counter := int64(0)
	if cnt.Delta != nil {
		counter = *cnt.Delta + 1
	}
	collector.Collector.Collect(collector.MetricJSON{ID: "PollCount", MType: "counter", Delta: PtrInt64(counter)})
}

func New(metricsCollector collectorImpl) *Aggregator {
	return &Aggregator{
		metricsCollector: metricsCollector,
	}
}

type Aggregator struct {
	metricsCollector collectorImpl
}

type collectorImpl interface {
	Collect(json collector.MetricJSON) error
}

func PtrFloat64(f float64) *float64 {
	return &f
}

func PtrInt64(i int64) *int64 {
	return &i
}
