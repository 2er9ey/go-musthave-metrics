package main

import (
	"net/http"

	"github.com/2er9ey/go-musthave-metrics/internal/handler"
	"github.com/2er9ey/go-musthave-metrics/internal/repository"
	"github.com/2er9ey/go-musthave-metrics/internal/service"
)

func main() {
	repo := repository.NewMemoryStorage()
	service := service.NewMetricService(repo)
	metricsHadler := handler.NewMetricHandler(service)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /update/{metricType}/{metricName}/{metricValue}", metricsHadler.PostUpdate)
	mux.HandleFunc("POST /update/{metricType}/", metricsHadler.StatusNotFound)
	mux.HandleFunc("POST /update/{metricType}", metricsHadler.StatusNotFound)
	mux.HandleFunc("POST /update/", metricsHadler.StatusBadRequest)
	mux.HandleFunc("POST /update", metricsHadler.StatusBadRequest)

	err := http.ListenAndServe(`localhost:8080`, mux)
	if err != nil {
		panic(err)
	}
}
