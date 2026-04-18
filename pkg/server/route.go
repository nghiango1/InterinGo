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
	"github.com/gorilla/websocket"
)

const API_ROUTE = "/api"
const WS_ROUTE = "/ws"

// This is enforce by build script, which copy over website built static file
// into `content/dist`
const WEBSITE_FILEPATH = "content/dist"

//go:embed all:content
var embedContent embed.FS

// To make Gin (RESTful API server focus) to better handle request as a
// static fileserver
// - Middleware is use, check if user request for a specific route (api, ws)
// - If there no prefix of URL.Path is found, middleware will serve the file
// with the same name, then drop the handle chain via c.Abort()
func (s *Server) registerFileServerMiddleware() {
	fsys, err := fs.Sub(embedContent, WEBSITE_FILEPATH)
	if err != nil {
		log.Fatalln("Failed to embed folder, got error: ", err)
		return
	}

	serveFile := func(c *gin.Context, filePath string, status int) {
		data, err := fs.ReadFile(fsys, filePath)
		if err != nil {
			log.Printf("[ERROR] Error when reading file, got %s", err.Error())
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

	// Middleware register
	s.ginEngine.Use(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, API_ROUTE) || strings.HasPrefix(c.Request.URL.Path, WS_ROUTE) {
			c.Next()
			return
		}

		// Sanitize: clean and strip leading slash
		urlPath := path.Clean(strings.Trim(c.Request.URL.Path, "/"))
		if urlPath == "." {
			urlPath = ""
		}
		candidates := []string{"index.html"}

		// By default, I already build all page as <path>/index.html, this
		// make sure we check others file server options
		// - <path> -> <path> (literal file - js, css, ...)
		// - <path> -> <path>.html (or it could be .html)
		// - <path> -> <path>/index.html (or it could be a path)
		if urlPath != "" {
			candidates = []string{
				urlPath,
				urlPath + ".html",
				urlPath + "/index.html",
			}
		}

		for _, candidate := range candidates {
			println(candidate)
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

		// This have been setup as default fallback, it can handle render
		// 404 page base on url.path via svelte routing support
		serveFile(c, "200.html", http.StatusNotFound)
		c.Abort()
	})
}

func apiRoute(s *Server) {
	s.ginEngine.POST("/api/repl/:id/evaluate", func(c *gin.Context) {
		runtimeId := c.Param("id")
		// Input validate
		var req core.EvaluateRequest
		err := c.BindJSON(&req)
		if err != nil {
			fmt.Println("[ERRPR] API error, can't parse JSON value, got: ", err.Error())
			errorResp := common.NewBadRequestErrorResponse("Invalid JSON", nil)
			c.JSON(http.StatusBadRequest, errorResp)
		}

		req.RuntimeId = runtimeId

		if s.serviceCore == nil {
			fmt.Println("[ERRPR] API error, serviceCore didn't init yet")
			c.JSON(http.StatusInternalServerError, common.NewErrorResponse(500))
		}

		res, evalErr := s.serviceCore.EvaluateHandlerV2(req)

		// Return
		if evalErr != nil {
			c.JSON(evalErr.GetType(), evalErr)
		} else if res != nil {
			c.JSON(http.StatusOK, res)
		}
	})
	s.ginEngine.POST("/api/repl", func(c *gin.Context) {
		// Input validate
		var req core.CreateReplRuntimeRequest
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

		res, evalErr := s.serviceCore.CreateReplRuntime(req)

		// Return
		if evalErr != nil {
			c.JSON(evalErr.GetType(), evalErr)
		} else if res != nil {
			c.JSON(http.StatusOK, res)
		}
	})
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
	pongWait   = 60 * time.Second
	pingPeriod = 50 * time.Second
)

func (s *Server) handleWebSocket(c *gin.Context) {
	s.upgrader.CheckOrigin = func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		return origin == "https://nghiango.asia" || origin == "http://localhost:8080" // Dev
	}
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

	// Start a goroutine to send pings.
	go func() {
		ticker := time.NewTicker(pingPeriod)
		defer ticker.Stop()
		for range ticker.C {
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}()

	client, err := s.serviceCore.WebsocketConnectionCreate(conn)
	if err != nil {
		conn.WriteMessage(websocket.CloseMessage, []byte("LUCKY USER WITH 128 BIT ID COLLISION!!!"))
		conn.Close()
	}

	for {
		messageType, message, err := conn.ReadMessage()

		if err != nil {
			log.Printf("Read error: %v", err)
			break
		}
		if messageType == websocket.TextMessage {
			s.serviceCore.WebsocketMessageHandler(client, message)
		}
		log.Printf("Received: %s", message)

		if messageType == websocket.CloseMessage {
			s.serviceCore.WebsocketConnectionCleanup(client)
		}

		if err != nil {
			log.Printf("Write error: %v", err)
			break
		}
	}
}

func Route(s *Server) {
	s.registerFileServerMiddleware()
	apiRoute(s)

	// setup websocket
	s.ginEngine.GET("/ws", s.handleWebSocket)
}
