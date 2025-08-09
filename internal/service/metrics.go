package service

import (
	"os"
	"strconv"

	"github.com/2er9ey/go-musthave-metrics/internal/models"
	"github.com/2er9ey/go-musthave-metrics/internal/repository"
)

type MetricService struct {
	repo repository.MetricsRepositoryInterface
}

func NewMetricService(repo repository.MetricsRepositoryInterface) *MetricService {
	return &MetricService{
		repo: repo,
	}
}

func (ms *MetricService) Set(mID string, mType string, mValue string) error {
	var metric models.Metrics
	var err error
	switch mType {
	case models.Counter:
		value, errConv := strconv.ParseInt(mValue, 10, 64)
		if errConv == nil {
			metric = models.NewMetricCounter(mID, value)
		} else {
			err = errConv
		}
	case models.Gauge:
		value, errConv := strconv.ParseFloat(mValue, 64)
		if errConv == nil {
			metric = models.NewMetricGauge(mID, value)
		} else {
			err = errConv
		}
	default:
		err = os.ErrInvalid
	}
	if err == nil {
		return ms.repo.Set(metric)
	}
	return err
}

func (ms *MetricService) Get(mID string, mType string) (string, error) {
	return ms.repo.GetString(mID, mType)
}

func (ms *MetricService) GetAll() []models.Metrics {
	return ms.repo.GetAllMetric()
}
