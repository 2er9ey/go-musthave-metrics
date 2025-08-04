package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/2er9ey/go-musthave-metrics/internal/agent"
	"github.com/2er9ey/go-musthave-metrics/internal/repository"
)

var mutex sync.Mutex
var pollInterval = 2
var reportInterval = 10

func main() {
	cm := agent.NewCollectionMectics()
	var repo repository.MetricsRepositoryInterface = repository.NewMemoryStorage()

	var wg sync.WaitGroup
	go getMetrics(repo, cm)
	sendMetrics(repo)
	wg.Wait()
	fmt.Println("All workers are done!")
}

func getMetrics(repo repository.MetricsRepositoryInterface, collectMetrics *[]agent.CollectMetric) {
	for {
		mutex.Lock()
		agent.CollectorMetrics(repo, collectMetrics)
		mutex.Unlock()
		fmt.Println("Метрики собраны")
		time.Sleep(time.Duration(pollInterval) * time.Second)
	}
}

func sendMetrics(repo repository.MetricsRepositoryInterface) {
	for {
		mutex.Lock()
		metrics := repo.GetAllMetric()
		mutex.Unlock()
		for key, value := range metrics {
			response, err := http.Post("http://127.0.0.1:8080/update/"+value.MType+"/"+key+"/"+value.String(), "text/plain", nil)
			if err != nil {
				fmt.Println("Ошибка отправки метрик")
				break
			}
			response.Body.Close()
		}
		fmt.Println("Метрики отправлены")
		time.Sleep(time.Duration(reportInterval) * time.Second)
	}
}
