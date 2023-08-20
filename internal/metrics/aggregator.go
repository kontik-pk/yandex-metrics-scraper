package aggregator

import (
	"github.com/kontik-pk/yandex-metrics-scraper/internal/collector"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"math/rand"
	"runtime"
	"strconv"
)

// AggregateRuntimeMetrics - a method for capturing and upserting runtime metrics.
func (a *Aggregator) AggregateRuntimeMetrics() {
	metrics := runtime.MemStats{}
	runtime.ReadMemStats(&metrics)

	//TODO: мерзкий парсинг структуры, можно ли тут улучшить?
	a.metricsCollector.UpsertMetric(collector.StoredMetric{ID: "Alloc", MType: "gauge", GaugeValue: collector.PtrFloat64(float64(metrics.Alloc)), TextValue: collector.PtrString(strconv.FormatFloat(float64(metrics.Alloc), 'f', 11, 64))})
	a.metricsCollector.UpsertMetric(collector.StoredMetric{ID: "BuckHashSys", MType: "gauge", GaugeValue: collector.PtrFloat64(float64(metrics.BuckHashSys)), TextValue: collector.PtrString(strconv.FormatFloat(float64(metrics.BuckHashSys), 'f', 11, 64))})
	a.metricsCollector.UpsertMetric(collector.StoredMetric{ID: "Frees", MType: "gauge", GaugeValue: collector.PtrFloat64(float64(metrics.Frees)), TextValue: collector.PtrString(strconv.FormatFloat(float64(metrics.Frees), 'f', 11, 64))})
	a.metricsCollector.UpsertMetric(collector.StoredMetric{ID: "GCCPUFraction", MType: "gauge", GaugeValue: &metrics.GCCPUFraction, TextValue: collector.PtrString(strconv.FormatFloat(metrics.GCCPUFraction, 'f', 11, 64))})
	a.metricsCollector.UpsertMetric(collector.StoredMetric{ID: "GCSys", MType: "gauge", GaugeValue: collector.PtrFloat64(float64(metrics.GCSys)), TextValue: collector.PtrString(strconv.FormatFloat(float64(metrics.GCSys), 'f', 11, 64))})
	a.metricsCollector.UpsertMetric(collector.StoredMetric{ID: "HeapAlloc", MType: "gauge", GaugeValue: collector.PtrFloat64(float64(metrics.HeapAlloc)), TextValue: collector.PtrString(strconv.FormatFloat(float64(metrics.HeapAlloc), 'f', 11, 64))})
	a.metricsCollector.UpsertMetric(collector.StoredMetric{ID: "HeapIdle", MType: "gauge", GaugeValue: collector.PtrFloat64(float64(metrics.HeapIdle)), TextValue: collector.PtrString(strconv.FormatFloat(float64(metrics.HeapIdle), 'f', 11, 64))})
	a.metricsCollector.UpsertMetric(collector.StoredMetric{ID: "HeapInuse", MType: "gauge", GaugeValue: collector.PtrFloat64(float64(metrics.HeapInuse)), TextValue: collector.PtrString(strconv.FormatFloat(float64(metrics.HeapInuse), 'f', 11, 64))})
	a.metricsCollector.UpsertMetric(collector.StoredMetric{ID: "HeapObjects", MType: "gauge", GaugeValue: collector.PtrFloat64(float64(metrics.HeapObjects)), TextValue: collector.PtrString(strconv.FormatFloat(float64(metrics.HeapObjects), 'f', 11, 64))})
	a.metricsCollector.UpsertMetric(collector.StoredMetric{ID: "HeapReleased", MType: "gauge", GaugeValue: collector.PtrFloat64(float64(metrics.HeapReleased)), TextValue: collector.PtrString(strconv.FormatFloat(float64(metrics.HeapReleased), 'f', 11, 64))})
	a.metricsCollector.UpsertMetric(collector.StoredMetric{ID: "HeapSys", MType: "gauge", GaugeValue: collector.PtrFloat64(float64(metrics.HeapSys)), TextValue: collector.PtrString(strconv.FormatFloat(float64(metrics.HeapSys), 'f', 11, 64))})
	a.metricsCollector.UpsertMetric(collector.StoredMetric{ID: "Lookups", MType: "gauge", GaugeValue: collector.PtrFloat64(float64(metrics.Lookups)), TextValue: collector.PtrString(strconv.FormatFloat(float64(metrics.Lookups), 'f', 11, 64))})
	a.metricsCollector.UpsertMetric(collector.StoredMetric{ID: "MCacheInuse", MType: "gauge", GaugeValue: collector.PtrFloat64(float64(metrics.MCacheInuse)), TextValue: collector.PtrString(strconv.FormatFloat(float64(metrics.MCacheInuse), 'f', 11, 64))})
	a.metricsCollector.UpsertMetric(collector.StoredMetric{ID: "MCacheSys", MType: "gauge", GaugeValue: collector.PtrFloat64(float64(metrics.MCacheSys)), TextValue: collector.PtrString(strconv.FormatFloat(float64(metrics.MCacheSys), 'f', 11, 64))})
	a.metricsCollector.UpsertMetric(collector.StoredMetric{ID: "MSpanInuse", MType: "gauge", GaugeValue: collector.PtrFloat64(float64(metrics.MSpanInuse)), TextValue: collector.PtrString(strconv.FormatFloat(float64(metrics.MSpanInuse), 'f', 11, 64))})
	a.metricsCollector.UpsertMetric(collector.StoredMetric{ID: "MSpanSys", MType: "gauge", GaugeValue: collector.PtrFloat64(float64(metrics.MSpanSys)), TextValue: collector.PtrString(strconv.FormatFloat(float64(metrics.MSpanSys), 'f', 11, 64))})
	a.metricsCollector.UpsertMetric(collector.StoredMetric{ID: "Mallocs", MType: "gauge", GaugeValue: collector.PtrFloat64(float64(metrics.Mallocs)), TextValue: collector.PtrString(strconv.FormatFloat(float64(metrics.Mallocs), 'f', 11, 64))})
	a.metricsCollector.UpsertMetric(collector.StoredMetric{ID: "NextGC", MType: "gauge", GaugeValue: collector.PtrFloat64(float64(metrics.NextGC)), TextValue: collector.PtrString(strconv.FormatFloat(float64(metrics.NextGC), 'f', 11, 64))})
	a.metricsCollector.UpsertMetric(collector.StoredMetric{ID: "NumForcedGC", MType: "gauge", GaugeValue: collector.PtrFloat64(float64(metrics.NumForcedGC)), TextValue: collector.PtrString(strconv.FormatFloat(float64(metrics.NumForcedGC), 'f', 11, 64))})
	a.metricsCollector.UpsertMetric(collector.StoredMetric{ID: "NumGC", MType: "gauge", GaugeValue: collector.PtrFloat64(float64(metrics.NumGC)), TextValue: collector.PtrString(strconv.FormatFloat(float64(metrics.NumGC), 'f', 11, 64))})
	a.metricsCollector.UpsertMetric(collector.StoredMetric{ID: "OtherSys", MType: "gauge", GaugeValue: collector.PtrFloat64(float64(metrics.OtherSys)), TextValue: collector.PtrString(strconv.FormatFloat(float64(metrics.OtherSys), 'f', 11, 64))})
	a.metricsCollector.UpsertMetric(collector.StoredMetric{ID: "PauseTotalNs", MType: "gauge", GaugeValue: collector.PtrFloat64(float64(metrics.PauseTotalNs)), TextValue: collector.PtrString(strconv.FormatFloat(float64(metrics.PauseTotalNs), 'f', 11, 64))})
	a.metricsCollector.UpsertMetric(collector.StoredMetric{ID: "StackInuse", MType: "gauge", GaugeValue: collector.PtrFloat64(float64(metrics.StackInuse)), TextValue: collector.PtrString(strconv.FormatFloat(float64(metrics.StackInuse), 'f', 11, 64))})
	a.metricsCollector.UpsertMetric(collector.StoredMetric{ID: "StackSys", MType: "gauge", GaugeValue: collector.PtrFloat64(float64(metrics.StackSys)), TextValue: collector.PtrString(strconv.FormatFloat(float64(metrics.StackSys), 'f', 11, 64))})
	a.metricsCollector.UpsertMetric(collector.StoredMetric{ID: "Sys", MType: "gauge", GaugeValue: collector.PtrFloat64(float64(metrics.Sys)), TextValue: collector.PtrString(strconv.FormatFloat(float64(metrics.Sys), 'f', 11, 64))})
	a.metricsCollector.UpsertMetric(collector.StoredMetric{ID: "TotalAlloc", MType: "gauge", GaugeValue: collector.PtrFloat64(float64(metrics.TotalAlloc)), TextValue: collector.PtrString(strconv.FormatFloat(float64(metrics.TotalAlloc), 'f', 11, 64))})
	a.metricsCollector.UpsertMetric(collector.StoredMetric{ID: "RandomValue", MType: "gauge", GaugeValue: collector.PtrFloat64(float64(rand.Int())), TextValue: collector.PtrString(strconv.FormatFloat(float64(rand.Int()), 'f', 11, 64))})
	a.metricsCollector.UpsertMetric(collector.StoredMetric{ID: "LastGC", MType: "gauge", GaugeValue: collector.PtrFloat64(float64(metrics.LastGC)), TextValue: collector.PtrString(strconv.FormatFloat(float64(metrics.LastGC), 'f', 11, 64))})

	cnt, _ := a.metricsCollector.GetMetric("PollCount")
	counter := int64(0)
	if cnt.CounterValue != nil {
		counter = *cnt.CounterValue + 1
	}
	a.metricsCollector.UpsertMetric(collector.StoredMetric{ID: "PollCount", MType: "counter", CounterValue: collector.PtrInt64(counter), TextValue: collector.PtrString(strconv.Itoa(int(counter)))})
}

