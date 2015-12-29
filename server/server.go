package server

import (
	"log"
	"os"

	"github.com/cleitonmarx/webcrawler/config"
	"github.com/cleitonmarx/webcrawler/server/controllers"
	"github.com/codegangsta/negroni"
	"github.com/dimfeld/httptreemux"
)

//Server encapsulates a HTTP server
type Server struct {
	configuration config.EnvironmentConfig
	router        *httptreemux.TreeMux
	httpServer    *negroni.Negroni
	logger        *log.Logger
}

//Init initialize all server objects instances
func (s *Server) Init() {

	s.logger = log.New(os.Stdout, "[WebCrawler] ", 0)
	mainController := controllers.NewMainController()

	s.router = httptreemux.New()
	s.router.GET("/", mainController.GetHandler)

	s.httpServer = negroni.Classic()
	s.httpServer.UseHandler(s.router)
}

func (s *Server) Run() {
	s.logger.Printf("Environment: %s", s.configuration.Name)

	s.httpServer.Run(s.configuration.HTTPServer.GetFormatedAddress())
}

//New creates a new server instance
func New(environmentConfig config.EnvironmentConfig) *Server {
	return &Server{configuration: environmentConfig}
}
