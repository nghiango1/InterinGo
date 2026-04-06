package server

import (
	"interingo/pkg/service/core"
	"log"

	"github.com/gin-gonic/gin"
)

type Server struct {
	ginEngine   *gin.Engine
	serviceCore *core.ServiceCore
}

func NewServer() *Server {
	return &Server{
		ginEngine:   gin.Default(),
		serviceCore: core.NewServiceCore(nil),
	}
}

func (s *Server) Start(listenAdrr string) {
	log.Println("Started listening on", listenAdrr)

	// Create a Gin router with default middleware (logger and recovery)
	// Now start handing data
	Route(s)

	// Spinning up the server.
	err := s.ginEngine.Run(listenAdrr)
	if err != nil {
		log.Fatal(err)
	}
}
