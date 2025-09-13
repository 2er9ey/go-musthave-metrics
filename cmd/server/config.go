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
	databaseDSN     string
}

func parseConfig() (Config, error) {
	var cfgEnv struct {
		ListenEndpoint  string `env:"ADDRESS"`
		StoreInterval   string `env:"STORE_INTERVAL"`
		FileStoragePath string `env:"FILE_STORAGE_PATH"`
		RestoreMetrics  string `env:"RESTORE"`
		DatabaseDSN     string `env:"DATABASE_DSN"`
	}
	var conf Config

	flag.StringVar(&conf.listenEndpoint, "a", "localhost:8080", "Адрес и порт для работы севрера")
	flag.StringVar(&conf.logLevel, "l", "debug", "Уровень журналирования")
	flag.IntVar(&conf.storeInterval, "i", 300, "Интервал сохранения значений метрик")
	flag.StringVar(&conf.fileStoragePath, "f", "metrics.dat", "Имя файла для сохранения значения метрик")
	//	flag.StringVar(&conf.databaseDSN, "d", "host=127.0.0.1 user=video password=XXXXXXXX dbname=video sslmode=disable", "Строка подключения к базе данных")
	flag.StringVar(&conf.databaseDSN, "d", "", "Строка подключения к базе данных")
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

	if cfgEnv.DatabaseDSN != "" {
		conf.databaseDSN = cfgEnv.DatabaseDSN
	}

	if cfgEnv.StoreInterval != "" {
		conf.storeInterval, _ = strconv.Atoi(cfgEnv.StoreInterval)
	}

	if cfgEnv.RestoreMetrics != "" {
		conf.restoreMetrics, _ = strconv.ParseBool(cfgEnv.RestoreMetrics)
	}

	return conf, nil
}
