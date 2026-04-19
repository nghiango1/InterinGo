package server

import (
	"interingo/pkg/service/core"
	service_v1 "interingo/pkg/service/v1"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Server struct {
	ginEngine   *gin.Engine
	serviceCore *core.ServiceCore
	serviceV1   service_v1.InteringoServiceV1
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

	core := core.NewServiceCore(nil)

	return &Server{
		ginEngine:   gin.Default(),
		serviceCore: core,
		serviceV1:   service_v1.NewInteringoServiceV1(core),
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
