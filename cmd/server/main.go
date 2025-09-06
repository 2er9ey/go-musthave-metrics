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
	config, configError := parseConfig()

	if configError != nil {
		fmt.Println("Ошибка чтения конфигурации", configError)
		return
	}

	if err := logger.Initialize(config.logLevel); err != nil {
		fmt.Println("Ошибка журнала", err)
		return
	}

	repo := repository.NewMemoryStorage()
	service := service.NewMetricService(repo, config.storeInterval, config.fileStoragePath)
	if config.restoreMetrics {
		service.LoadMetrics(config.fileStoragePath)
	}
	if config.storeInterval > 0 {
		service.RunSaver()
	}
	metricsHandler := handler.NewMetricHandler(service)

	router := SetupRouter(*metricsHandler)
	logger.Log.Info("Startin server listen on", zap.String("listenEndpoint", config.listenEndpoint))
	router.Run(config.listenEndpoint)
}
