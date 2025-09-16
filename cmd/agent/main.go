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
	repo, _ := repository.NewMemoryStorage()

	var wg sync.WaitGroup
	go getMetrics(repo, cm)
	time.Sleep(config.reportInterval)
	senderMetrics(repo)
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

func senderMetrics(repo repository.MetricsRepositoryInterface) {
	for {
		buf := GetMetricsBunch(repo)
		sendBunchMetricsCompressedWithRetry(repo, 4, buf)
	}
}

func sendBunchMetricsCompressedWithRetry(repo repository.MetricsRepositoryInterface, maxReties int, buf *bytes.Buffer) {
	retryTimeout := 1
	request, err := http.NewRequest("POST", "http://"+config.serverEndpoint+"/updates", buf)
	if err != nil {
		return
	}
	request.Header.Set("Content-Encoding", "gzip")
	request.Header.Set("Content-type", "application/json")
	for retry := 0; retry < 4; retry++ {
		fmt.Println("Try ", retry)
		resp, err2 := http.DefaultClient.Do(request)
		if err2 != nil {
			switch err2 := err2.(type) {
			case *url.Error:
				_, ok := err2.Err.(*net.OpError)
				if !err2.Timeout() && !ok {
					break
				}
			default:
				break
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
	time.Sleep(config.reportInterval)
}

func GetMetricsBunch(repo repository.MetricsRepositoryInterface) *bytes.Buffer {
	mutex.Lock()
	metrics := repo.GetAllMetric()
	mutex.Unlock()
	jsonValue, _ := json.Marshal(metrics)
	buf := bytes.NewBuffer(nil)
	zb := gzip.NewWriter(buf)
	zb.Write(jsonValue)
	zb.Close()
	return buf
}
