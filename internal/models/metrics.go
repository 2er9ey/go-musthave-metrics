package models

import (
	"strconv"
)

const (
	Counter = "counter"
	Gauge   = "gauge"
)

// NOTE: Не усложняем пример, вводя иерархическую вложенность структур.
// Органичиваясь плоской моделью.
// Delta и Value объявлены через указатели,
// что бы отличать значение "0", от не заданного значения
// и соответственно не кодировать в структуру.
type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
	Hash  string   `json:"hash,omitempty"`
}

func NewMetricGauge(key string, v float64) Metrics {
	return Metrics{ID: key, MType: Gauge, Value: &v}
}

func NewMetricCounter(key string, v int64) Metrics {
	return Metrics{ID: key, MType: Counter, Delta: &v}
}

func (m Metrics) String() string {
	res := "nil"
	switch m.MType {
	case Gauge:
		if m.Value != nil {
			res = strconv.FormatFloat(*(m.Value), 'f', -1, 64)
		}
	case Counter:
		if m.Delta != nil {
			res = strconv.FormatInt(*(m.Delta), 10)
		}
	}
	return res
}
