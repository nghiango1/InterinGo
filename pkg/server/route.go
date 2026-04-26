package server

import (
	"embed"
	"io/fs"
	"log"
	"log/slog"
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
const FALLBACK_PAGE = "index.html"

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
			// Server content dist may not existed - Making fallback still return 500
			if filePath == FALLBACK_PAGE {
				slog.Warn("Server mode was not packed with embeded WebUI")
				c.Data(status, "plan/text", []byte("Server is up and running"))
			} else {
				slog.Error("Error when reading file", "error", err)
				c.AbortWithStatus(http.StatusInternalServerError)
			}
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

		isFile := func(path string) bool {
			fileInfo, err := fs.Stat(fsys, path)
			if err != nil {
				return false
			}
			return !fileInfo.IsDir()
		}

		for _, candidate := range candidates {
			// fs.FS rejects any path containing ".." at the API level
			if candidate == FALLBACK_PAGE || isFile(candidate) {
				serveFile(c, candidate, http.StatusOK)
				c.Abort()
				return
			}
		}

		// This have been setup as default fallback, it can handle render
		// 404 page base on url.path via svelte routing support
		serveFile(c, FALLBACK_PAGE, http.StatusNotFound)
		c.Abort()
	})
}

func apiRoute(s *Server) {
	// v1 - original one REPL for all appoarch
	s.ginEngine.POST("/api/evaluate", s.serviceV1.EvaluateHandler)

	// v2 - Create multiple REPL for each session
	s.ginEngine.POST("/api/repl/:id/evaluate", s.serviceV2.EvaluateHandler)
	s.ginEngine.POST("/api/repl", s.serviceV2.CreateReplRuntimeHandler)
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
		slog.Error("WebSocket upgrade error", "error", err)
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
			slog.Error("Read error", "error", err)
			break
		}
		if messageType == websocket.TextMessage {
			s.serviceCore.WebsocketReceivedTextMessageHandler(client, message)
		}
		slog.Debug("Websocket received", "message", message)

		if messageType == websocket.CloseMessage {
			s.serviceCore.WebsocketConnectionCleanup(client)
		}

		if err != nil {
			slog.Error("Websocket write error", "error", err)
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
