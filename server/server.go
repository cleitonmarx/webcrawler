package server

import (
	"log"
	"os"

	"github.com/cleitonmarx/webcrawler/config"
	"github.com/cleitonmarx/webcrawler/server/controllers"
	"github.com/cleitonmarx/webcrawler/services"
	"github.com/codegangsta/negroni"
	"github.com/dimfeld/httptreemux"
)

//Server encapsulates a HTTP server
type Server struct {
	configuration config.EnvironmentConfig
	router        *httptreemux.TreeMux
	HttpServer    *negroni.Negroni
	fetcher       services.UrlFetcher
	logger        *log.Logger
}

//Init initialize all server objects instances
func (s *Server) Init() {

	s.logger = log.New(os.Stdout, "[WebCrawler] ", log.Ldate|log.Ltime)
	mainController := controllers.NewMainController(s.fetcher, s.logger)

	s.router = httptreemux.New()
	s.router.GET("/", mainController.GetHandler)
	s.router.POST("/crawler", mainController.CrawlerHandler)
	s.router.NotFoundHandler = mainController.NotFoundHandler
	s.router.MethodNotAllowedHandler = mainController.MethodNotAllowedHandler

	s.HttpServer = negroni.Classic()
	s.HttpServer.UseHandler(s.router)
}

//Run the web server
func (s *Server) Run() {
	s.logger.Printf("Environment: %s", s.configuration.Name)
	s.HttpServer.Run(s.configuration.HTTPServer.GetFormatedAddress())
}

//getCurrentConfig gets the configuration variables from the current environment
func getCurrentConfig(configRepository config.Repository) (config.EnvironmentConfig, error) {
	systemConfig, err := configRepository.GetSystemConfiguration()
	if err != nil {
		return config.EnvironmentConfig{}, err
	}

	currentConfig, err := systemConfig.GetCurrentEnvironmentConfig()
	if err != nil {
		return config.EnvironmentConfig{}, err
	}

	return currentConfig, nil
}

//New creates a new server instance
func New(configRepository config.Repository, fetcher services.UrlFetcher) (*Server, error) {
	environmentConfig, err := getCurrentConfig(configRepository)
	if err != nil {
		return nil, err
	}
	return &Server{
		configuration: environmentConfig,
		fetcher:       fetcher,
	}, nil
}
