package models

type Metric struct {
	ID    string   `json:"id" db:"id"`
	MType string   `json:"type" db:"m_type"`
	Delta *int64   `json:"delta,omitempty" db:"m_delta"`
	Value *float64 `json:"value,omitempty" db:"m_value"`
}
