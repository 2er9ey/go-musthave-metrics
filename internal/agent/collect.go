package agent

import (
	"errors"
	"math/rand"
	"reflect"
	"runtime"

	"github.com/2er9ey/go-musthave-metrics/internal/models"
	"github.com/2er9ey/go-musthave-metrics/internal/repository"
)

func CollectorMetrics(repo repository.MetricsRepositoryInterface, collectMetrics *[]CollectMetric) error {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	for _, v := range *collectMetrics {
		var metric models.Metrics
		err := error(nil)
		switch v.collectType {
		case collectTypeMemStat:
			metric, err = collectGetMemStatMectric(ms, v.ID, v.MType)
		case collectTypeRandom:
			metric, err = collectGetRandomMectric(v.ID, v.MType)
		case collectTypeConst1:
			metric, err = collectGetConst1Mectric(v.ID, v.MType)
		default:
			err = errors.New("invalid collect type")
		}
		if err == nil {
			repo.Set(metric)
		} else {
			return err
		}
	}
	return nil
}

func collectGetMemStatMectric(ms runtime.MemStats, ID string, MType string) (models.Metrics, error) {
	var metric models.Metrics
	rms := reflect.ValueOf(ms)
	memField := rms.FieldByName(ID)

	if memField.IsValid() {
		switch MType {
		case models.Gauge:
			if memField.Kind() == reflect.Float32 || memField.Kind() == reflect.Float64 {
				metric = models.NewMetricGauge(ID, memField.Float())
			} else {
				metric = models.NewMetricGauge(ID, float64(memField.Uint()))
			}
		case models.Counter:
			if memField.Kind() == reflect.Float32 || memField.Kind() == reflect.Float64 {
				metric = models.NewMetricCounter(ID, int64(memField.Float()))
			} else {
				metric = models.NewMetricCounter(ID, int64(memField.Uint()))
			}
		default:
			return metric, errors.New("invalid metric type")
		}
		return metric, nil
	}
	return metric, errors.New("metric not found")
}

func collectGetRandomMectric(ID string, MType string) (models.Metrics, error) {
	var metric models.Metrics
	switch MType {
	case models.Gauge:
		metric = models.NewMetricGauge(ID, rand.Float64())
	case models.Counter:
		metric = models.NewMetricCounter(ID, rand.Int63())
	default:
		return metric, errors.New("invalid metric type")
	}
	return metric, nil
}

func collectGetConst1Mectric(ID string, MType string) (models.Metrics, error) {
	var metric models.Metrics
	switch MType {
	case models.Gauge:
		metric = models.NewMetricGauge(ID, float64(1))
	case models.Counter:
		metric = models.NewMetricCounter(ID, int64(1))
	default:
		return metric, errors.New("invalid metric type")
	}
	return metric, nil
}
