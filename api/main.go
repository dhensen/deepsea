package main

import (
	"flag"
	"fmt"
	"html"
	"log"
	"net/http"

	mux "github.com/gorilla/mux"
)

var kubeconfig *string

func main() {
	kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	flag.Parse()
	if *kubeconfig == "" {
		panic("-kubeconfig not specified")
	}

	r := mux.NewRouter()
	r.HandleFunc("/", Index)

	dir := "static"
	// This will serve files under http://localhost:8000/static/<filename>
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(dir))))

	// Login
	r.HandleFunc("/login", LoginHandler).Methods("POST")

	// Domain endpoints
	s := r.PathPrefix("/domains").Subrouter()
	s.Methods("GET").HandlerFunc(ListDomains)
	s.Methods("POST").HandlerFunc(BuyDomain)

	// List container presets
	r.HandleFunc("/container-presets", ListContainerPresets).Methods("GET")

	// Container endpoints
	s = r.PathPrefix("/containers").Subrouter()
	s.Methods("POST").HandlerFunc(AddContainer)
	s.Methods("GET").HandlerFunc(ListContainers)

	// Backups
	s = r.PathPrefix("/backups/{id:[0-9]+}").Subrouter()
	s.Methods("GET").HandlerFunc(Authenticated(ListBackups))
	s.Methods("POST").HandlerFunc(CreateBackup)

	log.Println("Starting server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

// Home "/" handler
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

func Authenticated(f func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		f(w, r)
	}
}
