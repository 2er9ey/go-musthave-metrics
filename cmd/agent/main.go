package main

import (
	"bytes"
	"compress/gzip"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/2er9ey/go-musthave-metrics/internal/agent"
	"github.com/2er9ey/go-musthave-metrics/internal/models"
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
	cm2 := agent.NewPSCollectionMetrics()
	repo, _ := repository.NewMemoryStorage()

	var wg sync.WaitGroup
	go getMetrics(repo, cm)
	go getMetrics(repo, cm2)
	time.Sleep(config.reportInterval)
	if config.rateLimit > 0 {
		senderMetricsCompressedSeparately(repo, config.signingKey, config.rateLimit)
	} else {
		senderMetricsCompressed(repo, config.signingKey)
	}
	wg.Wait()
	// fmt.Println("All workers are done!")
}

func getMetrics(repo repository.MetricsRepositoryInterface, collectMetrics *[]agent.CollectMetric) {
	for {
		//		mutex.Lock()
		agent.CollectorMetrics(repo, collectMetrics)
		//		mutex.Unlock()
		//		fmt.Println("Метрики собраны")
		time.Sleep(config.pollInterval)
	}
}

func senderMetricsCompressed(repo repository.MetricsRepositoryInterface, key string) {
	for {
		metrics := repo.GetAllMetric()
		buf := bytes.NewBuffer(nil)
		PrepareBuf(buf, metrics)
		sendBunchMetricsWithRetry(4, buf, key)
		time.Sleep(config.reportInterval)
	}
}

func senderMetricsCompressedSeparately(repo repository.MetricsRepositoryInterface, key string, rateLimit int) {
	var bufferPool = sync.Pool{
		// Функция New сработает, если в пуле нет объекта
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}
	var metricPool = sync.Pool{
		// Функция New сработает, если в пуле нет объекта
		New: func() interface{} {
			return new([]models.Metrics)
		},
	}

	ch1 := make(chan models.Metrics)
	defer close(ch1)

	for i := range rateLimit {
		fmt.Println("Running sending worker N", i)
		go func() {
			for {
				metric := <-ch1
				sendMetricsPtr := metricPool.Get().(*[]models.Metrics)
				sendMetrics := append(*sendMetricsPtr, metric)
				buf := bufferPool.Get().(*bytes.Buffer)
				buf.Reset()
				PrepareBuf(buf, sendMetrics)
				sendBunchMetricsWithRetry(4, buf, key)
				bufferPool.Put(buf)
				bufferPool.Put(sendMetricsPtr)
			}
		}()
	}

	for {
		metrics := repo.GetAllMetric()
		for _, v := range metrics {
			ch1 <- v
		}
		time.Sleep(config.reportInterval)
	}
}

func sendBunchMetricsWithRetry(maxReties int, buf *bytes.Buffer, key string) {
	retryTimeout := 1
	request, err := http.NewRequest("POST", "http://"+config.serverEndpoint+"/updates", buf)
	if err != nil {
		return
	}
	request.Header.Set("Content-Encoding", "gzip")
	request.Header.Set("Content-type", "application/json")
	if key != "" {
		h := hmac.New(sha256.New, []byte(key))
		h.Write(buf.Bytes())
		dst := hex.EncodeToString(h.Sum(nil))
		request.Header.Set("HashSHA256", dst)
	}
	for retry := range maxReties {
		fmt.Println("Try ", retry)
		resp, err2 := http.DefaultClient.Do(request)
		if err2 == nil {
			resp.Body.Close()
			break
		}
		switch err2.(type) {
		case *net.OpError:
		default:
			break
		}
		if retry < 3 {
			fmt.Printf("Sleeping %d seconds\n", retryTimeout)
			time.Sleep(time.Duration(retryTimeout) * time.Second)
			fmt.Println("Wakeup")
			retryTimeout = retryTimeout + 2
		}
	}
}

func PrepareBuf(buf *bytes.Buffer, metrics []models.Metrics) *bytes.Buffer {
	jsonValue, _ := json.Marshal(metrics)
	zb := gzip.NewWriter(buf)
	zb.Write(jsonValue)
	zb.Close()
	return buf
}
