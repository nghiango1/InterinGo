package server

import (
	"bytes"
	"embed"
	"fmt"
	"interingo/pkg/repl"
	"interingo/pkg/service/common"
	"log"
	"net/http"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

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

//go:embed all:content
var embedContent embed.FS

func pageRoute(r *gin.Engine) {
	webpage, err := static.EmbedFolder(embedContent, "content/dist")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[INFO] Server static FS `%v` at `%v`\n", "/", "content/dist")
	r.Use(static.Serve("/", webpage))
}

func apiRoute(r *gin.Engine) {
	r.POST("/api/evaluate", EvaluateHandler)
}

func Route(r *gin.Engine) {
	pageRoute(r)
	apiRoute(r)
}
