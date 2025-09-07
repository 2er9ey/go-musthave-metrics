package repository

import "github.com/2er9ey/go-musthave-metrics/internal/models"

type MetricsRepositoryInterface interface {
	Set(models.Metrics) error
	SetBunch([]models.Metrics) error
	GetString(string, string) (string, error)
	GetMetric(string, string) (models.Metrics, error)
	GetAllMetric() []models.Metrics
	LoadMetrics(string) error
	SaveMetrics(string) error
}
