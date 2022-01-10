package main

import (
	"flag"

	"github.com/deesel/wol/internal/api"
	"github.com/deesel/wol/internal/config"
	l "github.com/deesel/wol/internal/logger"
)

const version = "1.0.0"

func main() {
	var configFile string
	var logLevel string

	flag.StringVar(&configFile, "config", "config.yml", "configuration file")
	flag.StringVar(&logLevel, "log", "info", "logging level")
	flag.Parse()

	logger := l.New().SetLevel(logLevel)
	logger.Infof("Starting WOL version: %s", version)
	logger.Debugf("Parsing configuration file: %s", configFile)

	cfg, err := config.New(configFile)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Infof("Running API server, listening on: %s:%d", cfg.Server.Address, cfg.Server.Port)
	err = api.New(cfg).Run()
	if err != nil {
		logger.Fatal(err)
	}
}
