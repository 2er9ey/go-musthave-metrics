package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
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
	sendBunchMetricsCompressed(repo)
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
			//			metricValue := handler.MetricRequest{ID: value.ID, MType: value.MType, Value: value.String()}
			jsonValue, _ := json.Marshal(value)
			fmt.Println(">", jsonValue)
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
			//			metricValue := handler.MetricRequest{ID: value.ID, MType: value.MType, Value: value.String()}
			jsonValue, _ := json.Marshal(value)
			//fmt.Println(">", string(jsonValue))
			buf := bytes.NewBuffer(nil)
			zb := gzip.NewWriter(buf)
			zb.Write(jsonValue)
			zb.Close()
			request, err := http.NewRequest("POST", "http://"+config.serverEndpoint+"/update", buf)
			if err != nil {
				break
			}
			request.Header.Set("Content-Encoding", "gzip")
			request.Header.Set("Content-type", "application/json")
			resp, err2 := http.DefaultClient.Do(request)
			if err2 != nil {
				break
			}
			resp.Body.Close()
		}
		//		fmt.Println("Метрики отправлены")
		time.Sleep(config.reportInterval)
	}
}

func sendBunchMetricsCompressed(repo repository.MetricsRepositoryInterface) {
	for {
		mutex.Lock()
		metrics := repo.GetAllMetric()
		mutex.Unlock()
		jsonValue, _ := json.Marshal(metrics)
		buf := bytes.NewBuffer(nil)
		zb := gzip.NewWriter(buf)
		zb.Write(jsonValue)
		zb.Close()
		retryTimeout := 1
		for retry := 0; retry < 4; retry++ {
			fmt.Println("Try ", retry)
			request, err := http.NewRequest("POST", "http://"+config.serverEndpoint+"/updates", buf)
			if err != nil {
				break
			}
			request.Header.Set("Content-Encoding", "gzip")
			request.Header.Set("Content-type", "application/json")
			resp, err2 := http.DefaultClient.Do(request)
			if err2 != nil {
				switch err2 := err2.(type) {
				case *url.Error:
					if err2.Timeout() {
						fmt.Printf("timeout: %s", err2.Err)
					} else if err2, ok := err2.Err.(*net.OpError); ok {
						fmt.Printf("net error: %s\n", err2)
					} else {
						fmt.Printf("original error: %T\n", err2)
					}
				default:
					fmt.Printf("unknown error: %v\n", err2)
				}
			} else {
				resp.Body.Close()
				break
			}
			if retry < 3 {
				fmt.Printf("Sleeping %d seconds\n", retryTimeout)
				time.Sleep(time.Duration(retryTimeout) * time.Second)
				fmt.Println("Wakeup")
				retryTimeout = retryTimeout + 2
			}
		}
		//		fmt.Println("Метрики отправлены")
		time.Sleep(config.reportInterval)
	}
}
