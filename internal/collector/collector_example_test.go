package collector

import "fmt"

func Example_collect() {
	metricsCollector := collector{
		Metrics: make([]StoredMetric, 0),
	}
	if err := metricsCollector.Collect(MetricRequest{
		ID:    "metricName",
		MType: Gauge,
		Delta: PtrInt64(100500),
	}, "100500"); err != nil {
		panic(err)
	}
	fmt.Println(metricsCollector.Metrics[0].ID)
	fmt.Println(metricsCollector.Metrics[0].MType)
	fmt.Println(*metricsCollector.Metrics[0].TextValue)
	fmt.Println(*metricsCollector.Metrics[0].GaugeValue)
	// Output:
	// metricName
	// gauge
	// 100500
	// 100500
}

func Example_getMetricsJson() {
	metricsCollector := collector{
		Metrics: []StoredMetric{
			{
				ID:           "metricName",
				MType:        Counter,
				TextValue:    PtrString("100500"),
				CounterValue: PtrInt64(100500),
			},
		},
	}
	metric, err := metricsCollector.GetMetricJSON("metricName")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(metric))
	// Output:
	// {"id":"metricName","type":"counter","counter_value":100500,"text_value":"100500"}
}

func Example_getAvailableMetrics() {
	metricsCollector := collector{
		Metrics: []StoredMetric{
			{
				ID:           "metricCounterName",
				MType:        Counter,
				TextValue:    PtrString("100500"),
				CounterValue: PtrInt64(100500),
			},
			{
				ID:         "metricGaugeName",
				MType:      Gauge,
				TextValue:  PtrString("64.502"),
				GaugeValue: PtrFloat64(64.502),
			},
		},
	}
	fmt.Println(metricsCollector.GetAvailableMetrics())
	// Output:
	// [metricCounterName metricGaugeName]
}

func Example_getMetric() {
	metricsCollector := collector{
		Metrics: []StoredMetric{
			{
				ID:           "metricCounterName",
				MType:        Counter,
				TextValue:    PtrString("100500"),
				CounterValue: PtrInt64(100500),
			},
			{
				ID:         "metricGaugeName",
				MType:      Gauge,
				TextValue:  PtrString("64.502"),
				GaugeValue: PtrFloat64(64.502),
			},
		},
	}
	metric, err := metricsCollector.GetMetric("metricGaugeName")
	if err != nil {
		panic(err)
	}
	fmt.Println(metric.ID)
	fmt.Println(metric.MType)
	fmt.Println(*metric.GaugeValue)
	fmt.Println(*metric.TextValue)
	// Output:
	// metricGaugeName
	// gauge
	// 64.502
	// 64.502
}

func Example_upsertMetric() {
	metricsCollector := collector{
		Metrics: []StoredMetric{},
	}
	metricToSave := StoredMetric{
		ID:         "metricName",
		MType:      Gauge,
		GaugeValue: PtrFloat64(10.882900000),
		TextValue:  PtrString("10.882900000"),
	}
	metricsCollector.UpsertMetric(metricToSave)
	fmt.Println(metricsCollector.Metrics[0].ID)
	fmt.Println(*metricsCollector.Metrics[0].GaugeValue)
	fmt.Println(*metricsCollector.Metrics[0].TextValue)
	// Output:
	// metricName
	// 10.8829
	// 10.882900000
}
