package collector

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCollector_Collect(t *testing.T) {
	testCases := []struct {
		name          string
		storage       collector
		request       MetricRequest
		metricValue   string
		expected      []StoredMetric
		expectedError error
	}{
		{
			name:    "case0",
			storage: collector{[]StoredMetric{}},
			request: MetricRequest{
				ID:    "Alloc",
				MType: "gauge",
				Value: PtrFloat64(1),
			},
			metricValue: "1",
			expected: []StoredMetric{
				{
					ID:         "Alloc",
					MType:      "gauge",
					GaugeValue: PtrFloat64(1),
					TextValue:  PtrString("1"),
				},
			},
		},
		{
			name: "case1",
			storage: collector{[]StoredMetric{
				{
					ID:         "Alloc",
					MType:      "gauge",
					GaugeValue: PtrFloat64(10),
					TextValue:  PtrString("10"),
				},
				{
					ID:         "GCCPUFraction",
					MType:      "gauge",
					GaugeValue: PtrFloat64(5.543),
					TextValue:  PtrString("5.543"),
				},
			}},
			request: MetricRequest{
				ID:    "Alloc",
				MType: "gauge",
				Value: PtrFloat64(1),
			},
			metricValue: "1",
			expected: []StoredMetric{
				{
					ID:         "Alloc",
					MType:      "gauge",
					GaugeValue: PtrFloat64(1),
					TextValue:  PtrString("1"),
				},
				{
					ID:         "GCCPUFraction",
					MType:      "gauge",
					GaugeValue: PtrFloat64(5.543),
					TextValue:  PtrString("5.543"),
				},
			},
		},
		{
			name: "case2",
			storage: collector{Metrics: []StoredMetric{
				{
					ID:         "Alloc",
					MType:      "gauge",
					GaugeValue: PtrFloat64(3),
					TextValue:  PtrString("3"),
				},
				{
					ID:         "Sys",
					MType:      "gauge",
					GaugeValue: PtrFloat64(5),
					TextValue:  PtrString("5"),
				},
				{
					ID:           "Counter",
					MType:        "counter",
					CounterValue: PtrInt64(5),
					TextValue:    PtrString("5"),
				},
			}},
			request: MetricRequest{
				ID:    "Counter",
				MType: "counter",
				Delta: PtrInt64(10),
			},
			metricValue: "10",
			expected: []StoredMetric{
				{
					ID:         "Alloc",
					MType:      "gauge",
					GaugeValue: PtrFloat64(3),
					TextValue:  PtrString("3"),
				},
				{
					ID:         "Sys",
					MType:      "gauge",
					GaugeValue: PtrFloat64(5),
					TextValue:  PtrString("5"),
				},
				{
					ID:           "Counter",
					MType:        "counter",
					CounterValue: PtrInt64(15),
					TextValue:    PtrString("15"),
				},
			},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.storage.Collect(tt.request, tt.metricValue)
			if tt.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.expectedError.Error())
			}
			assert.Equal(t, tt.expected, tt.storage.Metrics)
		})
	}
}
