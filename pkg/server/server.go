package server

import (
	"interingo/pkg/service/core"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Server struct {
	ginEngine   *gin.Engine
	serviceCore *core.ServiceCore
	upgrader    websocket.Upgrader
}

func NewServer() *Server {
	var upgrader = websocket.Upgrader{
		// Allow all origins for development; restrict this in production.
		// reuse GIN_MODE flag
		CheckOrigin: func(r *http.Request) bool {
			return os.Getenv("GIN_MODE") != "release"
		},
	}

	return &Server{
		ginEngine:   gin.Default(),
		serviceCore: core.NewServiceCore(nil),
		upgrader:    upgrader,
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
