package collector

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
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
		{
			name: "case3",
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
				ID:    "",
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
					CounterValue: PtrInt64(5),
					TextValue:    PtrString("5"),
				},
			},
			expectedError: ErrBadRequest,
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

func TestCollector_GetAvailableMetrics(t *testing.T) {
	testCases := []struct {
		name            string
		collector       collector
		expectedMetrics []string
	}{
		{
			name: "positive",
			collector: collector{
				Metrics: []StoredMetric{
					{
						ID: "metric1",
					},
					{
						ID: "metric2",
					},
					{
						ID: "metric3",
					},
				},
			},
			expectedMetrics: []string{
				"metric1",
				"metric2",
				"metric3",
			},
		},
		{
			name: "positive: no metrics",
			collector: collector{
				Metrics: []StoredMetric{},
			},
			expectedMetrics: []string{},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			metrics := tt.collector.GetAvailableMetrics()
			assert.Equal(t, metrics, tt.expectedMetrics)
		})
	}
}

func TestCollector_GetMetricJSON(t *testing.T) {
	testCases := []struct {
		name           string
		metricName     string
		collector      collector
		expectedMetric *StoredMetric
		expectedError  error
	}{
		{
			name:       "positive",
			metricName: "metric1",
			collector: collector{
				Metrics: []StoredMetric{
					{
						ID:         "metric1",
						MType:      "gauge",
						GaugeValue: PtrFloat64(64.2),
						TextValue:  PtrString("64"),
					},
					{
						ID:         "metric2",
						MType:      "gauge",
						GaugeValue: PtrFloat64(128.2),
						TextValue:  PtrString("128"),
					},
					{
						ID:           "metric3",
						MType:        "counter",
						CounterValue: PtrInt64(64),
						TextValue:    PtrString("64"),
					},
				},
			},
			expectedMetric: &StoredMetric{
				ID:         "metric1",
				MType:      "gauge",
				GaugeValue: PtrFloat64(64.2),
				TextValue:  PtrString("64"),
			},
		},
		{
			name:       "negative: not found",
			metricName: "metric4",
			collector: collector{
				Metrics: []StoredMetric{
					{
						ID:         "metric1",
						MType:      "gauge",
						GaugeValue: PtrFloat64(64.2),
						TextValue:  PtrString("64"),
					},
					{
						ID:         "metric2",
						MType:      "gauge",
						GaugeValue: PtrFloat64(128.2),
						TextValue:  PtrString("128"),
					},
					{
						ID:           "metric3",
						MType:        "counter",
						CounterValue: PtrInt64(64),
						TextValue:    PtrString("64"),
					},
				},
			},
			expectedMetric: nil,
			expectedError:  ErrNotFound,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			metric, err := tt.collector.GetMetricJSON(tt.metricName)
			if tt.expectedError == nil {
				expected, _ := json.Marshal(tt.expectedMetric)
				assert.NoError(t, err)
				assert.Equal(t, expected, metric)
			} else {
				assert.EqualError(t, err, tt.expectedError.Error())
			}
		})
	}
}
