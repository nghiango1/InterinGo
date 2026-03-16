package server

import (
	"bytes"
	"embed"
	"fmt"
	"interingo/pkg/repl"
	"interingo/pkg/server/render"
	"interingo/pkg/service/common"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

// Handler functions.
func HomeHandle(c *gin.Context) {
	c.HTML(http.StatusOK, "", Home())
}

func InfoHandler(c *gin.Context) {
	component := Info("<p>This is infomation about Authors of InterinGo language<p>")
	info, err := os.ReadFile("server/assets/resume.md")
	if err == nil {
		c.HTML(http.StatusOK, "", Info(string(mdToHTML(info))))
	} else {
		c.HTML(http.StatusOK, "", component)
	}
}

func NotFoundHandler(c *gin.Context) {
	c.HTML(http.StatusNotFound, "", NotFound())
}

func EvaluateHandler(c *gin.Context) {
	// Input validate
	var req common.EvalRequest
	errs := c.BindJSON(&req)
	if errs != nil {
		fmt.Println("API error, can't parse form value")
	}

	// Handling eval
	buf := bytes.Buffer{}
	repl.Handle(req.Data, &buf)

	// Return
	resp := common.EvalResponse{
		Result: buf.String(),
	}
	c.JSON(http.StatusOK, resp)
}

//go:embed content/**/*
var embedContent embed.FS

func pageRoute(r *gin.Engine) {
	// Isolate assets static file (css, js) from embeded content
	subFS, err := fs.Sub(embedContent, "content/assets")
	if err != nil {
		log.Fatal(err)
	}
	// Gin only serve the file from FS directly, there isn't params to
	// enforce traversal
	// - Direct http.FS(embedContent) will not work
	r.StaticFS("/assets", http.FS(subFS))
    r.NoRoute(NotFoundHandler)

	// traversal(r, "content/dist", "/")
	webpage, err := static.EmbedFolder(embedContent, "content/dist")
	log.Printf("[INFO] Server static FS `%v` at `%v`\n", "/", "content/dist")
	r.Use(static.Serve("/", webpage))

	// Templ render
	ginHtmlRenderer := r.HTMLRender
	r.HTMLRender = &render.HTMLTemplRenderer{FallbackHtmlRenderer: ginHtmlRenderer}

	// Registering our handler functions, and creating paths.
	r.GET("/", HomeHandle)
	// r.GET("/docs", DocsHandler)
	// populateHandle("", allDocs) - Will embed later
	r.GET("/info", InfoHandler)
	r.GET("/404", NotFoundHandler)
}

func apiRoute(r *gin.Engine) {
	r.POST("/api/evaluate", EvaluateHandler)
}

func Route(r *gin.Engine) {
	pageRoute(r)
	apiRoute(r)
}
