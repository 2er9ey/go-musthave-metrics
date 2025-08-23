package repository

import "github.com/2er9ey/go-musthave-metrics/internal/models"

type MetricsRepositoryInterface interface {
	Set(m models.Metrics) error
	GetString(metricKey string, metricType string) (string, error)
	GetMetric(metricKey string, metricType string) (models.Metrics, error)
	GetAllMetric() []models.Metrics
	LoadMetrics(string) error
	SaveMetrics(string) error
}
