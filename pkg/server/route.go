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

	"github.com/gin-gonic/gin"
)

const API_ROUTE = "/api"

// This is enforce by build script, which copy over website built static file
// into `content/dist`
const WEBSITE_FILEPATH = "content/dist"

//go:embed all:content
var embedContent embed.FS
var coreInstance *core.Core

func EvaluateHandler(c *gin.Context) {
	// Input validate
	var req core.EvaluateRequest
	errs := c.BindJSON(&req)
	if errs != nil {
		fmt.Println("API error, can't parse form value")
		errorResp := common.NewBadRequestErrorResponse("Invalid JSON", nil)
		c.JSON(http.StatusBadRequest, errorResp)
	}

	if coreInstance == nil {
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(http.StatusInternalServerError))
		return
	}

	// Handling eval
	result, error := coreInstance.Eval(req)
	if error != nil {
		c.JSON(http.StatusBadRequest, error)
	} else {
		c.JSON(http.StatusOK, result)
	}
}

// Any call which doesn't match `/api` route will be handle with static fileserver
func pageRoute(r *gin.Engine) {
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

	r.Use(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, API_ROUTE) {
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

func apiRoute(r *gin.Engine) {
	r.POST("/api/evaluate", EvaluateHandler)
}

func Route(r *gin.Engine) {
	coreInstance = core.NewCore()
	pageRoute(r)
	apiRoute(r)
}
