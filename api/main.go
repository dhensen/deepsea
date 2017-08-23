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

	// Domain endpoints
	s := r.PathPrefix("/domains").Subrouter()
	s.Methods("GET").HandlerFunc(ListDomains)
	s.Methods("POST").HandlerFunc(BuyDomain)

	// List container presets
	r.HandleFunc("/container-presets", ListContainerPresets).Methods("GET")

	// Container endpoints
	s = r.PathPrefix("/containers").Subrouter()
	r.HandleFunc("/containers", AddContainer).Methods("POST")
	r.HandleFunc("/containers", ListContainers).Methods("GET")

	// Backups
	s = r.PathPrefix("/backups").Subrouter()
	r.HandleFunc("/backups/{id:[0-9]+}", ListBackups).Methods("GET")
	// Pass in a webhook to callback when backup is created
	r.HandleFunc("/backups", CreateBackup).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", r))
}

// Home "/" handler
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}
