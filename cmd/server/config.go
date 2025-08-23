package main

import (
	"flag"
	"strconv"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	listenEndpoint  string
	logLevel        string
	storeInterval   int
	fileStoragePath string
	restoreMetrics  bool
}

func parseConfig() (Config, error) {
	var cfgEnv struct {
		ListenEndpoint  string `env:"ADDRESS"`
		StoreInterval   string `env:"STORE_INTERVAL"`
		FileStoragePath string `env:"FILE_STORAGE_PATH"`
		RestoreMetrics  string `env:"RESTORE"`
	}
	var conf Config

	flag.StringVar(&conf.listenEndpoint, "a", "localhost:8080", "Адрес и порт для работы севрера")
	flag.StringVar(&conf.logLevel, "l", "info", "Уровень журналирования")
	flag.IntVar(&conf.storeInterval, "i", 300, "Интервал сохранения значений метрик")
	flag.StringVar(&conf.fileStoragePath, "f", "metrics.dat", "Имя файла для сохранения значения метрик")
	flag.BoolVar(&conf.restoreMetrics, "r", false, "Считать значения метрик при старте сервера")
	flag.Parse()

	if err := env.Parse(&cfgEnv); err != nil {
		return conf, err
	}

	if cfgEnv.ListenEndpoint != "" {
		conf.listenEndpoint = cfgEnv.ListenEndpoint
	}

	if cfgEnv.FileStoragePath != "" {
		conf.fileStoragePath = cfgEnv.FileStoragePath
	}

	if cfgEnv.StoreInterval != "" {
		conf.storeInterval, _ = strconv.Atoi(cfgEnv.StoreInterval)
	}

	if cfgEnv.RestoreMetrics != "" {
		conf.restoreMetrics, _ = strconv.ParseBool(cfgEnv.RestoreMetrics)
	}

	return conf, nil
}
