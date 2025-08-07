package repository

import (
	"fmt"
	"os"

	"github.com/2er9ey/go-musthave-metrics/internal/models"
)

type MemoryStorage struct {
	metrics map[string]models.Metrics
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		metrics: make(map[string]models.Metrics),
	}
}

func (ms *MemoryStorage) PrintAll() {
	for _, v := range ms.metrics {
		fmt.Println(v)
	}
}

func (ms *MemoryStorage) Set(m models.Metrics) error {
	metric, exists := ms.metrics[m.ID]
	if exists {
		if metric.MType != m.MType {
			return os.ErrInvalid
		}
		if metric.MType == models.Gauge {
			ms.metrics[m.ID] = m
		} else {
			*(ms.metrics[m.ID].Delta) += *(m.Delta)
		}
	} else {
		ms.metrics[m.ID] = m
	}
	return nil
}

func (ms *MemoryStorage) GetMetric(metricKey string) (models.Metrics, error) {
	metric, exists := ms.metrics[metricKey]
	if !exists {
		return metric, os.ErrNotExist
	}
	return metric, nil
}

func (ms *MemoryStorage) GetString(metricKey string) (string, error) {
	metric, exists := ms.metrics[metricKey]
	if !exists {
		return metric.String(), os.ErrNotExist
	}
	return metric.String(), nil
}

func (ms *MemoryStorage) GetAllMetric() map[string]models.Metrics {
	return ms.metrics
}
