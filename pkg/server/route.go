package server

import (
	"bytes"
	"fmt"
	"interingo/pkg/repl"
	"interingo/pkg/server/render"
	"interingo/pkg/service/common"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// Handler functions.
func HomeHandle(c *gin.Context) {
	c.HTML(http.StatusOK,"", Home())
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
	c.HTML(http.StatusOK, "", NotFound())
}

func EvaluateHandler(c *gin.Context) {
	// Input validate
	var req common.EvalRequest;
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

func Route(r *gin.Engine) {
	ginHtmlRenderer := r.HTMLRender
	r.HTMLRender = &render.HTMLTemplRenderer{FallbackHtmlRenderer: ginHtmlRenderer}

	// Registering our handler functions, and creating paths.
	r.GET("/", HomeHandle)
	// r.GET("/docs", DocsHandler)
	// populateHandle("", allDocs) - Will embed later
	r.GET("/info", InfoHandler)
	r.GET("/404", NotFoundHandler)
	r.POST("/api/evaluate", EvaluateHandler)


	// Static assets file like javascript and css
	staticFileHandler := http.FileServer(http.Dir("./server/assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", staticFileHandler))

}
