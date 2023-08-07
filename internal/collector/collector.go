package collector

import (
	"encoding/json"
	"errors"
	"strconv"
)

var (
	ErrBadRequest     = errors.New("bad request")
	ErrNotImplemented = errors.New("not implemented")
	ErrNotFound       = errors.New("not found")
)

var Collector = collector{
	Metrics: make([]StoredMetric, 0),
}

// GetMetric - a method for gettong metric value by metric name.
// Returns the struct of type StoredMetric
func (c *collector) GetMetric(metricName string) (StoredMetric, error) {
	for _, m := range c.Metrics {
		if m.ID == metricName {
			return m, nil
		}
	}
	return StoredMetric{}, ErrNotFound
}

// GetMetricJSON - a metric for getting metric value by metric name.
// Returns the json.
func (c *collector) GetMetricJSON(metricName string) ([]byte, error) {
	for _, m := range c.Metrics {
		if m.ID == metricName {
			resultJSON, err := json.Marshal(m)
			if err != nil {
				return nil, ErrBadRequest
			}
			return resultJSON, nil
		}
	}
	return nil, ErrNotFound
}

// GetAvailableMetrics - a method for getting all metrics names from
// the server. Returns the list of strings.
func (c *collector) GetAvailableMetrics() []string {
	names := make([]string, 0, len(c.Metrics))
	for _, m := range c.Metrics {
		names = append(names, m.ID)
	}
	return names
}

// UpsertMetric - a method for upserting new metric.
func (c *collector) UpsertMetric(metric StoredMetric) {
	for i, m := range c.Metrics {
		if m.ID == metric.ID {
			c.Metrics[i] = metric
			return
		}
	}
	c.Metrics = append(c.Metrics, metric)
}

// Collect - a method for upserting metric from the MetricRequest.
func (c *collector) Collect(metric MetricRequest, metricValue string) error {
	if (metric.Delta != nil && *metric.Delta < 0) || (metric.Value != nil && *metric.Value < 0) || metric.ID == "" {
		return ErrBadRequest
	}

	switch metric.MType {
	case Counter:
		v, err := c.GetMetric(metric.ID)
		if err != nil {
			if !errors.Is(err, ErrNotFound) {
				return err
			}
		}
		value, err := strconv.Atoi(metricValue)
		if err != nil {
			return ErrBadRequest
		}
		if v.CounterValue != nil {
			value = value + int(*v.CounterValue)
		}
		c.UpsertMetric(StoredMetric{
			ID:           metric.ID,
			MType:        metric.MType,
			CounterValue: PtrInt64(int64(value)),
			TextValue:    PtrString(strconv.Itoa(value)),
		})
	case Gauge:
		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			return ErrBadRequest
		}
		c.UpsertMetric(StoredMetric{
			ID:         metric.ID,
			MType:      metric.MType,
			GaugeValue: &value,
			TextValue:  &metricValue,
		})
	default:
		return ErrNotImplemented
	}
	return nil
}

func PtrFloat64(f float64) *float64 {
	return &f
}

func PtrInt64(i int64) *int64 {
	return &i
}

func PtrString(s string) *string {
	return &s
}
