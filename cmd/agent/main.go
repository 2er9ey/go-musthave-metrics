package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"strconv"
	"sync"
	"time"
)

type HostMetric struct {
	Alloc         uint64
	BuckHashSys   uint64
	Frees         uint64
	GCCPUFraction float64
	GCSys         uint64
	HeapAlloc     uint64
	HeapIdle      uint64
	HeapInuse     uint64
	HeapObjects   uint64
	HeapReleased  uint64
	HeapSys       uint64
	LastGC        uint64
	Lookups       uint64
	MCacheInuse   uint64
	MCacheSys     uint64
	MSpanInuse    uint64
	MSpanSys      uint64
	Mallocs       uint64
	NextGC        uint64
	NumForcedGC   uint32
	NumGC         uint32
	OtherSys      uint64
	PauseTotalNs  uint64
	StackInuse    uint64
	StackSys      uint64
	Sys           uint64
	TotalAlloc    uint64
	PollCount     uint64
	RandomValue   float64
}

var metricsValues HostMetric
var mutex sync.Mutex
var pollInterval = 2
var reportInterval = 10

func main() {
	var wg sync.WaitGroup
	// getMetrics()
	// sendGaugeMetric("RandomValue", 0.123456)
	// sendCounterMetric("RandomValue", 0)
	go getMetrics()
	time.Sleep(5 * time.Second)
	sendMetrics()
	wg.Wait()
	fmt.Println("All workers are done!")
}

func getMetrics() {
	for {
		mutex.Lock()
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		metricsValues.Alloc = m.Alloc
		metricsValues.BuckHashSys = m.BuckHashSys
		metricsValues.Frees = m.Frees
		metricsValues.GCCPUFraction = m.GCCPUFraction
		metricsValues.GCSys = m.GCSys
		metricsValues.HeapAlloc = m.HeapAlloc
		metricsValues.HeapIdle = m.HeapIdle
		metricsValues.HeapInuse = m.HeapInuse
		metricsValues.HeapObjects = m.HeapObjects
		metricsValues.HeapReleased = m.HeapReleased
		metricsValues.HeapSys = m.HeapSys
		metricsValues.LastGC = m.LastGC
		metricsValues.Lookups = m.Lookups
		metricsValues.MCacheInuse = m.MCacheInuse
		metricsValues.MCacheSys = m.MCacheSys
		metricsValues.MSpanInuse = m.MSpanInuse
		metricsValues.MSpanSys = m.MSpanSys
		metricsValues.Mallocs = m.Mallocs
		metricsValues.NextGC = m.NextGC
		metricsValues.NumForcedGC = m.NumForcedGC
		metricsValues.NumGC = m.NumGC
		metricsValues.OtherSys = m.OtherSys
		metricsValues.PauseTotalNs = m.PauseTotalNs
		metricsValues.StackInuse = m.StackInuse
		metricsValues.StackSys = m.StackSys
		metricsValues.Sys = m.Sys
		metricsValues.TotalAlloc = m.TotalAlloc
		metricsValues.PollCount += 1
		metricsValues.RandomValue = rand.Float64()
		mutex.Unlock()
		//		fmt.Println(metricsValues)
		time.Sleep(time.Duration(pollInterval) * time.Second)
	}
}

func sendMetrics() {
	for {
		mutex.Lock()
		sendGaugeMetric("Alloc", float64(metricsValues.Alloc))
		sendGaugeMetric("BuckHashSys", float64(metricsValues.BuckHashSys))
		sendGaugeMetric("Frees", float64(metricsValues.Frees))
		sendGaugeMetric("GCCPUFraction", float64(metricsValues.GCCPUFraction))
		sendGaugeMetric("GCSys", float64(metricsValues.GCSys))
		sendGaugeMetric("HeapAlloc", float64(metricsValues.HeapAlloc))
		sendGaugeMetric("HeapIdle", float64(metricsValues.HeapIdle))
		sendGaugeMetric("HeapInuse", float64(metricsValues.HeapInuse))
		sendGaugeMetric("HeapObjects", float64(metricsValues.HeapObjects))
		sendGaugeMetric("HeapReleased", float64(metricsValues.HeapReleased))
		sendGaugeMetric("HeapSys", float64(metricsValues.HeapSys))
		sendGaugeMetric("LastGC", float64(metricsValues.LastGC))
		sendGaugeMetric("Lookups", float64(metricsValues.Lookups))
		sendGaugeMetric("MCacheInuse", float64(metricsValues.MCacheInuse))
		sendGaugeMetric("MCacheSys", float64(metricsValues.MCacheSys))
		sendGaugeMetric("MSpanInuse", float64(metricsValues.MSpanInuse))
		sendGaugeMetric("MSpanSys", float64(metricsValues.MSpanSys))
		sendGaugeMetric("Mallocs", float64(metricsValues.Mallocs))
		sendGaugeMetric("NextGC", float64(metricsValues.NextGC))
		sendGaugeMetric("NumForcedGC", float64(metricsValues.NumForcedGC))
		sendGaugeMetric("NumGC", float64(metricsValues.NumGC))
		sendGaugeMetric("OtherSys", float64(metricsValues.OtherSys))
		sendGaugeMetric("PauseTotalNs", float64(metricsValues.PauseTotalNs))
		sendGaugeMetric("StackInuse", float64(metricsValues.StackInuse))
		sendGaugeMetric("StackSys", float64(metricsValues.StackSys))
		sendGaugeMetric("Sys", float64(metricsValues.Sys))
		sendGaugeMetric("TotalAlloc", float64(metricsValues.TotalAlloc))
		sendGaugeMetric("RandomValue", float64(metricsValues.RandomValue))
		sendCounterMetric("PollCount", int64(metricsValues.PollCount))
		mutex.Unlock()
		//		fmt.Println("Sending metrics")
		time.Sleep(time.Duration(reportInterval) * time.Second)
	}

}

func sendGaugeMetric(name string, value float64) {
	response, err := http.Post("http://127.0.0.1:8080/update/gauge/"+name+"/"+strconv.FormatFloat(value, 'E', 5, 64), "text/plain", nil)

	if err != nil {
		return
	}

	defer response.Body.Close()

	// fmt.Printf("Status Code: %d\r\n", response.StatusCode)
	// for k, v := range response.Header {
	// 	// заголовок может иметь несколько значений,
	// 	// но для простоты запросим только первое
	// 	fmt.Printf("%s: %v\r\n", k, v[0])
	// }

}

func sendCounterMetric(name string, value int64) {
	response, err := http.Post("http://127.0.0.1:8080/update/gauge/"+name+"/"+strconv.FormatInt(value, 10), "text/plain", nil)

	if err != nil {
		return
	}

	defer response.Body.Close()

	// fmt.Printf("Status Code: %d\r\n", response.StatusCode)
	// for k, v := range response.Header {
	// 	// заголовок может иметь несколько значений,
	// 	// но для простоты запросим только первое
	// 	fmt.Printf("%s: %v\r\n", k, v[0])
	// }

}
