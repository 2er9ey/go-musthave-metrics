package repository

import "github.com/2er9ey/go-musthave-metrics/internal/models"

type MetricsRepositoryInterface interface {
	Set(m models.Metrics) error
	GetString(metricKey string) (string, error)
	GetMetric(metricKey string) (models.Metrics, error)
	GetAllMetric() map[string]models.Metrics
}
