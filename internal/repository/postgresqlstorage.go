package repository

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/2er9ey/go-musthave-metrics/internal/models"
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

func (ms *PostreSQLStorage) CreateTables() error {
	_, err := ms.db.ExecContext(ms.ctx, "CREATE TABLE IF NOT EXISTS metrics (id SERIAL PRIMARY KEY, metric_id VARCHAR(255) NOT NULL, metric_type  VARCHAR(255) NOT NULL, metric_delta bigint, metric_value double precision, metric_hash  VARCHAR(255));")
	if err != nil {
		return err
	}
	return nil
}

func (ms *PostreSQLStorage) PrintAll() {
}

func (ms *PostreSQLStorage) SetMetric(m models.Metrics) error {
	_, err := ms.db.ExecContext(ms.ctx, "INSERT INTO metrics (metric_id, metric_type, metric_delta, metric_value, metric_hash) values ( $1 , $2, $3, $4, $5 );", m.ID, m.MType, m.Delta, m.Value, m.Hash)
	if err != nil {
		return err
	}
	return nil
}

func (ms *PostreSQLStorage) SetMetrics(metrics []models.Metrics) error {
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
	row := ms.db.QueryRowContext(ms.ctx, "SELECT metric_id, metric_type, metric_delta, metric_value, metric_hash from metrics where metric_id = $1 and metric_type = $2", metricKey, metricType)
	if row.Err() != nil {
		return models.Metrics{}, row.Err()
	}
	var metric models.Metrics
	err := row.Scan(&metric.ID, &metric.MType, &metric.Delta, &metric.Value, &metric.Hash)
	if err != nil {
		return models.Metrics{}, err
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
