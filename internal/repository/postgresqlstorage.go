package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/2er9ey/go-musthave-metrics/internal/dbutils"
	"github.com/2er9ey/go-musthave-metrics/internal/logger"
	"github.com/2er9ey/go-musthave-metrics/internal/models"
)

type PostgreSQLStorage struct {
	databaseDSN string
	db          *sql.DB
	ctx         context.Context
}

func NewPostgreSQLStorage(ctx context.Context, databaseDSN string) (*PostgreSQLStorage, error) {
	if databaseDSN == "" {
		return nil, errors.New("wrong database dsn")
	}
	ps := &PostgreSQLStorage{}
	db, err := sql.Open("pgx", databaseDSN)
	if err != nil {
		return nil, err
	}
	ps.db = db
	ps.databaseDSN = databaseDSN
	ps.ctx = ctx
	return ps, nil
}

func (ms *PostgreSQLStorage) Close() {
	ms.db.Close()
}

func (ms *PostgreSQLStorage) PrintAll() {
}

func (ms *PostgreSQLStorage) SetMetric(m models.Metrics) error {
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

	_, err := dbutils.ExecContextWithRetry(ms.ctx, ms.db, 4, sql, m.ID, m.MType, m.Delta, m.Value, m.Hash)

	return err
}

func (ms *PostgreSQLStorage) SetMetrics(metrics []models.Metrics) error {
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

func (ms *PostgreSQLStorage) GetMetric(metricKey string, metricType string) (models.Metrics, error) {
	logger.Log.Debug("PSQL: select " + metricKey + "/" + metricType)
	var metric models.Metrics
	row := dbutils.QueryRowContextWithRetry(ms.ctx, ms.db, 4, "SELECT metric_id, metric_type, metric_delta, metric_value, metric_hash from metrics where metric_id = $1 and metric_type = $2", metricKey, metricType)
	if row.Err() == nil {
		err := row.Scan(&metric.ID, &metric.MType, &metric.Delta, &metric.Value, &metric.Hash)
		if err != nil {
			return models.Metrics{}, err
		}
	}
	return metric, nil
}

func (ms *PostgreSQLStorage) GetMetricString(metricKey string, metricType string) (string, error) {
	metric, err := ms.GetMetric(metricKey, metricType)
	if err == nil {
		return metric.String(), err
	}
	return "", err
}

func (ms *PostgreSQLStorage) GetAllMetric() []models.Metrics {
	rows, err := dbutils.QueryContextWithRetry(ms.ctx, ms.db, 4, "SELECT metric_id, metric_type, metric_delta, metric_value, metric_hash from metrics ORDER BY id")
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

func (ms *PostgreSQLStorage) Ping() (bool, error) {
	ctx, cancel := context.WithTimeout(ms.ctx, 1*time.Second)
	defer cancel()

	if err := ms.db.PingContext(ctx); err != nil {
		return false, err
	}
	return true, nil
}
