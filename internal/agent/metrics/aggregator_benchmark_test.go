package aggregator

import (
	"github.com/kontik-pk/yandex-metrics-scraper/internal/agent/collector"
	"testing"
)

func BenchmarkAggregator_AggregateGopsutilMetrics(b *testing.B) {
	b.Run("aggregate gopsutil metrics benchmark", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			metricsCollector := collector.Collector()
			metricsAggregator := New(metricsCollector)
			metricsAggregator.AggregateGopsutilMetrics()
		}
	})
}

func BenchmarkAggregator_AggregateRuntimeMetrics(b *testing.B) {
	b.Run("aggregate runtime metrics benchmark", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			metricsCollector := collector.Collector()
			metricsAggregator := New(metricsCollector)
			metricsAggregator.AggregateRuntimeMetrics()
		}
	})
}
