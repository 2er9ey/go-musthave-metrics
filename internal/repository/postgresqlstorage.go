package repository

import (
	"bufio"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"os"
	"strings"

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

func (ms *PostreSQLStorage) Set(m models.Metrics) error {
	_, err := ms.db.ExecContext(ms.ctx, "INSERT INTO metrics (metric_id, metric_type, metric_delta, metric_value, metric_hash) values ( $1 , $2, $3, $4, $5 );", m.ID, m.MType, m.Delta, m.Value, m.Hash)
	if err != nil {
		return err
	}
	return nil
}

func (ms *PostreSQLStorage) SetBunch(metrics []models.Metrics) error {
	tx, err := ms.db.Begin()
	if err != nil {
		return err
	}
	for _, m := range metrics {
		err := ms.Set(m)
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

func (ms *PostreSQLStorage) GetString(metricKey string, metricType string) (string, error) {
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

func (ms *PostreSQLStorage) LoadMetrics(filename string) error {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0666)
	if err != nil {
		return errors.New("can't open file")
	}

	buf := bufio.NewReader(file)
	line, err := buf.ReadBytes('\n')
	if err != nil {
		return err
	}

	if string(line) != "[\n" {
		return errors.New("wrong file")
	}

	var metrics []models.Metrics

	for {
		var metric models.Metrics
		line, err = buf.ReadBytes('\n')
		if err != nil {
			return err
		}
		stringLine := strings.TrimRight(strings.TrimSpace(string(line)), ",")
		if stringLine == "]" {
			break
		}
		json.Unmarshal([]byte(stringLine), &metric)
		metrics = append(metrics, metric)
	}
	return ms.SetBunch(metrics)
}

func (ms *PostreSQLStorage) SaveMetrics(filename string) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	metrics := ms.GetAllMetric()

	file.WriteString("[\n")

	for _, value := range metrics {
		jsonString, _ := json.Marshal(value)
		file.WriteString("  ")
		file.Write(jsonString)
		file.WriteString(",\n")
	}
	file.WriteString("]\n")
	file.Close()
	return nil
}