// AggregateGopsutilMetrics - a method for capturing and upserting gopsutil metrics.
func (a *Aggregator) AggregateGopsutilMetrics() {
	v, _ := mem.VirtualMemory()
	cp, _ := cpu.Percent(0, false)
	a.metricsCollector.UpsertMetric(collector.StoredMetric{ID: "FreeMemory", MType: "gauge", GaugeValue: collector.PtrFloat64(float64(v.Free)), TextValue: collector.PtrString(strconv.FormatFloat(float64(v.Free), 'f', 11, 64))})
	a.metricsCollector.UpsertMetric(collector.StoredMetric{ID: "TotalMemory", MType: "gauge", GaugeValue: collector.PtrFloat64(float64(v.Total)), TextValue: collector.PtrString(strconv.FormatFloat(float64(v.Total), 'f', 11, 64))})
	a.metricsCollector.UpsertMetric(collector.StoredMetric{ID: "CPUutilization1", MType: "gauge", GaugeValue: collector.PtrFloat64(cp[0]), TextValue: collector.PtrString(strconv.FormatFloat(cp[0], 'f', 11, 64))})
}

// New is a function for creating `aggregator` object
func New(metricsCollector metricsCollector) *Aggregator {
	return &Aggregator{
		metricsCollector: metricsCollector,
	}
}

// Aggregator get metrics from runtime and gopsutil and upsert it
// to the metricsCollector.
type Aggregator struct {
	metricsCollector metricsCollector
}

type metricsCollector interface {
	UpsertMetric(metric collector.StoredMetric)
	GetMetric(metricName string) (collector.StoredMetric, error)
}
