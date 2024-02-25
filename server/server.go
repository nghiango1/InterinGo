package server

import (
	"context"
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
func HomeHandle(w http.ResponseWriter, r *http.Request) {
	component := Home()
	component.Render(context.Background(), w)
}

func InfoHandle(w http.ResponseWriter, r *http.Request) {
	component := Info()
	component.Render(context.Background(), w)
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
	http.HandleFunc("/", HomeHandle)
	http.HandleFunc("/info", InfoHandle)
	http.HandleFunc("/hello", HelloHandle)
	http.HandleFunc("/greet", GreetingHandle)

	// Static assets file like javascript and css
	staticFileHandler := http.FileServer(http.Dir("./server/assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", staticFileHandler))


	// Spinning up the server.
	err := http.ListenAndServe(listenAdrr, nil)
	if err != nil {
		log.Fatal(err)
	}
}
