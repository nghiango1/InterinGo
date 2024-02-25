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

// Handler functions.
func Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Homepage")
}

func Info(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Info page")
}

func HelloHandle(w http.ResponseWriter, r *http.Request) {
	component := Hello("Lenn")
	component.Render(context.Background(), w)
}

func GreetingHandle(w http.ResponseWriter, r *http.Request) {
	person := Person{
		Name: "Michel",
	}
	component := Greeting(person)
	component.Render(context.Background(), w)
}

func Start(listenAdrr string) {
	log.Println("Started listening on", listenAdrr)

	// Registering our handler functions, and creating paths.
	http.HandleFunc("/", Home)
	http.HandleFunc("/info", Info)
	http.HandleFunc("/hello", HelloHandle)
	http.HandleFunc("/greet", GreetingHandle)

	// Spinning up the server.
	err := http.ListenAndServe(listenAdrr, nil)
	if err != nil {
		log.Fatal(err)
	}
}
