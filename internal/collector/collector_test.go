package collector

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"runtime"
	"strconv"
	"testing"
)

func TestCollector_Collect(t *testing.T) {
	testCases := []struct {
		name     string
		storage  collector
		metric   runtime.MemStats
		expected memStorage
	}{
		{
			name:    "case0",
			storage: collector{storage: &memStorage{gauges: map[string]string{}, counters: map[string]int{}}},
			metric:  runtime.MemStats{Alloc: 1, Sys: 1, GCCPUFraction: 5.543},
			expected: memStorage{
				gauges: map[string]string{
					"Alloc":         "1",
					"Sys":           "1",
					"GCCPUFraction": "5.543",
				},
			},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.storage.Collect("Alloc", "gauge", strconv.FormatUint(tt.metric.Alloc, 10))
			assert.NoError(t, err)
			err = tt.storage.Collect("Sys", "gauge", strconv.FormatUint(tt.metric.Sys, 10))
			assert.NoError(t, err)
			err = tt.storage.Collect("GCCPUFraction", "gauge", fmt.Sprintf("%.3f", tt.metric.GCCPUFraction))
			assert.NoError(t, err)

			assert.Equal(t, tt.expected.gauges, tt.storage.GetGauges())
		})
	}
}
