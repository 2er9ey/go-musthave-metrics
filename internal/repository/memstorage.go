package repository

import (
	"errors"
	"fmt"
	"maps"
	"slices"
	"sync"

	"github.com/2er9ey/go-musthave-metrics/internal/models"
)

type MemoryStorage struct {
	gaugeMutex     sync.RWMutex
	counterMutex   sync.RWMutex
	metricsGauge   map[string]models.Metrics
	metricsCounter map[string]models.Metrics
}

func NewMemoryStorage() (*MemoryStorage, error) {
	return &MemoryStorage{
		metricsGauge:   make(map[string]models.Metrics),
		metricsCounter: make(map[string]models.Metrics),
	}, nil
}

func (ms *MemoryStorage) PrintAll() {
	fmt.Println("Gauges:")
	for _, v := range ms.metricsGauge {
		fmt.Println(v)
	}
	fmt.Println("Counters:")
	for _, v := range ms.metricsCounter {
		fmt.Println(v)
	}
}

func (ms *MemoryStorage) SetMetric(m models.Metrics) error {
	switch m.MType {
	case models.Gauge:
		ms.gaugeMutex.Lock()
		ms.metricsGauge[m.ID] = m
		ms.gaugeMutex.Unlock()
	case models.Counter:
		ms.counterMutex.Lock()
		_, exists := ms.metricsCounter[m.ID]
		if exists {
			*(ms.metricsCounter[m.ID].Delta) += *(m.Delta)
		} else {
			ms.metricsCounter[m.ID] = m
		}
		ms.counterMutex.Unlock()
	default:
		return errors.New("invalid metric type")
	}
	return nil
}

func (ms *MemoryStorage) SetMetrics(metrics []models.Metrics) error {
	for _, m := range metrics {
		err := ms.SetMetric(m)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ms *MemoryStorage) GetMetric(metricKey string, metricType string) (models.Metrics, error) {
	var metric models.Metrics
	var exists bool
	switch metricType {
	case models.Gauge:
		ms.gaugeMutex.RLock()
		metric, exists = ms.metricsGauge[metricKey]
		ms.gaugeMutex.RUnlock()
	case models.Counter:
		ms.counterMutex.RLock()
		metric, exists = ms.metricsCounter[metricKey]
		ms.counterMutex.RUnlock()
	default:
		return metric, errors.New("invalid metric type")
	}

	if !exists {
		return metric, errors.New("metric does not exists")
	}

	return metric, nil
}

func (ms *MemoryStorage) GetMetricString(metricKey string, metricType string) (string, error) {
	metric, err := ms.GetMetric(metricKey, metricType)
	if err != nil {
		return "", errors.New("invalid metric type")
	}
	return metric.String(), nil
}

func (ms *MemoryStorage) GetAllMetric() []models.Metrics {
	ms.gaugeMutex.RLock()
	gauges := slices.Collect(maps.Values(ms.metricsGauge))
	ms.gaugeMutex.RUnlock()
	ms.counterMutex.RLock()
	counters := slices.Collect(maps.Values(ms.metricsCounter))
	ms.counterMutex.RUnlock()
	return append(gauges, counters...)
}

func (ms *MemoryStorage) Ping() (bool, error) {
	return true, nil
}
