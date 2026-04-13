package server

import (
	"bytes"
	"encoding/json"
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

type BodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w BodyWriter) Write(b []byte) (int, error) {
	return w.body.Write(b)
}

// Middleware so that we can handle or mask non related data to be return to
// client
func messageBodyEnforce(c *gin.Context) {
	wb := &BodyWriter{
		body:           &bytes.Buffer{},
		ResponseWriter: c.Writer,
	}
	c.Writer = wb

	c.Next()

	data := wb.body.String()

	obj := make(map[string]interface{})
	json.Unmarshal([]byte(data), &obj)

	// Make sure that error response have correct number of field
	switch c.Writer.Status() {
	case 400:
		obj["type"] = c.Writer.Status()
		if obj["code"] == nil {
			obj["type"] = "bad_request"
		}
	case 500:
		// Not leak server information
		obj["code"] = "internal_error"
		obj["message"] = "Internal server error"
	}

	updatedBody, _ := json.Marshal(obj)

	wb.ResponseWriter.WriteString(string(updatedBody))
	wb.body.Reset()
}

func (s *Server) Start(listenAdrr string) {
	log.Println("Started listening on", listenAdrr)

	s.ginEngine.Use(messageBodyEnforce)

	// Create a Gin router with default middleware (logger and recovery)
	// Now start handing data
	Route(s)

	// Spinning up the server.
	err := s.ginEngine.Run(listenAdrr)
	if err != nil {
		log.Fatal(err)
	}
}
