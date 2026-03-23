package server

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

// Populating and cache all md pages in docs
var mdPages map[string]string

type Linked struct {
	docs       []string
	nestedLink map[string]*Linked
}

func (lk *Linked) add(fullpath string, filename string) {
	splited := strings.Split(fullpath, "/")
	splited = splited[0 : len(splited)-1]
	var curr *Linked = lk
	for _, steppath := range splited {
		next, ok := curr.nestedLink[steppath]
		if !ok {
			curr.nestedLink[steppath] = &Linked{
				nestedLink: make(map[string]*Linked),
			}
			curr = curr.nestedLink[steppath]
		} else {
			curr = next
		}
	}
	curr.docs = append(curr.docs, filename)
}

var allDocs *Linked

var docsPath = "server/docs/"

func Init() {
}

func isMDextension(fileInfo os.FileInfo) bool {
	splitedName := strings.Split(fileInfo.Name(), ".")
	return len(splitedName) > 1 && splitedName[len(splitedName)-1] == "md"
}

func getPage(fileName string) (string, bool) {
	pageContent, ok := mdPages[fileName]

	if !ok {
		fmt.Println("Can't find entry\n", fileName)
		return "", false
	}

	return pageContent, true
}

func mdToHTML(md []byte) []byte {
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}

var tmplt *template.Template

type News struct {
	Headline string
	Body     string
}

func Start(listenAdrr string) {
	log.Println("Started listening on", listenAdrr)

	// Create a Gin router with default middleware (logger and recovery)
	r := gin.Default()

	// Now start handing data
	Route(r)

	// Spinning up the server.
	err := r.Run(listenAdrr)
	if err != nil {
		log.Fatal(err)
	}
}
