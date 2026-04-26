package server

import (
	"interingo/pkg/service/core"
	service_v1 "interingo/pkg/service/v1"
	service_v2 "interingo/pkg/service/v2"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Server struct {
	ginEngine   *gin.Engine
	serviceCore *core.ServiceCore
	serviceV1   service_v1.InteringoServiceV1
	serviceV2   service_v2.InteringoServiceV2
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

	ginEngine := gin.New()
	ginEngine.Use(gin.Recovery())
	ginEngine.Use(func(c *gin.Context) {
		start := time.Now()

		c.Next()

		slog.Info("http_request",
			"status", c.Writer.Status(),
			"latency", time.Since(start),
			"ip", c.ClientIP(),
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
		)
	})
	return &Server{
		ginEngine:   ginEngine,
		serviceCore: core,
		serviceV1:   service_v1.NewInteringoServiceV1(core),
		serviceV2:   service_v2.NewInteringoServiceV2(core),
		upgrader:    upgrader,
	}
}

func (s *Server) Start(listenAdrr string) {
	slog.Info("Server started listening", "address", listenAdrr)

	// Create a Gin router with default middleware (logger and recovery)
	// Now start handing data
	Route(s)

	// Spinning up the server.
	err := s.ginEngine.Run(listenAdrr)
	if err != nil {
		log.Fatal(err)
	}
}
