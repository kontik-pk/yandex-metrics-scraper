package collector

import (
	"bufio"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestCollector_Collect(t *testing.T) {
	testCases := []struct {
		name          string
		storage       collector
		metricName    string
		metricType    string
		metricValue   string
		expected      memStorage
		expectedError error
	}{
		{
			name:        "case0",
			storage:     collector{storage: &memStorage{Gauges: map[string]string{}, Counters: map[string]int{}}},
			metricName:  "Alloc",
			metricType:  "gauge",
			metricValue: "1",
			expected: memStorage{
				Gauges: map[string]string{
					"Alloc": "1",
				},
			},
		},
		{
			name: "case1",
			storage: collector{storage: &memStorage{Gauges: map[string]string{
				"Alloc":         "3",
				"GCCPUFraction": "5.543",
			}, Counters: map[string]int{}}},
			metricName:  "Alloc",
			metricType:  "gauge",
			metricValue: "1",
			expected: memStorage{
				Gauges: map[string]string{
					"Alloc":         "1",
					"GCCPUFraction": "5.543",
				},
			},
		},
		{
			name: "case3",
			storage: collector{storage: &memStorage{Gauges: map[string]string{
				"Alloc": "3",
				"Sys":   "5",
			}, Counters: map[string]int{
				"Counter": 5,
			}}},
			metricName:  "Counter",
			metricType:  "counter",
			metricValue: "10",
			expected: memStorage{
				Gauges: map[string]string{
					"Alloc": "3",
					"Sys":   "5",
				},
				Counters: map[string]int{
					"Counter": 15,
				},
			},
		},
		{
			name: "case4",
			storage: collector{storage: &memStorage{Gauges: map[string]string{
				"Alloc": "3",
				"Sys":   "5",
			}, Counters: map[string]int{}}},
			metricName:  "Counter",
			metricType:  "counter",
			metricValue: "10",
			expected: memStorage{
				Gauges: map[string]string{
					"Alloc": "3",
					"Sys":   "5",
				},
				Counters: map[string]int{
					"Counter": 10,
				},
			},
		},
		{
			name:        "case5",
			storage:     collector{storage: &memStorage{Gauges: map[string]string{}, Counters: map[string]int{}}},
			metricName:  "Alloc",
			metricType:  "gauge",
			metricValue: "1.0000000",
			expected: memStorage{
				Gauges: map[string]string{
					"Alloc": "1.0000000",
				},
			},
		},
		{
			name:        "case5",
			storage:     collector{storage: &memStorage{Gauges: map[string]string{}, Counters: map[string]int{}}},
			metricName:  "Alloc",
			metricType:  "gauge",
			metricValue: "invalid",
			expected: memStorage{
				map[string]int{},
				map[string]string{},
			},
			expectedError: ErrBadRequest,
		},
		{
			name:        "case5",
			storage:     collector{storage: &memStorage{Gauges: map[string]string{}, Counters: map[string]int{}}},
			metricName:  "Alloc",
			metricType:  "invalid",
			metricValue: "15",
			expected: memStorage{
				map[string]int{},
				map[string]string{},
			},
			expectedError: ErrNotImplemented,
		},
		{
			name:        "case5",
			storage:     collector{storage: &memStorage{Gauges: map[string]string{}, Counters: map[string]int{}}},
			metricName:  "Alloc",
			metricType:  "counter",
			metricValue: "17.001",
			expected: memStorage{
				map[string]int{},
				map[string]string{},
			},
			expectedError: ErrBadRequest,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.storage.Collect(tt.metricName, tt.metricType, tt.metricValue)
			if tt.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.expectedError.Error())
			}
			assert.Equal(t, tt.expected.Gauges, tt.storage.GetGauges())
		})
	}
}

