package main

import (
	"flag"
	"net/http"
	"sync"
	"time"

	"github.com/2er9ey/go-musthave-metrics/internal/agent"
	"github.com/2er9ey/go-musthave-metrics/internal/repository"
)

var mutex sync.Mutex
var pollInterval = 2
var reportInterval = 10
var listenEndpoint = "localhost:8080"

func main() {
	flag.StringVar(&listenEndpoint, "a", "localhost:8080", "Адрес и порт для работы севрера")
	flag.IntVar(&pollInterval, "p", 2, "Время опроса метрик")
	flag.IntVar(&reportInterval, "r", 10, "Время отправки метрик на сервер")
	flag.Parse()

	cm := agent.NewCollectionMetrics()
	var repo repository.MetricsRepositoryInterface = repository.NewMemoryStorage()

	var wg sync.WaitGroup
	go getMetrics(repo, cm)
	time.Sleep(time.Duration(reportInterval) * time.Second)
	sendMetrics(repo)
	wg.Wait()
	// fmt.Println("All workers are done!")
}

func getMetrics(repo repository.MetricsRepositoryInterface, collectMetrics *[]agent.CollectMetric) {
	for {
		mutex.Lock()
		agent.CollectorMetrics(repo, collectMetrics)
		mutex.Unlock()
		//		fmt.Println("Метрики собраны")
		time.Sleep(time.Duration(pollInterval) * time.Second)
	}
}

func sendMetrics(repo repository.MetricsRepositoryInterface) {
	for {
		mutex.Lock()
		metrics := repo.GetAllMetric()
		mutex.Unlock()
		for _, value := range metrics {
			response, err := http.Post("http://"+listenEndpoint+"/update/"+value.MType+"/"+value.ID+"/"+value.String(), "text/plain", nil)
			if err != nil {
				//				fmt.Println("Ошибка отправки метрик")
				break
			}
			response.Body.Close()
		}
		//		fmt.Println("Метрики отправлены")
		time.Sleep(time.Duration(reportInterval) * time.Second)
	}
}
