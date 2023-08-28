package models

type Metric struct {
	Name  string `json:"id"`
	MType string `json:"type"`
	MetricValue
}

type MetricValue struct {
	Gauge   `json:"delta,omitempty"`
	Counter `json:"value,omitempty"`
}

type Gauge float64
type Counter int64
