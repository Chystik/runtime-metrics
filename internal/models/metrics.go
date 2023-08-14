package models

type Metric struct {
	Name string
	MetricValue
}

type MetricValue struct {
	Gauge
	Counter
}

type Gauge float64
type Counter int64
