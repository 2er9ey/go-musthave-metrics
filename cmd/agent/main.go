package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/2er9ey/go-musthave-metrics/internal/agent"
	"github.com/2er9ey/go-musthave-metrics/internal/repository"
)

var mutex sync.Mutex
var config Config

func main() {
	var configError error

	config, configError = parseConfig()

	if configError != nil {
		fmt.Println("Ошибка чтения конфигурации", configError)
		return
	}

	cm := agent.NewCollectionMetrics()
	var repo repository.MetricsRepositoryInterface = repository.NewMemoryStorage()

	var wg sync.WaitGroup
	go getMetrics(repo, cm)
	time.Sleep(config.reportInterval)
	sendMetricsCompressed(repo)
	wg.Wait()
	// fmt.Println("All workers are done!")
}

func getMetrics(repo repository.MetricsRepositoryInterface, collectMetrics *[]agent.CollectMetric) {
	for {
		mutex.Lock()
		agent.CollectorMetrics(repo, collectMetrics)
		mutex.Unlock()
		//		fmt.Println("Метрики собраны")
		time.Sleep(config.pollInterval)
	}
}

func sendMetrics(repo repository.MetricsRepositoryInterface) {
	for {
		mutex.Lock()
		metrics := repo.GetAllMetric()
		mutex.Unlock()
		for _, value := range metrics {
			jsonValue, _ := json.Marshal(value)
			response, err := http.Post("http://"+config.serverEndpoint+"/update", "application/json", bytes.NewBuffer(jsonValue))
			if err != nil {
				//				fmt.Println("Ошибка отправки метрик")
				break
			}
			response.Body.Close()
		}
		//		fmt.Println("Метрики отправлены")
		time.Sleep(config.reportInterval)
	}
}

func sendMetricsCompressed(repo repository.MetricsRepositoryInterface) {
	for {
		mutex.Lock()
		metrics := repo.GetAllMetric()
		mutex.Unlock()
		for _, value := range metrics {
			jsonValue, _ := json.Marshal(value)
			buf := bytes.NewBuffer(nil)
			zb := gzip.NewWriter(buf)
			zb.Write(jsonValue)
			zb.Close()
			request, err := http.NewRequest("POST", "http://"+config.serverEndpoint+"/update", buf)
			if err != nil {
				//				fmt.Println("Ошибка отправки метрик")
				break
			}
			request.Header.Set("Content-Encoding", "gzip")
			request.Header.Set("Content-type", "application/json")
			resp, _ := http.DefaultClient.Do(request)
			resp.Body.Close()
		}
		//		fmt.Println("Метрики отправлены")
		time.Sleep(config.reportInterval)
	}
}
