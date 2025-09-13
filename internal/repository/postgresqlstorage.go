package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"

	"github.com/2er9ey/go-musthave-metrics/internal/logger"
	"github.com/2er9ey/go-musthave-metrics/internal/models"
	"github.com/2er9ey/go-musthave-metrics/internal/pgerrors"
)

type PostreSQLStorage struct {
	databaseDSN string
	db          *sql.DB
	ctx         context.Context
}

func NewPostgreSQLStorage(ctx context.Context, databaseDSN string) *PostreSQLStorage {
	if databaseDSN == "" {
		return nil
	}
	ps := &PostreSQLStorage{}
	db, err := sql.Open("pgx", databaseDSN)
	if err != nil {
		panic(err)
	}
	ps.db = db
	ps.databaseDSN = databaseDSN
	ps.ctx = ctx
	return ps
}

func (ms *PostreSQLStorage) Close() {
	ms.db.Close()
}

func (ms *PostreSQLStorage) PrintAll() {
}

func (ms *PostreSQLStorage) SetMetric(m models.Metrics) error {
	logger.Log.Debug("PSQL: insert " + m.String())
	sql := "INSERT INTO metrics (metric_id, metric_type, metric_delta, metric_value, metric_hash) values ( $1 , $2, $3, $4, $5 ) on conflict (metric_id, metric_type) do update set "
	switch m.MType {
	case models.Gauge:
		sql = sql + "metric_value = EXCLUDED.metric_value"
	case models.Counter:
		sql = sql + "metric_delta = metrics.metric_delta + EXCLUDED.metric_delta"
	default:
		return errors.New("invalid metric type")
	}

	sql = sql + ";"

	logger.Log.Debug("SQL:", zap.String("sql", sql))
	classifier := pgerrors.NewPostgresErrorClassifier()
	for attempt := 0; attempt < 4; attempt++ {
		_, err := ms.db.ExecContext(ms.ctx, sql, m.ID, m.MType, m.Delta, m.Value, m.Hash)
		if err == nil {
			break
		}

		classification := classifier.Classify(err)
		if classification == pgerrors.NonRetriable {
			// Нет смысла повторять, возвращаем ошибку
			fmt.Printf("Непредвиденная ошибка: %v\n", err)
			return err
		}
		time.Sleep(time.Duration(1+(attempt*2)) * time.Second)
	}
	return nil
}

func (ms *PostreSQLStorage) SetMetrics(metrics []models.Metrics) error {
	logger.Log.Debug("PSQL: SetMetrics")
	tx, err := ms.db.Begin()
	if err != nil {
		return err
	}
	for _, m := range metrics {
		err := ms.SetMetric(m)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (ms *PostreSQLStorage) GetMetric(metricKey string, metricType string) (models.Metrics, error) {
	logger.Log.Debug("PSQL: select " + metricKey + "/" + metricType)
	classifier := pgerrors.NewPostgresErrorClassifier()
	var metric models.Metrics
	for attempt := 0; attempt < 4; attempt++ {
		row := ms.db.QueryRowContext(ms.ctx, "SELECT metric_id, metric_type, metric_delta, metric_value, metric_hash from metrics where metric_id = $1 and metric_type = $2", metricKey, metricType)
		if row.Err() == nil {
			err := row.Scan(&metric.ID, &metric.MType, &metric.Delta, &metric.Value, &metric.Hash)
			if err != nil {
				return models.Metrics{}, err
			}
			break
		}
		classification := classifier.Classify(row.Err())
		if classification == pgerrors.NonRetriable || attempt == 3 {
			return models.Metrics{}, row.Err()
		}
		time.Sleep(time.Duration(1+(attempt*2)) * time.Second)
	}

	return metric, nil

}

func (ms *PostreSQLStorage) GetMetricString(metricKey string, metricType string) (string, error) {
	return "", nil
}

func (ms *PostreSQLStorage) GetAllMetric() []models.Metrics {
	rows, err := ms.db.QueryContext(ms.ctx, "SELECT metric_id, metric_type, metric_delta, metric_value, metric_hash from metrics ORDER BY id")
	if err != nil || rows.Err() != nil {
		return nil
	}

	defer rows.Close()

	metrics := make([]models.Metrics, 0)

	for rows.Next() {
		var metric models.Metrics
		err = rows.Scan(&metric.ID, &metric.MType, &metric.Delta, &metric.Value, &metric.Hash)
		if err != nil {
			return nil
		}

		metrics = append(metrics, metric)
	}

	return metrics
}

func (ms *PostreSQLStorage) Ping() (bool, error) {
	ctx, cancel := context.WithTimeout(ms.ctx, 1*time.Second)
	defer cancel()

	if err := ms.db.PingContext(ctx); err != nil {
		return false, err
	}
	return true, nil
}
