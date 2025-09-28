package agent

import "github.com/2er9ey/go-musthave-metrics/internal/models"

const (
	collectTypePsUtils = "PSUtils"
	collectTypeMemStat = "MemStat"
	collectTypeRandom  = "Random"
	collectTypeConst1  = "CollectConst1"
)

type CollectMetric struct {
	ID          string
	MType       string
	collectType string
}

func NewCollectionMetrics() *[]CollectMetric {
	return &([]CollectMetric{
		{ID: "Alloc", collectType: collectTypeMemStat, MType: models.Gauge},
		{ID: "BuckHashSys", collectType: collectTypeMemStat, MType: models.Gauge},
		{ID: "Frees", collectType: collectTypeMemStat, MType: models.Gauge},
		{ID: "GCCPUFraction", collectType: collectTypeMemStat, MType: models.Gauge},
		{ID: "GCSys", collectType: collectTypeMemStat, MType: models.Gauge},
		{ID: "HeapAlloc", collectType: collectTypeMemStat, MType: models.Gauge},
		{ID: "HeapIdle", collectType: collectTypeMemStat, MType: models.Gauge},
		{ID: "HeapInuse", collectType: collectTypeMemStat, MType: models.Gauge},
		{ID: "HeapObjects", collectType: collectTypeMemStat, MType: models.Gauge},
		{ID: "HeapReleased", collectType: collectTypeMemStat, MType: models.Gauge},
		{ID: "HeapSys", collectType: collectTypeMemStat, MType: models.Gauge},
		{ID: "LastGC", collectType: collectTypeMemStat, MType: models.Gauge},
		{ID: "Lookups", collectType: collectTypeMemStat, MType: models.Gauge},
		{ID: "MCacheInuse", collectType: collectTypeMemStat, MType: models.Gauge},
		{ID: "MCacheSys", collectType: collectTypeMemStat, MType: models.Gauge},
		{ID: "MSpanInuse", collectType: collectTypeMemStat, MType: models.Gauge},
		{ID: "MSpanSys", collectType: collectTypeMemStat, MType: models.Gauge},
		{ID: "Mallocs", collectType: collectTypeMemStat, MType: models.Gauge},
		{ID: "NextGC", collectType: collectTypeMemStat, MType: models.Gauge},
		{ID: "NumForcedGC", collectType: collectTypeMemStat, MType: models.Gauge},
		{ID: "NumGC", collectType: collectTypeMemStat, MType: models.Gauge},
		{ID: "OtherSys", collectType: collectTypeMemStat, MType: models.Gauge},
		{ID: "PauseTotalNs", collectType: collectTypeMemStat, MType: models.Gauge},
		{ID: "StackInuse", collectType: collectTypeMemStat, MType: models.Gauge},
		{ID: "StackSys", collectType: collectTypeMemStat, MType: models.Gauge},
		{ID: "Sys", collectType: collectTypeMemStat, MType: models.Gauge},
		{ID: "TotalAlloc", collectType: collectTypeMemStat, MType: models.Gauge},
		{ID: "RandomValue", collectType: collectTypeRandom, MType: models.Gauge},
		{ID: "PollCount", collectType: collectTypeConst1, MType: models.Counter},
	})
}

func NewPSCollectionMetrics() *[]CollectMetric {
	return &([]CollectMetric{
		{ID: "TotalMemory", collectType: collectTypePsUtils, MType: models.Gauge},
		{ID: "FreeMemory", collectType: collectTypePsUtils, MType: models.Gauge},
		{ID: "CPUutilization", collectType: collectTypePsUtils, MType: models.Gauge},
	})
}
