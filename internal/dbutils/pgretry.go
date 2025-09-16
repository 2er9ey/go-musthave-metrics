package dbutils

import (
	"context"
	"database/sql"
	"time"

	"github.com/2er9ey/go-musthave-metrics/internal/logger"
	"go.uber.org/zap"
)

func ExecContextWithRetry(ctx context.Context, db *sql.DB, maxRetries int, sql string, arguments ...interface{}) (sql.Result, error) {
	logger.Log.Debug("SQL:", zap.String("sql", sql))
	result, err := db.ExecContext(ctx, sql, arguments...)
	if err == nil {
		return result, err
	}

	classifier := NewPostgresErrorClassifier()
	classification := classifier.Classify(err)
	for attempt := 0; err != nil && classification == Retriable && attempt < maxRetries; attempt++ {
		time.Sleep(time.Duration(1+(attempt*2)) * time.Second)
		result, err = db.ExecContext(ctx, sql, arguments...)
		if err == nil {
			break
		}
		classification = classifier.Classify(err)
	}
	return result, err
}

func QueryRowContextWithRetry(ctx context.Context, db *sql.DB, maxRetries int, sql string, arguments ...interface{}) *sql.Row {
	logger.Log.Debug("SQL:", zap.String("sql", sql))
	row := db.QueryRowContext(ctx, sql, arguments...)
	if row.Err() == nil {
		return row
	}
	classifier := NewPostgresErrorClassifier()
	classification := classifier.Classify(row.Err())

	for attempt := 0; row.Err() != nil && classification == Retriable && attempt < maxRetries; attempt++ {
		time.Sleep(time.Duration(1+(attempt*2)) * time.Second)
		row = db.QueryRowContext(ctx, sql, arguments...)
		if row.Err() == nil {
			break
		}
		classification = classifier.Classify(row.Err())
	}
	return row
}

func QueryContextWithRetry(ctx context.Context, db *sql.DB, maxRetries int, sql string, arguments ...interface{}) (*sql.Rows, error) {
	logger.Log.Debug("SQL:", zap.String("sql", sql))
	rows, err := db.QueryContext(ctx, sql, arguments...)
	if err == nil {
		return rows, err
	}
	classifier := NewPostgresErrorClassifier()
	classification := classifier.Classify(err)
	for attempt := 0; err != nil && classification == Retriable && attempt < maxRetries; attempt++ {
		time.Sleep(time.Duration(1+(attempt*2)) * time.Second)
		rows, err = db.QueryContext(ctx, sql, arguments...)
		if err == nil {
			break
		}
		classification = classifier.Classify(err)
	}
	return rows, err
}
