package agent

import (
	"errors"
	"math/rand"
	"reflect"
	"runtime"
	"strconv"
	"time"

	"github.com/2er9ey/go-musthave-metrics/internal/models"
	"github.com/2er9ey/go-musthave-metrics/internal/repository"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

func CollectorMetrics(repo repository.MetricsRepositoryInterface, collectMetrics *[]CollectMetric) error {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	for _, v := range *collectMetrics {
		var metrics []models.Metrics
		var metric models.Metrics
		err := error(nil)
		switch v.collectType {
		case collectTypeMemStat:
			metric, err = collectGetMemStatMectric(ms, v.ID, v.MType)
			metrics = []models.Metrics{metric}
		case collectTypeRandom:
			metric, err = collectGetRandomMectric(v.ID, v.MType)
			metrics = []models.Metrics{metric}
		case collectTypeConst1:
			metric, err = collectGetConst1Mectric(v.ID, v.MType)
			metrics = []models.Metrics{metric}
		case collectTypePsUtils:
			metrics, err = collectGetPSUtilsMetrics(v.ID)
		default:
			err = errors.New("invalid collect type")
		}
		if err == nil {
			repo.SetMetrics(metrics)
		} else {
			return err
		}
	}
	return nil
}

func collectGetPSUtilsMetrics(ID string) ([]models.Metrics, error) {
	var metrics []models.Metrics
	err := error(nil)
	v, err := mem.VirtualMemory()
	if err != nil {
		err = errors.New("error getting virtual memory info")
	}
	switch ID {
	case "TotalMemory":
		metric := models.NewMetricGauge(ID, float64(v.Total))
		metrics = []models.Metrics{metric}
	case "FreeMemory":
		metric := models.NewMetricGauge(ID, float64(v.Free))
		metrics = []models.Metrics{metric}
	case "CPUutilization":
		cpudata, err := cpu.Percent(time.Second, true)
		if err == nil {
			for k, v := range cpudata {
				metric := models.NewMetricGauge(ID+strconv.Itoa(k+1), float64(v))
				metrics = append(metrics, metric)
			}
		}
	default:
		err = errors.New("invalid metric type")
	}

	return metrics, err
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
