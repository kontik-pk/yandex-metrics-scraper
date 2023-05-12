package collector

import (
	"encoding/json"
	"errors"
)

var (
	ErrBadRequest     = errors.New("bad request")
	ErrNotImplemented = errors.New("not implemented")
	ErrNotFound       = errors.New("not found")
)

var Collector = collector{
	Metrics: make([]MetricJSON, 0),
}

func (c *collector) Collect(metric MetricJSON) error {
	if (metric.Delta != nil && *metric.Delta < 0) || (metric.Value != nil && *metric.Value < 0) {
		return ErrBadRequest
	}
	switch metric.MType {
	case "counter":
		v, err := c.GetMetric(metric.ID)
		if err != nil {
			if !errors.Is(err, ErrNotFound) {
				return err
			}
		}
		if v.Delta != nil {
			*metric.Delta += *v.Delta
		}
		c.UpsertMetric(metric)

	case "gauge":
		c.UpsertMetric(metric)
	default:
		return ErrNotImplemented
	}
	return nil
}

func (c *collector) GetMetric(metricName string) (MetricJSON, error) {
	for _, m := range c.Metrics {
		if m.ID == metricName {
			return m, nil
		}
	}
	return MetricJSON{}, ErrNotFound
}

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

func (c *collector) GetAvailableMetrics() []string {
	names := make([]string, 0)
	for _, m := range c.Metrics {
		names = append(names, m.ID)
	}
	return names
}

func (c *collector) UpsertMetric(metric MetricJSON) {
	for i, m := range c.Metrics {
		if m.ID == metric.ID {
			c.Metrics[i] = metric
			return
		}
	}
	c.Metrics = append(c.Metrics, metric)
}
