package main

import (
	"fmt"

	"github.com/2er9ey/go-musthave-metrics/internal/handler"
	"github.com/2er9ey/go-musthave-metrics/internal/logger"
	"github.com/2er9ey/go-musthave-metrics/internal/repository"
	"github.com/2er9ey/go-musthave-metrics/internal/service"
	"go.uber.org/zap"
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

	if err := logger.Initialize(config.logLevel); err != nil {
		fmt.Println("Ошибка журнала", err)
		return
	}

	router := SetupRouter(*metricsHandler)
	logger.Log.Info("Startin server listen on", zap.String("listenEndpoint", config.listenEndpoint))
	router.Run(config.listenEndpoint)
}
