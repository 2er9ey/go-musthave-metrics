package repository

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"os"
	"slices"
	"strings"

	"github.com/2er9ey/go-musthave-metrics/internal/models"
)

type MemoryStorage struct {
	metricsGauge   map[string]models.Metrics
	metricsCounter map[string]models.Metrics
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		metricsGauge:   make(map[string]models.Metrics),
		metricsCounter: make(map[string]models.Metrics),
	}
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

func (ms *MemoryStorage) Set(m models.Metrics) error {
	switch m.MType {
	case models.Gauge:
		ms.metricsGauge[m.ID] = m
	case models.Counter:
		_, exists := ms.metricsCounter[m.ID]
		if exists {
			*(ms.metricsCounter[m.ID].Delta) += *(m.Delta)
		} else {
			ms.metricsCounter[m.ID] = m
		}
	default:
		return errors.New("invalid metric type")
	}
	return nil
}

func (ms *MemoryStorage) GetMetric(metricKey string, metricType string) (models.Metrics, error) {
	var metric models.Metrics
	var exists bool
	switch metricType {
	case models.Gauge:
		metric, exists = ms.metricsGauge[metricKey]
	case models.Counter:
		metric, exists = ms.metricsCounter[metricKey]
	default:
		return metric, errors.New("invalid metric type")
	}

	if !exists {
		return metric, errors.New("metric does not exists")
	}

	return metric, nil
}

func (ms *MemoryStorage) GetString(metricKey string, metricType string) (string, error) {
	var metric models.Metrics
	var exists bool
	switch metricType {
	case models.Gauge:
		metric, exists = ms.metricsGauge[metricKey]
	case models.Counter:
		metric, exists = ms.metricsCounter[metricKey]
	default:
		return "", errors.New("invalid metric type")
	}
	if !exists {
		return "", errors.New("metric does not exists")
	}
	return metric.String(), nil
}

func (ms *MemoryStorage) GetAllMetric() []models.Metrics {
	gauges := slices.Collect(maps.Values(ms.metricsGauge))
	counters := slices.Collect(maps.Values(ms.metricsCounter))
	return append(gauges, counters...)
}

func (ms *MemoryStorage) LoadMetrics(filename string) error {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0666)
	if err != nil {
		return errors.New("can't open file")
	}

	buf := bufio.NewReader(file)
	line, err := buf.ReadBytes('\n')
	if err != nil {
		return err
	}

	if string(line) != "[\n" {
		return errors.New("wrong file")
	}

	for {
		var metric models.Metrics
		line, err = buf.ReadBytes('\n')
		if err != nil {
			return err
		}
		stringLine := strings.TrimRight(strings.TrimSpace(string(line)), ",")
		if stringLine == "]" {
			break
		}
		json.Unmarshal([]byte(stringLine), &metric)
		switch metric.MType {
		case models.Counter:
			ms.metricsCounter[metric.ID] = metric
		case models.Gauge:
			ms.metricsGauge[metric.ID] = metric
		default:
			return errors.New("wrong file")
		}
	}
	return nil
}

func (ms *MemoryStorage) SaveMetrics(filename string) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	file.WriteString("[\n")

	for _, value := range ms.metricsCounter {
		jsonString, _ := json.Marshal(value)
		file.WriteString("  ")
		file.Write(jsonString)
		file.WriteString(",\n")
	}
	for _, value := range ms.metricsGauge {
		jsonString, _ := json.Marshal(value)
		file.WriteString("  ")
		file.Write(jsonString)
		file.WriteString(",\n")
	}
	file.WriteString("]\n")
	file.Close()
	return nil
}
