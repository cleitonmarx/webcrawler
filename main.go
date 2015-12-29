package main

import (
	"os"
	"strings"

	"github.com/cleitonmarx/webcrawler/config"
	"github.com/cleitonmarx/webcrawler/infrastructure"
	"github.com/cleitonmarx/webcrawler/server"
)

func main() {
	var configFileRepository config.Repository
	configFileRepository = infrastructure.NewConfigFileRepository(
		strings.Join([]string{os.Getenv("GOPATH"), "/src/github.com/cleitonmarx/webcrawler/webcrawler.json"}, ""),
		//"../webcrawler.json",
	)
	currentConfig := getCurrentConfig(configFileRepository)
	appServer := server.New(currentConfig)
	appServer.Init()
	appServer.Run()
}

//getCurrentConfig gets the configuration variables from the current environment
func getCurrentConfig(configRepository config.Repository) config.EnvironmentConfig {
	systemConfig, err := configRepository.GetSystemConfiguration()
	handleError(err)

	currentConfig, err := systemConfig.GetCurrentEnvironmentConfig()
	handleError(err)

	return currentConfig
}

//handleError stops the program execution
func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
