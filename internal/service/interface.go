package service

import "github.com/2er9ey/go-musthave-metrics/internal/models"

type MetricServiceInterface interface {
	Set(string, string, string) error
	SetBunch([]models.Metrics) error
	Get(string, string) (string, error)
	GetMetric(string, string) (models.Metrics, error)
	GetAll() []models.Metrics
	Ping() (bool, error)
	// LoadMetrics(string) error
	// SaveMetrics(string) error
	// RunSaver()
}
