package main

import (
	"github.com/p2sub/p2sub/configuration"
)

func main() {
	//Read configuration from agrument & file
	config, err := loadConfig()
	if err == nil {
		if config.NodeType == configuration.NodeNotary {
			startNotary(config)
		} else {
			startMaster(config)
		}
	} else {
		sugar.Fatal(err)
	}
}
