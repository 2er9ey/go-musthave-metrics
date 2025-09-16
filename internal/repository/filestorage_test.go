package repository

import (
	"os"
	"testing"
	"time"

	"github.com/2er9ey/go-musthave-metrics/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestFileStorageCreateStorage1(t *testing.T) {
	fs, _ := NewFileStorage("", 0, false)
	assert.Nil(t, fs, "FileStrorage storage must be nil")
}

func TestFileStorageCreateStorage2(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "tmpfile-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	t.Logf("Temp filename: %s", tmpFile.Name())

	fs, _ := NewFileStorage(tmpFile.Name(), 0, false)
	assert.NotNil(t, fs, "FileStrorage storage must not be nil")
}

func TestFileStorageSaveMetrics1(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "tmpfile-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	t.Logf("Temp filename: %s", tmpFile.Name())

	fs, _ := NewFileStorage(tmpFile.Name(), 0, false)
	assert.NotNil(t, fs, "FileStrorage storage must not be nil")

	err = fs.SaveMetrics()
	assert.Nil(t, err, "Error saving file")
}

func TestFileStorageRunSaver1(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "tmpfile-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	t.Logf("Temp filename: %s", tmpFile.Name())

	fs, _ := NewFileStorage(tmpFile.Name(), 3, false)
	assert.NotNil(t, fs, "FileStrorage storage must not be nil")

	metric := models.NewMetricCounter("123", 345)

	fs.SetMetric(metric)

	time.Sleep(time.Duration(10) * time.Second)
}
