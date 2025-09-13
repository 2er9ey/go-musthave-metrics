package main

import (
	"context"
	"fmt"
	"os"
	"reflect"

	"github.com/2er9ey/go-musthave-metrics/internal/handler"
	"github.com/2er9ey/go-musthave-metrics/internal/logger"
	"github.com/2er9ey/go-musthave-metrics/internal/repository"
	"github.com/2er9ey/go-musthave-metrics/internal/service"
	"go.uber.org/zap"
)

func main() {
	config, configError := parseConfig()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if configError != nil {
		fmt.Println("Ошибка чтения конфигурации", configError)
		return
	}

	if err := logger.Initialize(config.logLevel); err != nil {
		fmt.Println("Ошибка журнала", err)
		return
	}
	defer logger.Log.Sync()

	logger.Log.Info("Config: ", zap.String("databaseDSN", config.databaseDSN),
		zap.Int("StoreInterval", config.storeInterval),
		zap.String("FileStoragePath", config.fileStoragePath),
		zap.Bool("RestoreMetrics", config.restoreMetrics))

	var repo repository.MetricsRepositoryInterface
	if config.databaseDSN == "" {
		if config.fileStoragePath == "" {
			repo = repository.NewMemoryStorage()
		} else {
			repo = repository.NewFileStorage(config.fileStoragePath, config.storeInterval, config.restoreMetrics)
		}
	} else {
		ps := repository.NewPostgreSQLStorage(ctx, config.databaseDSN)
		if ps != nil {
			ps.CreateTables()
			defer ps.Close()
		}
		repo = ps
	}
	v := reflect.ValueOf(repo)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		logger.Log.Error("Ошибка создания репозитория метрик")
		os.Exit(1)
	}
	service := service.NewMetricService(ctx, repo, config.storeInterval, config.fileStoragePath, config.databaseDSN)
	metricsHandler := handler.NewMetricHandler(service)

	router := SetupRouter(*metricsHandler)
	logger.Log.Info("Starting server listen on", zap.String("listenEndpoint", config.listenEndpoint))
	router.Run(config.listenEndpoint)
}
