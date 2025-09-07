package repository

import (
	"context"
	"database/sql"

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
	defer ps.db.Close()

	return ps
}

func (ms *PostreSQLStorage) PrintAll() {
}

func (ms *PostreSQLStorage) Set(m models.Metrics) error {
	_, err := ms.db.ExecContext(ms.ctx, "insert into metrics (metric_id, metric_type, metric_delta, metric_value, metric_hash) values (?,?,?,?,?)",
		m.ID, m.MType, m.Delta, m.Value, m.Hash)
	if err != nil {
		return err
	}
	return nil
}

func (ms *PostreSQLStorage) GetMetric(metricKey string, metricType string) (models.Metrics, error) {
	return models.Metrics{}, nil

}

func (ms *PostreSQLStorage) GetString(metricKey string, metricType string) (string, error) {
	return "", nil
}

func (ms *PostreSQLStorage) GetAllMetric() []models.Metrics {
	rows, err := ms.db.QueryContext(ms.ctx, "SELECT metric_id, metric_type, metric_delta, metric_value, metric_hash from metrics ORDER BY id")
	if err != nil {
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

func (ms *PostreSQLStorage) LoadMetrics(filename string) error {
	return nil
}

func (ms *PostreSQLStorage) SaveMetrics(filename string) error {
	return nil
}
