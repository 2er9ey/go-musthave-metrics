package main

import (
	"flag"
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	serverEndpoint string
	signingKey     string
	pollInterval   time.Duration
	reportInterval time.Duration
}

func parseConfig() (Config, error) {
	var cfgEnv struct {
		ServerEndpoint string `env:"ADDRESS"`
		PollInterval   string `env:"POLL_INTERVAL"`
		ReportInterval string `env:"REPORT_INTERVAL"`
		SigningKey     string `env:"KEY"`
	}
	var cfgFlag struct {
		ServerEndpoint string
		PollInterval   string
		ReportInterval string
		SigningKey     string
	}
	var conf Config

	flag.StringVar(&cfgFlag.ServerEndpoint, "a", "localhost:8080", "Адрес и порт для работы севрера")
	flag.StringVar(&cfgFlag.PollInterval, "p", "2s", "Время опроса метрик")
	flag.StringVar(&cfgFlag.ReportInterval, "r", "10s", "Время отправки метрик на сервер")
	flag.StringVar(&cfgFlag.SigningKey, "k", "", "Ключ для подписи")
	flag.Parse()

	if err := env.Parse(&cfgEnv); err != nil {
		return conf, err
	}

	if cfgEnv.ServerEndpoint != "" {
		conf.serverEndpoint = cfgEnv.ServerEndpoint
	} else {
		conf.serverEndpoint = cfgFlag.ServerEndpoint
	}

	if cfgEnv.SigningKey != "" {
		conf.signingKey = cfgEnv.SigningKey
	} else {
		conf.signingKey = cfgFlag.SigningKey
	}

	var tmpPollInterval string
	if cfgEnv.PollInterval != "" {
		tmpPollInterval = cfgEnv.PollInterval
	} else {
		tmpPollInterval = cfgFlag.PollInterval
	}

	if res, err := time.ParseDuration(tmpPollInterval); err != nil {
		res, err = time.ParseDuration(tmpPollInterval + "s")
		if err != nil {
			return conf, err
		}
		conf.pollInterval = res
	} else {
		conf.pollInterval = res
	}

	var tmpReportInterval string
	if cfgEnv.ReportInterval != "" {
		tmpReportInterval = cfgEnv.ReportInterval
	} else {
		tmpReportInterval = cfgFlag.ReportInterval
	}

	if res, err := time.ParseDuration(tmpReportInterval); err != nil {
		res, err = time.ParseDuration(tmpReportInterval + "s")
		if err != nil {
			return conf, err
		}
		conf.reportInterval = res
	} else {
		conf.reportInterval = res
	}

	return conf, nil
}
