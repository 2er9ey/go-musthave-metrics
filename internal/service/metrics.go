package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

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
			select {
			case <-ms.ctx.Done():
				return
			default:
				time.Sleep(time.Duration(ms.saveInterval) * time.Second)
			}
			logger.Log.Debug("Saving metrics")
			ms.SaveMetrics(ms.storageFilename)
		}
	}()
}

func (ms *MetricService) DBChekConnection() (bool, error) {
	ctx, cancel := context.WithTimeout(ms.ctx, 1*time.Second)
	defer cancel()

	fmt.Println("DBString = ", ms.databaseDSN)

	if err := ms.db.PingContext(ctx); err != nil {
		return false, err
	}
	return true, nil
}
