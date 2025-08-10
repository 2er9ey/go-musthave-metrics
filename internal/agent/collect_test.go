package agent

import (
	"runtime"
	"testing"

	"github.com/2er9ey/go-musthave-metrics/internal/models"
)

func TestMemStatWithCorrectMectricGauge(t *testing.T) {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	metric, err := collectGetMemStatMectric(ms, "Alloc", models.Gauge)
	if err != nil {
		t.Fatal("Ошибки быть не должно!")
	}
	if metric.ID != "Alloc" || metric.MType != models.Gauge ||
		metric.Value == nil || metric.Delta != nil {
		t.Fatal("Метрика заполнена неверно")
	}
}

func TestMemStatWithCorrectMectricCounter(t *testing.T) {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	metric, err := collectGetMemStatMectric(ms, "Alloc", models.Counter)
	if err != nil {
		t.Fatal("Ошибки быть не должно!")
	}
	if metric.ID != "Alloc" || metric.MType != models.Counter ||
		metric.Value != nil || metric.Delta == nil {
		t.Fatal("Метрика заполнена неверно")
	}
}

func TestMemStatWithIncorrectMectric(t *testing.T) {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	_, err := collectGetMemStatMectric(ms, "XXX", models.Gauge)
	if err == nil {
		t.Fatal("Должна быть ошибка")
	}
}

func TestRandomMectricGauge(t *testing.T) {
	metric, err := collectGetRandomMectric("XXX", models.Gauge)
	if err != nil {
		t.Fatal("Ошибки быть не должно")
	}
	if metric.ID != "XXX" || metric.MType != models.Gauge ||
		metric.Value == nil || metric.Delta != nil {
		t.Fatal("Метрика заполнена неверно")
	}
}

func TestRandomMectricCounter(t *testing.T) {
	metric, err := collectGetRandomMectric("XXX", models.Counter)
	if err != nil {
		t.Fatal("Ошибки быть не должно")
	}
	if metric.ID != "XXX" || metric.MType != models.Counter ||
		metric.Value != nil || metric.Delta == nil {
		t.Fatal("Метрика заполнена неверно")
	}
}
