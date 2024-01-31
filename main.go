package main

import (
	"os"
	"szg/configuration"
	"szg/process"
)

func main() {
	processManager := process.NewProcessManager()

	if len(os.Args) < 2 {
		panic("No configuration file provided")
	}

	configurationFile := os.Args[1]
	if configurationFile == "" {
		panic("No configuration file provided")
	}

	if _, err := os.Stat(configurationFile); os.IsNotExist(err) {
		panic("Configuration file does not exist")
	}

	reader := configuration.NewReader()
	activeConf, err := reader.Read(configurationFile)
	if err != nil {
		panic(err)
	}

	processManager.RunFromConfiguration(activeConf)
	processManager.Loop()
}
