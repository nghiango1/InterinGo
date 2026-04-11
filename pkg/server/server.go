package server

import (
	"interingo/pkg/service/core"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	// Allow all origins for development; restrict this in production.
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Read error: %v", err)
			break
		}
		log.Printf("Received: %s", message)

		if err := conn.WriteMessage(messageType, message); err != nil {
			log.Printf("Write error: %v", err)
			break
		}
	}
}

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

	// setup websocket
	s.ginEngine.GET("/ws", handleWebSocket)

	// Spinning up the server.
	err := s.ginEngine.Run(listenAdrr)
	if err != nil {
		log.Fatal(err)
	}
}
