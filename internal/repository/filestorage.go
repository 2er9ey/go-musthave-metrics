package repository

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/2er9ey/go-musthave-metrics/internal/models"
)

type FileStorage struct {
	mutex         sync.Mutex
	ms            *MemoryStorage
	fileName      string
	storeInterval int
}

func NewFileStorage(fileName string, storeInterval int, restoreMetric bool) (*FileStorage, error) {
	if fileName == "" {
		return nil, errors.New("filename must not be empty")
	}
	ms, err := NewMemoryStorage()
	if err != nil {
		return nil, errors.New("out of memory")
	}

	fs := &FileStorage{
		ms:            ms,
		fileName:      fileName,
		storeInterval: storeInterval,
	}

	if restoreMetric {
		fs.LoadMetrics(fileName)
		// if err != nil {
		// 	return nil, err
		// }
	}
	if storeInterval > 0 {
		fs.RunSaver()
	}
	return fs, err
}

func (fs *FileStorage) SetMetric(m models.Metrics) error {
	err := fs.ms.SetMetric(m)
	if err != nil {
		return err
	}
	if fs.storeInterval == 0 {
		fs.SaveMetrics()
	}
	return nil
}

func (fs *FileStorage) SetMetrics(metrics []models.Metrics) error {
	err := fs.ms.SetMetrics(metrics)
	if err != nil {
		return err
	}
	if fs.storeInterval == 0 {
		fs.SaveMetrics()
	}
	return nil
}

func (fs *FileStorage) GetMetric(metricKey string, metricType string) (models.Metrics, error) {
	if fs == nil {
		os.Exit(10)
	}
	if fs.ms == nil {
		return models.Metrics{}, nil
	}
	return fs.ms.GetMetric(metricKey, metricType)
}

func (fs *FileStorage) GetMetricString(metricKey string, metricType string) (string, error) {
	if fs.ms == nil {
		return "", nil
	}
	return fs.ms.GetMetricString(metricKey, metricType)
}

func (fs *FileStorage) GetAllMetric() []models.Metrics {
	if fs.ms == nil {
		return make([]models.Metrics, 0)
	}
	return fs.ms.GetAllMetric()
}

func (fs *FileStorage) LoadMetrics(filename string) error {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()

	file, err := os.OpenFile(filename, os.O_RDONLY, 0666)
	if err != nil {
		return errors.New("can't open file " + filename)
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
	return fs.ms.SetMetrics(metrics)
}

func (fs *FileStorage) SaveMetrics() error {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()

	file, err := os.OpenFile(fs.fileName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	metrics := fs.ms.GetAllMetric()

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

func (fs *FileStorage) RunSaver() {
	if fs.storeInterval <= 0 {
		return
	}
	go func() {
		for {
			time.Sleep(time.Duration(fs.storeInterval) * time.Second)
			fs.SaveMetrics()
		}
		//			logger.Log.Debug("Saving metrics")
	}()
}

func (fs *FileStorage) Ping() (bool, error) {
	return true, nil
}
