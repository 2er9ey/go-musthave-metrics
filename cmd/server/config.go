package main

import (
	"flag"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	listenEndpoint string
	logLevel       string
}

func parseConfig() (Config, error) {
	var cfgEnv struct {
		ListenEndpoint string `env:"ADDRESS"`
	}
	var conf Config

	flag.StringVar(&conf.listenEndpoint, "a", "localhost:8080", "Адрес и порт для работы севрера")
	flag.StringVar(&conf.logLevel, "l", "info", "Уровень журналирования")
	flag.Parse()

	if err := env.Parse(&cfgEnv); err != nil {
		return conf, err
	}

	if cfgEnv.ListenEndpoint != "" {
		conf.listenEndpoint = cfgEnv.ListenEndpoint
	}

	return conf, nil
}
