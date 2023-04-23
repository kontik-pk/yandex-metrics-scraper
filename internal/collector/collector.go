package collector

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

var (
	ErrBadRequest     = errors.New("bad request")
	ErrNotImplemented = errors.New("not implemented")
	ErrNotFound       = errors.New("not found")
)

var Collector = collector{
	storage: &memStorage{
		counters: make(map[string]int),
		gauges:   make(map[string]string),
	},
}

func (c *collector) Collect(metricName string, metricType string, metricValue string) error {
	switch metricType {
	case "counter":
		value, err := strconv.Atoi(metricValue)
		if err != nil {
			return ErrBadRequest
		}
		c.storage.counters[metricName] += value
	case "gauge":
		_, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			return ErrBadRequest
		}
		c.storage.gauges[metricName] = metricValue
	default:
		return ErrNotImplemented
	}
	return nil
}

func (c *collector) CollectFromJSON(metric MetricJSON) error {
	metricValue := ""
	switch metric.MType {
	case "counter":
		metricValue = strconv.Itoa(int(*metric.Delta))
	case "gauge":
		metricValue = fmt.Sprintf("%.11f", *metric.Value)
	}

	return c.Collect(metric.ID, metric.MType, metricValue)
}

func (c *collector) GetMetricJSON(metricName string, metricType string) ([]byte, error) {
	updated, err := c.GetMetricByName(metricName, metricType)
	if err != nil {
		return nil, err
	}

	result := MetricJSON{
		ID:    metricName,
		MType: metricType,
	}
	switch result.MType {
	case "counter":
		counter, err := strconv.Atoi(updated)
		if err != nil {
			return nil, ErrBadRequest
		}
		c64 := int64(counter)
		result.Delta = &c64
	case "gauge":
		g, err := strconv.ParseFloat(updated, 64)
		if err != nil {
			return nil, ErrBadRequest
		}
		result.Value = &g
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		return nil, ErrBadRequest
	}
	return resultJSON, nil
}

func (c *collector) GetMetricByName(metricName string, metricType string) (string, error) {
	switch metricType {
	case "counter":
		value, ok := Collector.storage.counters[metricName]
		if !ok {
			return "", ErrNotFound
		}
		return strconv.Itoa(value), nil
	case "gauge":
		value, ok := Collector.storage.gauges[metricName]
		if !ok {
			return "", ErrNotFound
		}
		return value, nil
	default:
		return "", ErrNotImplemented
	}
}

func (c *collector) GetCounters() map[string]string {
	counters := make(map[string]string, 0)
	for name, value := range c.storage.counters {
		counters[name] = strconv.Itoa(value)
	}
	return counters
}

func (c *collector) GetGauges() map[string]string {
	gauges := make(map[string]string, 0)
	for name, value := range c.storage.gauges {
		gauges[name] = value
	}
	return gauges
}

func (c *collector) GetAvailableMetrics() []string {
	names := make([]string, 0)
	for cm := range c.storage.counters {
		names = append(names, cm)
	}
	for gm := range c.storage.gauges {
		names = append(names, gm)
	}
	return names
}

type MetricJSON struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

type collector struct {
	storage *memStorage
}

type memStorage struct {
	counters map[string]int
	gauges   map[string]string
}
