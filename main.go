package main

import (
	"os"
	"runtime"
	"strings"

	"github.com/cleitonmarx/webcrawler/infrastructure"
	"github.com/cleitonmarx/webcrawler/server"
)

func main() {
	//Use all CPU in goroutines
	runtime.GOMAXPROCS(runtime.NumCPU())

	configFileRepository := infrastructure.NewConfigFileRepository(
		strings.Join([]string{os.Getenv("GOPATH"), "/src/github.com/cleitonmarx/webcrawler/webcrawler.json"}, ""),
	)
	fetcher := infrastructure.NewHttpClientFetcher()

	appServer, err := server.New(configFileRepository, fetcher)
	handleError(err)
	appServer.Init()
	appServer.Run()

}

//handleError stops the program execution
func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
