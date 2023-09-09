package aggregator

import (
	collector2 "github.com/kontik-pk/yandex-metrics-scraper/internal/agent/collector"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAggregator_AggregateGopsutilMetrics(t *testing.T) {
	t.Run("aggregate gopsutil metrics test", func(t *testing.T) {
		metricsCollector := collector2.Collector()
		metricsCollector.Metrics = []collector2.StoredMetric{}
		metricsAggregator := New(metricsCollector)
		metricsAggregator.AggregateGopsutilMetrics()
		assert.Equal(t, metricsCollector.GetAvailableMetrics(), []string{"FreeMemory", "TotalMemory", "CPUutilization1"})
	})
}

func TestAggregator_AggregateRuntimeMetrics(t *testing.T) {
	t.Run("aggregate runtime metrics test", func(t *testing.T) {
		metricsCollector := collector2.Collector()
		metricsCollector.Metrics = []collector2.StoredMetric{}
		metricsAggregator := New(metricsCollector)
		metricsAggregator.AggregateRuntimeMetrics()
		assert.Equal(t, metricsCollector.GetAvailableMetrics(), []string{"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "GCSys", "HeapAlloc", "HeapIdle", "HeapInuse", "HeapObjects", "HeapReleased", "HeapSys", "Lookups", "MCacheInuse", "MCacheSys", "MSpanInuse", "MSpanSys", "Mallocs", "NextGC", "NumForcedGC", "NumGC", "OtherSys", "PauseTotalNs", "StackInuse", "StackSys", "Sys", "TotalAlloc", "RandomValue", "LastGC", "PollCount"})
	})
}
