package main

import (
	"flag"

	"github.com/deesel/wol/internal/api"
	"github.com/deesel/wol/internal/config"
	l "github.com/deesel/wol/internal/logger"
)

const VERSION = "1.0.0"

func main() {
	var configFile string
	var logLevel string

	flag.StringVar(&configFile, "config", "config.yml", "configuration file")
	flag.StringVar(&logLevel, "log", "info", "logging level")
	flag.Parse()

	logger, err := l.New(logLevel)
	if err != nil {
		panic(err)
	}

	logger.Infof("Starting WOL version: %s", VERSION)
	logger.Debugf("Parsing configuration file: %s", configFile)

	cfg, err := config.New(configFile)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Info("Running API server")
	api.New(cfg).Run()
}
