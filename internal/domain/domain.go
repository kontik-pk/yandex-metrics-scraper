package domain

type Metric struct {
	Value any
	MType string
}

type MemStorage struct {
	Metrics map[string]Metric
}