func TestCollector_GetMetric(t *testing.T) {
	_ = Collector.Collect("Counter1", "counter", "0")
	_ = Collector.Collect("Counter2", "counter", "15")
	_ = Collector.Collect("Gauge1", "gauge", "17.01")
	_ = Collector.Collect("Gauge2", "gauge", "18.00000")

	testCases := []struct {
		name          string
		metricName    string
		metricType    string
		expectedValue string
		expectedError error
	}{
		{
			name:          "case0",
			metricType:    "counter",
			metricName:    "Counter1",
			expectedValue: "0",
		},
		{
			name:          "case1",
			metricType:    "counter",
			metricName:    "Counter2",
			expectedValue: "15",
		},
		{
			name:          "case2",
			metricType:    "gauge",
			metricName:    "Gauge1",
			expectedValue: "17.01",
		},
		{
			name:          "case3",
			metricType:    "gauge",
			metricName:    "Gauge2",
			expectedValue: "18.00000",
		},
		{
			name:          "case4",
			metricType:    "gauge",
			metricName:    "Gauge3",
			expectedValue: "",
			expectedError: ErrNotFound,
		},
		{
			name:          "case5",
			metricType:    "invalid",
			metricName:    "Gauge2",
			expectedValue: "",
			expectedError: ErrNotImplemented,
		},
		{
			name:          "case6",
			metricType:    "counter",
			metricName:    "Counter3",
			expectedValue: "",
			expectedError: ErrNotFound,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			value, err := Collector.GetMetricByName(tt.metricName, tt.metricType)
			if tt.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.expectedError.Error())
			}
			assert.Equal(t, value, tt.expectedValue)
		})
	}
}

func TestCollector_CollectFromJSON(t *testing.T) {
	testCases := []struct {
		name          string
		metric        MetricJSON
		expectedError string
	}{
		{
			name: "positive (collect counter)",
			metric: MetricJSON{
				ID:    "metricValidName",
				MType: "counter",
				Delta: ptrInt(5),
			},
		},
		{
			name: "positive (collect gauge)",
			metric: MetricJSON{
				ID:    "metricValidName",
				MType: "gauge",
				Value: ptrFloat(5.727),
			},
		},
		{
			name: "negative (invalid metric type)",
			metric: MetricJSON{
				ID:    "metricValidName",
				MType: "invalid metric type",
				Value: ptrFloat(5.727),
			},
			expectedError: "not implemented",
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			testCollector := collector{
				storage: &memStorage{
					Counters: make(map[string]int),
					Gauges:   make(map[string]string),
				},
			}
			err := testCollector.CollectFromJSON(tt.metric)
			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCollector_GetAvailableMetrics(t *testing.T) {
	testCases := []struct {
		name             string
		metricsCollector collector
		expectedMetrics  []string
	}{
		{
			name: "case0",
			metricsCollector: collector{
				storage: &memStorage{
					Counters: map[string]int{
						"firstCounter":  1,
						"secondCounter": 2,
					},
					Gauges: map[string]string{
						"firstGauge":  "1.35",
						"secondGauge": "2.67",
					},
				},
			},
			expectedMetrics: []string{"firstCounter", "secondCounter", "firstGauge", "secondGauge"},
		},
		{
			name: "case1",
			metricsCollector: collector{
				storage: &memStorage{
					Counters: map[string]int{},
					Gauges: map[string]string{
						"firstGauge":  "1.35",
						"secondGauge": "2.67",
					},
				},
			},
			expectedMetrics: []string{"firstGauge", "secondGauge"},
		},
		{
			name: "case2",
			metricsCollector: collector{
				storage: &memStorage{
					Counters: map[string]int{
						"firstCounter":  1,
						"secondCounter": 2,
					},
					Gauges: map[string]string{},
				},
			},
			expectedMetrics: []string{"firstCounter", "secondCounter"},
		},
		{
			name: "case3",
			metricsCollector: collector{
				storage: &memStorage{
					Counters: map[string]int{},
					Gauges:   map[string]string{},
				},
			},
			expectedMetrics: []string{},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			actualMetrics := tt.metricsCollector.GetAvailableMetrics()
			assert.ElementsMatch(t, actualMetrics, tt.expectedMetrics)
		})
	}
}

func TestCollector_GetCounters(t *testing.T) {
	testCases := []struct {
		name             string
		metricsCollector collector
		expectedCounters map[string]string
	}{
		{
			name: "case0",
			metricsCollector: collector{
				storage: &memStorage{
					Counters: map[string]int{
						"firstCounter":  1,
						"secondCounter": 2,
					},
					Gauges: map[string]string{
						"firstGauge":  "1.35",
						"secondGauge": "2.67",
					},
				},
			},
			expectedCounters: map[string]string{
				"firstCounter":  "1",
				"secondCounter": "2",
			},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			actualMetrics := tt.metricsCollector.GetCounters()
			assert.Equal(t, actualMetrics, tt.expectedCounters)
		})
	}
}

func TestCollector_Save(t *testing.T) {
	testCases := []struct {
		name             string
		metricsCollector collector
		expectedError    string
	}{
		{
			name: "case0",
			metricsCollector: collector{
				storage: &memStorage{
					Counters: map[string]int{
						"firstCounter":  1,
						"secondCounter": 2,
					},
					Gauges: map[string]string{
						"firstGauge":  "1.35",
						"secondGauge": "2.67",
					},
				},
			},
		},
		{
			name: "case1",
			metricsCollector: collector{
				storage: &memStorage{
					Counters: map[string]int{
						"firstCounter": 1,
					},
					Gauges: map[string]string{},
				},
			},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			filePath := "/tmp/test_save.json"
			err := tt.metricsCollector.Save(filePath)
			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

			file, err := os.OpenFile(filePath, os.O_RDONLY, 0666)
			assert.NoError(t, err)
			scanner := bufio.NewScanner(file)
			assert.True(t, scanner.Scan())
			data := scanner.Bytes()
			metricsFromFile := memStorage{}
			err = json.Unmarshal(data, &metricsFromFile)
			assert.NoError(t, err)
			assert.Equal(t, metricsFromFile.Gauges, tt.metricsCollector.storage.Gauges)
			assert.Equal(t, metricsFromFile.Counters, tt.metricsCollector.storage.Counters)
		})
	}
}

