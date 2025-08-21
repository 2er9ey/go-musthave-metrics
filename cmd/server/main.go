package main

import (
	"fmt"

	"github.com/2er9ey/go-musthave-metrics/internal/handler"
	"github.com/2er9ey/go-musthave-metrics/internal/repository"
	"github.com/2er9ey/go-musthave-metrics/internal/service"
)

func main() {
	repo := repository.NewMemoryStorage()
	service := service.NewMetricService(repo)
	metricsHandler := handler.NewMetricHandler(service)

	config, configError := parseConfig()

	if configError != nil {
		fmt.Println("Ошибка чтения конфигурации", configError)
		return
	}

	router := SetupRouter(*metricsHandler)
	router.Run(config.listenEndpoint)
}
