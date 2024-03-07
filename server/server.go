package server

import (
	"context"
	"fmt"
	"html/template"
	"interingo/repl"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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
	fmt.Println("server init")
	mdPages = make(map[string]string)
	content, err := os.ReadFile("README.md")
	if err == nil {
		mdPages["/docs"] = string(mdToHTML(content))
	}
	allDocs = &Linked{
		nestedLink: make(map[string]*Linked),
	}
	err = filepath.Walk(docsPath, populatedPage)
	if err != nil {
		fmt.Printf("Read docs path error\n")
	}
}

func isMDextension(fileInfo os.FileInfo) bool {
	splitedName := strings.Split(fileInfo.Name(), ".")
	return len(splitedName) > 1 && splitedName[len(splitedName)-1] == "md"
}

func populatedPage(path string, file os.FileInfo, err error) error {
	if file.IsDir() {
		return nil
	}

	if !isMDextension(file) {
		return nil
	}

	docs, err := os.ReadFile(path)
	if err != nil {
		return nil
	}

	mdPages["/"+path] = string(mdToHTML(docs))
	allDocs.add(path, file.Name())
	fmt.Printf("Init poplated page %s\n", path)

	return nil
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

// Handler functions.
func HomeHandle(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(200)
		NotFoundHandler(w, r)
		return
	}

	component := Home()
	component.Render(context.Background(), w)
}

func InfoHandler(w http.ResponseWriter, r *http.Request) {
	component := Info("<p>This is infomation about Authors of InterinGo language<p>")
	info, err := os.ReadFile("server/assets/resume.md")
	if err == nil {
		component = Info(string(mdToHTML(info)))
	}
	component.Render(context.Background(), w)
}

func DocsHandler(w http.ResponseWriter, r *http.Request) {
	component := Docs("<p>This is documentations of InterinGo language<p>", allDocs)

	docs, ok := getPage(r.URL.Path)
	if ok {
		component = Docs(string(docs), allDocs)
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

func populateHandle(path string, linked *Linked) {
	for k := range linked.nestedLink {
		populateHandle(path+"/"+k, linked.nestedLink[k])
	}
	for _, file := range linked.docs {
		http.HandleFunc(path+"/"+file, DocsHandler)
	}
}

func Start(listenAdrr string) {
	log.Println("Started listening on", listenAdrr)

	// Registering our handler functions, and creating paths.
	http.HandleFunc("/", HomeHandle)
	http.HandleFunc("/docs", DocsHandler)
	populateHandle("", allDocs)
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
