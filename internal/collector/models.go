package collector

const (
	Counter = "counter"
	Gauge   = "gauge"
)

type MetricRequest struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

type StoredMetric struct {
	ID           string   `json:"id"`                      // имя метрики
	MType        string   `json:"type"`                    // параметр, принимающий значение gauge или counter
	CounterValue *int64   `json:"counter_value,omitempty"` // значение метрики в случае передачи counter
	GaugeValue   *float64 `json:"gauge_value,omitempty"`   // значение метрики в случае передачи gauge
	TextValue    *string  `json:"text_value,omitempty"`
}

type collector struct {
	Metrics []StoredMetric
}
