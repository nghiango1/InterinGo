package server

import (
	"bytes"
	"embed"
	"fmt"
	"interingo/pkg/repl"
	"interingo/pkg/service/common"
	"io/fs"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const API_ROUTE = "/api"

// This is enforce by build script, which copy over website built static file
// into `content/dist`
const WEBSITE_FILEPATH = "content/dist"

//go:embed all:content
var embedContent embed.FS

func EvaluateHandler(c *gin.Context) {
	// Input validate
	var req common.EvalRequest
	errs := c.BindJSON(&req)
	if errs != nil {
		fmt.Println("API error, can't parse form value")
		errorResp := common.NewBadRequestErrorResponse("Invalid JSON", nil)
		c.JSON(http.StatusBadRequest, errorResp)
	}

	log.Printf("[INFO] Eval request, got: %v", req)

	// Handling eval
	buf := bytes.Buffer{}
	repl.Handle(req.Data, &buf)

	// Return
	resp := common.EvalResponseSuccess{
		Output: buf.String(),
	}
	c.JSON(http.StatusOK, resp)
}

// Any call which doesn't match `/api` route will be handle with static fileserver
func pageRoute(r *gin.Engine) {
	fsys, err := fs.Sub(embedContent, WEBSITE_FILEPATH)
	if err != nil {
		log.Fatalln("Failed to embed folder, got error: ", err)
		return
	}

	fileserver := http.FileServer(http.FS(fsys))
	r.Use(
		func(c *gin.Context) {
			isApiCall := strings.HasPrefix(c.Request.URL.Path, API_ROUTE)
			if !isApiCall {
				fileserver.ServeHTTP(c.Writer, c.Request)
				c.Abort()
			}
		})
}

func apiRoute(r *gin.Engine) {
	r.POST("/api/evaluate", EvaluateHandler)
}

func Route(r *gin.Engine) {
	pageRoute(r)
	apiRoute(r)
}
