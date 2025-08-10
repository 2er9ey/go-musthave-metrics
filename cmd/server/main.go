package main

import (
	"flag"

	"github.com/2er9ey/go-musthave-metrics/internal/handler"
	"github.com/2er9ey/go-musthave-metrics/internal/repository"
	"github.com/2er9ey/go-musthave-metrics/internal/server"
	"github.com/2er9ey/go-musthave-metrics/internal/service"
)

func main() {
	repo := repository.NewMemoryStorage()
	service := service.NewMetricService(repo)
	metricsHandler := handler.NewMetricHandler(service)
	listenEndpoint := flag.String("a", "localhost:8080", "Адрес и порт для работы севрера")

	flag.Parse()

	router := server.SetupRouter(*metricsHandler)

	router.Run(*listenEndpoint)
}
