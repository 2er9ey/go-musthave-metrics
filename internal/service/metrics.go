package service

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"

	"github.com/2er9ey/go-musthave-metrics/internal/logger"
	"github.com/2er9ey/go-musthave-metrics/internal/models"
	"github.com/2er9ey/go-musthave-metrics/internal/repository"
)

type MetricService struct {
	ctx             context.Context
	repo            repository.MetricsRepositoryInterface
	saveInterval    int
	storageFilename string
	databaseDSN     string
	db              *sql.DB
}

func NewMetricService(ctx context.Context, repo repository.MetricsRepositoryInterface, saveInterval int,
	storageFilename string, databaseDSN string) *MetricService {
	db, err := sql.Open("pgx", databaseDSN)
	if err != nil {
		panic(err)
	}
	return &MetricService{
		ctx:             ctx,
		repo:            repo,
		saveInterval:    saveInterval,
		storageFilename: storageFilename,
		databaseDSN:     databaseDSN,
		db:              db,
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
		ms.repo.SetMetric(metric)
		return nil
	}
	return err
}

func (ms *MetricService) Get(mID string, mType string) (string, error) {
	return ms.repo.GetMetricString(mID, mType)
}

func (ms *MetricService) GetMetric(mID string, mType string) (models.Metrics, error) {
	logger.Log.Debug("Service: GetMetric", zap.String("mID", mID), zap.String("mType", mType))
	return ms.repo.GetMetric(mID, mType)
}

func (ms *MetricService) GetAll() []models.Metrics {
	return ms.repo.GetAllMetric()
}

func (ms *MetricService) SetBunch(metrics []models.Metrics) error {
	return ms.repo.SetMetrics(metrics)
}

func (ms *MetricService) Ping() (bool, error) {
	return ms.repo.Ping()
}
