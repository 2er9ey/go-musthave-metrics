package repository

import (
	"context"
	"testing"
	"time"

	"github.com/2er9ey/go-musthave-metrics/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestCreatePostgreSQLStorage(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	t.Run("CreateCorrectPostgreSQLStorage", func(t *testing.T) {
		ps := NewPostgreSQLStorage(ctx, "host=127.0.0.1 user=video password=XXXXXXXX dbname=video sslmode=disable")
		assert.NotNil(t, ps, "PostgreSQL storage is nil")
	})

	t.Run("CreateIncorrectPostgreSQLStorage", func(t *testing.T) {
		ps := NewPostgreSQLStorage(ctx, "")
		assert.Nil(t, ps, "PostgreSQL storage must be nil")
	})

	t.Run("CreateIncorrectDSN", func(t *testing.T) {
		ps := NewPostgreSQLStorage(ctx, "host=10.0.0.1")
		assert.NotNil(t, ps, "Must be not nil")
	})
}

func TestSetMetric(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	metric := models.NewMetricCounter("counter", 134)
	ps := NewPostgreSQLStorage(ctx, "host=127.0.0.1 user=video password=XXXXXXXX dbname=video sslmode=disable")

	err := ps.Set(metric)
	assert.NotNil(t, err, "Must be not nil")
}
