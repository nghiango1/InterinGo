package server

import (
	"embed"
	"fmt"
	"interingo/pkg/service/common"
	"interingo/pkg/service/core"
	"io/fs"
	"log"
	"mime"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const API_ROUTE = "/api"
const WS_ROUTE = "/ws"

// This is enforce by build script, which copy over website built static file
// into `content/dist`
const WEBSITE_FILEPATH = "content/dist"

//go:embed all:content
var embedContent embed.FS

// Any call which doesn't match `/api` route will be handle with static fileserver
func pageRoute(s *Server) {
	fsys, err := fs.Sub(embedContent, WEBSITE_FILEPATH)
	if err != nil {
		log.Fatalln("Failed to embed folder, got error: ", err)
		return
	}

	serveFile := func(c *gin.Context, filePath string, status int) {
		data, err := fs.ReadFile(fsys, filePath)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		ext := path.Ext(filePath)
		mimeType := mime.TypeByExtension(ext)
		if mimeType == "" {
			mimeType = "application/octet-stream"
		}
		c.Data(status, mimeType, data)
	}

	s.ginEngine.Use(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, API_ROUTE) {
			c.Next()
			return
		}
		if strings.HasPrefix(c.Request.URL.Path, WS_ROUTE) {
			c.Next()
			return
		}

		// Sanitize: clean and strip leading slash
		urlPath := path.Clean(strings.Trim(c.Request.URL.Path, "/"))
		if urlPath == "." {
			urlPath = ""
		}

		candidates := []string{"index.html"}
		if urlPath != "" {
			candidates = []string{
				urlPath,
				urlPath + ".html",
				urlPath + "/index.html",
			}
		}

		for _, candidate := range candidates {
			// fs.FS rejects any path containing ".." at the API level
			fileInfo, err := fs.Stat(fsys, candidate)
			if err != nil {
				continue
			}
			if !fileInfo.IsDir() {
				serveFile(c, candidate, http.StatusOK)
				c.Abort()
				return

			}
		}

		serveFile(c, "404.html", http.StatusNotFound)
		c.Abort()
	})
}

func apiRoute(s *Server) {
	s.ginEngine.POST("/api/evaluate", func(c *gin.Context) {
		// Input validate
		var req core.EvaluateRequest
		err := c.BindJSON(&req)
		if err != nil {
			fmt.Println("[ERRPR] API error, can't parse JSON value, got: ", err.Error())
			errorResp := common.NewBadRequestErrorResponse("Invalid JSON", nil)
			c.JSON(http.StatusBadRequest, errorResp)
		}

		if s.serviceCore == nil {
			fmt.Println("[ERRPR] API error, serviceCore didn't init yet")
			c.JSON(http.StatusInternalServerError, common.NewErrorResponse(500))
		}

		res, evalErr := s.serviceCore.EvaluateHandler(req)

		// Return
		if evalErr != nil {
			c.JSON(evalErr.GetType(), evalErr)
		} else if res != nil {
			c.JSON(http.StatusOK, res)
		}
	})
}

const (
	pongWait = 70 * time.Second
	// Client will ping, server don't expect to check health
	// pingPeriod = 60 * time.Second
)

func (s *Server) handleWebSocket(c *gin.Context) {
	conn, err := s.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		messageType, message, err := conn.ReadMessage()

		if err != nil {
			log.Printf("Read error: %v", err)
			break
		}
		log.Printf("Received: %s", message)

		s.serviceCore.WebsocketHandler(conn, messageType, message)

		// We not expect to return anything back to client
		// if err := conn.WriteMessage(messageType, message); err != nil {
		// 	log.Printf("Write error: %v", err)
		// 	break
		// }
	}
}

func Route(s *Server) {
	pageRoute(s)
	apiRoute(s)

	// setup websocket
	s.ginEngine.GET("/ws", s.handleWebSocket)
}
