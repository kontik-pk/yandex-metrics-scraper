package collector

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
)

var (
	ErrBadRequest     = errors.New("bad request")
	ErrNotImplemented = errors.New("not implemented")
	ErrNotFound       = errors.New("not found")
)

var Collector = collector{
	storage: &memStorage{
		Counters: make(map[string]int),
		Gauges:   make(map[string]string),
	},
}

func (c *collector) Collect(metricName string, metricType string, metricValue string) error {
	switch metricType {
	case "counter":
		value, err := strconv.Atoi(metricValue)
		if err != nil {
			return ErrBadRequest
		}
		c.storage.Counters[metricName] += value
	case "gauge":
		_, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			return ErrBadRequest
		}
		c.storage.Gauges[metricName] = metricValue
	default:
		return ErrNotImplemented
	}
	return nil
}

func (c *collector) Restore(filePath string) error {
	file, err := os.OpenFile(filePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		return err
	}
	data := scanner.Bytes()
	metricsFromFile := memStorage{}
	if err = json.Unmarshal(data, &metricsFromFile); err != nil {
		return err
	}
	c.decode(metricsFromFile)
	return nil
}

func (c *collector) Save(filePath string) error {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	writer := bufio.NewWriter(file)
	metricsData := c.encode()
	data, err := json.Marshal(&metricsData)
	if err != nil {
		return err
	}
	if _, err := writer.Write(data); err != nil {
		return err
	}
	if err := writer.WriteByte('\n'); err != nil {
		return err
	}
	if err := writer.Flush(); err != nil {
		return err
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
		value, ok := Collector.storage.Counters[metricName]
		if !ok {
			return "", ErrNotFound
		}
		return strconv.Itoa(value), nil
	case "gauge":
		value, ok := Collector.storage.Gauges[metricName]
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
	for name, value := range c.storage.Counters {
		counters[name] = strconv.Itoa(value)
	}
	return counters
}

func (c *collector) GetGauges() map[string]string {
	gauges := make(map[string]string, 0)
	for name, value := range c.storage.Gauges {
		gauges[name] = value
	}
	return gauges
}

func (c *collector) GetAvailableMetrics() []string {
	names := make([]string, 0)
	for cm := range c.storage.Counters {
		names = append(names, cm)
	}
	for gm := range c.storage.Gauges {
		names = append(names, gm)
	}
	return names
}

func (c *collector) encode() memStorage {
	return *c.storage
}

func (c *collector) decode(encoded memStorage) {
	c.storage = &encoded
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
	Counters map[string]int
	Gauges   map[string]string
}
