package main

import (
	"flag"
	"os"

	"github.com/p2sub/p2sub/configuration"
	"github.com/p2sub/p2sub/logger"
)

var sugar = logger.GetSugarLogger()

//Load configuration from file
func loadConfig() (conf *configuration.Config, err error) {
	configFile := flag.String("config", "", "Path to configuration file")
	flag.Parse()
	if *configFile == "" {
		flag.Usage()
		os.Exit(0)
	}
	sugar.Info("Load configuration from: ", *configFile)
	return configuration.Import(*configFile)
}
