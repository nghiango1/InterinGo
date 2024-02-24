package server

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var tmplt *template.Template

type News struct {
	Headline string
	Body     string
}

// Port we listen on.
const portNum string = ":8080"

// Handler functions.
func Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Homepage")
}

func Info(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Info page")
}

func HelloHandle(w http.ResponseWriter, r *http.Request) {
	component := Hello("John")
	component.Render(context.Background(), w)
}

func Start() {
	log.Println("Starting our simple http server.")

	// Registering our handler functions, and creating paths.
	http.HandleFunc("/", Home)
	http.HandleFunc("/info", Info)
	http.HandleFunc("/hello", HelloHandle)

	log.Println("Started on port", portNum)

	// Spinning up the server.
	err := http.ListenAndServe(portNum, nil)
	if err != nil {
		log.Fatal(err)
	}
}
