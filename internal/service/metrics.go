package service

import (
	"errors"
	"strconv"
	"time"

	"github.com/2er9ey/go-musthave-metrics/internal/logger"
	"github.com/2er9ey/go-musthave-metrics/internal/models"
	"github.com/2er9ey/go-musthave-metrics/internal/repository"
)

type MetricService struct {
	repo            repository.MetricsRepositoryInterface
	saveInterval    int
	storageFilename string
}

func NewMetricService(repo repository.MetricsRepositoryInterface, saveInterval int,
	storageFilename string) *MetricService {
	return &MetricService{
		repo:            repo,
		saveInterval:    saveInterval,
		storageFilename: storageFilename,
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
		err = errors.New("invalid metric type (" + mType + ")")
	}
	if err == nil {
		ms.repo.Set(metric)
		if ms.saveInterval == 0 {
			ms.repo.SaveMetrics(ms.storageFilename)
		}
		return nil
	}
	return err
}

func (ms *MetricService) Get(mID string, mType string) (string, error) {
	return ms.repo.GetString(mID, mType)
}

func (ms *MetricService) GetMetric(mID string, mType string) (models.Metrics, error) {
	return ms.repo.GetMetric(mID, mType)
}

func (ms *MetricService) GetAll() []models.Metrics {
	return ms.repo.GetAllMetric()
}

func (ms *MetricService) LoadMetrics(filename string) error {
	return ms.repo.LoadMetrics(filename)
}

func (ms *MetricService) SaveMetrics(filename string) error {
	return ms.repo.SaveMetrics(filename)
}

func (ms *MetricService) RunSaver() {
	if ms.saveInterval <= 0 {
		return
	}
	go func() {
		for {
			time.Sleep(time.Duration(ms.saveInterval) * time.Second)
			logger.Log.Debug("Saving metrics")
			ms.SaveMetrics(ms.storageFilename)
		}
	}()
}
