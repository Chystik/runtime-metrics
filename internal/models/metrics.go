package models

type (
	Metric struct {
		Name string
		MetricValue
	}

	MetricValue struct {
		Gauge
		Counter
	}

	Gauge   float64
	Counter int64
)
