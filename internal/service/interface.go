package service

import "github.com/2er9ey/go-musthave-metrics/internal/models"

type MetricServiceInterface interface {
	Set(string, string, string) error
	Get(string, string) (string, error)
	GetMetric(string, string) (models.Metrics, error)
	GetAll() []models.Metrics
}
