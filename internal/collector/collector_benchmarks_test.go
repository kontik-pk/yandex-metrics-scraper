package collector

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var testBenchCollector = collector{
	[]StoredMetric{
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
		{
			ID:           "IO",
			MType:        "counter",
			CounterValue: PtrInt64(5),
			TextValue:    PtrString("5"),
		},
		{
			ID:         "Mem",
			MType:      "gauge",
			GaugeValue: PtrFloat64(500.1992),
			TextValue:  PtrString("500.1992"),
		},
		{
			ID:           "Requests",
			MType:        "counter",
			CounterValue: PtrInt64(100500),
			TextValue:    PtrString("100500"),
		},
	},
}

func BenchmarkCollector_Collect(b *testing.B) {
	b.Run("collect benchmark", func(b *testing.B) {
		metric := MetricRequest{
			ID:    "new",
			MType: "gauge",
			Value: PtrFloat64(50.1001),
		}
		for i := 0; i < b.N; i++ {
			err := testBenchCollector.Collect(metric, "50.1001")
			assert.NoError(b, err)
		}
	})
}

func BenchmarkCollector_GetAvailableMetrics(b *testing.B) {
	b.Run("get available metrics benchmark", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			testBenchCollector.GetAvailableMetrics()
		}
	})
}

func BenchmarkCollector_GetMetric(b *testing.B) {
	b.Run("get metric benchmark", func(b *testing.B) {
		metricName := "Requests"
		for i := 0; i < b.N; i++ {
			_, err := testBenchCollector.GetMetric(metricName)
			assert.NoError(b, err)
		}
	})
}

func BenchmarkCollector_GetMetricJSON(b *testing.B) {
	b.Run("get metric json benchmark", func(b *testing.B) {
		metricName := "Requests"
		for i := 0; i < b.N; i++ {
			_, err := testBenchCollector.GetMetricJSON(metricName)
			assert.NoError(b, err)
		}
	})
}

func BenchmarkCollector_UpsertMetric(b *testing.B) {
	b.Run("upsert metric benchmark", func(b *testing.B) {
		metric := StoredMetric{
			ID:         "Alloc",
			MType:      "gauge",
			GaugeValue: PtrFloat64(3),
			TextValue:  PtrString("3"),
		}
		for i := 0; i < b.N; i++ {
			testBenchCollector.UpsertMetric(metric)
		}
	})
}