func TestCollector_Restore(t *testing.T) {
	t.Run("case0", func(t *testing.T) {
		filePath := "/tmp/test_save.json"

		previousCollector := collector{
			storage: &memStorage{
				Counters: map[string]int{
					"firstCounter":  1,
					"secondCounter": 2,
				},
				Gauges: map[string]string{
					"firstGauge":  "1.35",
					"secondGauge": "2.67",
				},
			},
		}
		err := previousCollector.Save(filePath)
		assert.NoError(t, err)

		newCollector := collector{
			storage: &memStorage{
				Counters: map[string]int{},
				Gauges:   map[string]string{},
			},
		}
		err = newCollector.Restore(filePath)
		assert.NoError(t, err)

		assert.Equal(t, newCollector, previousCollector)
	})
}

func TestCollector_GetMetricJSON(t *testing.T) {
	metricsCollector := collector{
		storage: &memStorage{
			Counters: map[string]int{
				"firstCounter":  1,
				"secondCounter": 2,
			},
			Gauges: map[string]string{
				"firstGauge":  "1.35",
				"secondGauge": "2.67",
			},
		},
	}

	testCases := []struct {
		name           string
		metricName     string
		metricType     string
		expectedResult *MetricJSON
		expectedError  string
	}{
		{
			name:       "positive (counter)",
			metricName: "firstCounter",
			metricType: "counter",
			expectedResult: &MetricJSON{
				ID:    "firstCounter",
				MType: "counter",
				Delta: ptrInt(1),
			},
		},
		{
			name:       "positive (gauge)",
			metricName: "firstGauge",
			metricType: "gauge",
			expectedResult: &MetricJSON{
				ID:    "firstGauge",
				MType: "gauge",
				Value: ptrFloat(1.35),
			},
		},
		{
			name:           "negative (invalid type)",
			metricName:     "firstCounter",
			metricType:     "invalid",
			expectedResult: nil,
			expectedError:  "not implemented",
		},
		{
			name:           "negative (not found)",
			metricName:     "invalid",
			metricType:     "counter",
			expectedResult: nil,
			expectedError:  "not found",
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := metricsCollector.GetMetricJSON(tt.metricName, tt.metricType)
			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
			if tt.expectedResult != nil {
				expected, err := json.Marshal(tt.expectedResult)
				assert.NoError(t, err)
				assert.Equal(t, actual, expected)
			} else {
				assert.Equal(t, actual, []byte(nil))
			}
		})
	}
}

func ptrInt(variable int64) *int64 {
	return &variable
}

func ptrFloat(variable float64) *float64 {
	return &variable
}
