package repository

import "github.com/2er9ey/go-musthave-metrics/internal/models"

type MetricsRepositoryInterface interface {
	SetMetric(models.Metrics) error
	SetMetrics([]models.Metrics) error
	GetMetricString(string, string) (string, error)
	GetMetric(string, string) (models.Metrics, error)
	GetAllMetric() []models.Metrics
	Ping() (bool, error)
}
