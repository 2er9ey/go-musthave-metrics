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
		defer ps.Close()
	})

	t.Run("CreateIncorrectPostgreSQLStorage", func(t *testing.T) {
		ps := NewPostgreSQLStorage(ctx, "")
		assert.Nil(t, ps, "PostgreSQL storage must be nil")
	})

	t.Run("CreateIncorrectDSN", func(t *testing.T) {
		ps := NewPostgreSQLStorage(ctx, "host=10.0.0.1")
		assert.NotNil(t, ps, "Must be not nil")
		defer ps.Close()
	})
}

func TestSetMetric(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	metric := models.NewMetricCounter("counter", 134)
	ps := NewPostgreSQLStorage(ctx, "host=127.0.0.1 user=video password=XXXXXXXX dbname=video sslmode=disable")
	defer ps.Close()

	if err := ps.db.PingContext(ctx); err != nil {
		return
	}

	err := ps.Set(metric)
	if err != nil {
		t.Log(err)
	}
	assert.Nil(t, err, "Must be nil")
}

func TestCreateTables(t *testing.T) {
	// ps := NewPostgreSQLStorage(context.Background(), "host=127.0.0.1 user=video password=XXXXXXXX dbname=video sslmode=disable")
	// defer ps.Close()

	// err := ps.CreateTables()
	// assert.Nil(t, err, "Must be nil")
}

func TestGetAllMetric(t *testing.T) {
	// ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	// defer cancel()
	// ps := NewPostgreSQLStorage(ctx, "host=127.0.0.1 user=video password=XXXXXXXX dbname=video sslmode=disable")
	// defer ps.Close()

	// metrics := ps.GetAllMetric()
	// assert.NotNil(t, err, "Must be not nil")
}

func TestSaveMetric(t *testing.T) {
	// ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	// defer cancel()
	// ps := NewPostgreSQLStorage(ctx, "host=127.0.0.1 user=video password=XXXXXXXX dbname=video sslmode=disable")
	// defer ps.Close()

	// err := ps.SaveMetrics("./x.dat")
	// assert.Nil(t, err, "Must be nil")
}
