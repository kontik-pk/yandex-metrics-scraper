package collector

import (
	"github.com/kontik-pk/yandex-metrics-scraper/internal/domain"
	"github.com/stretchr/testify/assert"
	"runtime"
	"testing"
)

func TestCollector_Collect(t *testing.T) {
	testCases := []struct {
		name     string
		storage  domain.MemStorage
		metric   runtime.MemStats
		expected domain.MemStorage
	}{
		{
			name:    "case0",
			storage: domain.MemStorage{Metrics: map[string]domain.Metric{}},
			metric:  runtime.MemStats{Alloc: 1, Sys: 1, GCCPUFraction: 5.543},
			expected: domain.MemStorage{Metrics: map[string]domain.Metric{
				"Alloc":         {MType: "gauge", Value: uint64(1)},
				"Sys":           {MType: "gauge", Value: uint64(1)},
				"GCCPUFraction": {MType: "gauge", Value: 5.543},
			}},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			metric := runtime.MemStats{Alloc: 1, Sys: 1, GCCPUFraction: 5.543}
			metricsCollector := New(&tt.storage)
			metricsCollector.Collect(&metric)
			assert.Equal(t, tt.expected.Metrics["Alloc"], tt.storage.Metrics["Alloc"])
			assert.Equal(t, tt.expected.Metrics["Sys"], tt.storage.Metrics["Sys"])
			assert.Equal(t, tt.expected.Metrics["GCCPUFraction"], tt.storage.Metrics["GCCPUFraction"])
		})
	}
}
