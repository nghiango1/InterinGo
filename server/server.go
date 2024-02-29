package server

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"main/repl"
	"net/http"
	"os"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

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

// Handler functions.
func HomeHandle(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(404)
		NotFoundHandler(w, r)
		return
	}

	component := Home()
	component.Render(context.Background(), w)
}

func InfoHandler(w http.ResponseWriter, r *http.Request) {
	component := Info("<p>This is REPL for InterinGo language<p>")
	info, err := os.ReadFile("server/assets/resume.md")
	if err == nil {
		component = Info(string(mdToHTML(info)))
	}
	component.Render(context.Background(), w)
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	component := NotFound()
	component.Render(context.Background(), w)
}

func EvaluateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		fmt.Fprintf(w, "Not support")
	}

	errs := r.ParseForm()
	if errs != nil {
		fmt.Fprintf(w, "API error, can't parse form value")
		fmt.Println("API error, can't parse form value")
		fmt.Println(r.Form)
	}

	if r.Form.Has("repl-input") {
		r.Form.Get("repl-input")
		repl.Handle(r.Form.Get("repl-input"), w)
	} else {
		fmt.Fprintf(w, "API error, need `repl-input` value to be set")
		fmt.Println("API error, need `repl-input` value to be set")
	}
}

func Start(listenAdrr string) {
	log.Println("Started listening on", listenAdrr)

	// Registering our handler functions, and creating paths.
	http.HandleFunc("/", HomeHandle)
	http.HandleFunc("/info", InfoHandler)
	http.HandleFunc("/404", NotFoundHandler)
	http.HandleFunc("/api/evaluate", EvaluateHandler)

	// Static assets file like javascript and css
	staticFileHandler := http.FileServer(http.Dir("./server/assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", staticFileHandler))

	// Spinning up the server.
	err := http.ListenAndServe(listenAdrr, nil)
	if err != nil {
		log.Fatal(err)
	}
}
