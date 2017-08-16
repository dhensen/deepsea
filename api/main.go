package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html"
	"log"
	"net/http"

	uuid "github.com/google/uuid"
	mux "github.com/gorilla/mux"
)

var kubeconfig *string

func main() {
	kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	flag.Parse()
	if *kubeconfig == "" {
		panic("-kubeconfig not specified")
	}

	router := mux.NewRouter()
	router.HandleFunc("/", Index)

	// Domains
	s := router.PathPrefix("/domains").Subrouter()
	s.Methods("GET").HandlerFunc(ListDomains)
	s.Methods("POST").HandlerFunc(BuyDomain)

	router.HandleFunc("/container_presets", ListContainerPresets).Methods("GET")
	router.HandleFunc("/container", AddContainer).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", router))
}

type Organisation struct {
	ID   uuid.UUID
	Name string
}

type User struct {
	ID   uuid.UUID
	Name string
}

type Domain struct {
	ID       uuid.UUID
	Name     string
	Provider string
}

type Domains []Domain

type DomainProvider struct {
	ID   uuid.UUID
	Name string
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

func ListDomains(w http.ResponseWriter, r *http.Request) {
	domains := Domains{
		Domain{Name: "dinohensen.nl", Provider: "transip"},
		Domain{Name: "ebrandlocal.com", Provider: "transip"},
	}

	json.NewEncoder(w).Encode(domains)
}

// BuyDomain buys a requested domain via a given domain provider
func BuyDomain(w http.ResponseWriter, r *http.Request) {
	domainProvider := r.FormValue("domainProvider")
	domainName := r.FormValue("domainName")

	// TODO: create a Transip domain provider that wraps a transip api client
	log.Printf("Buying domainName %s via domainProvider %s", domainName, domainProvider)

	// Fake it 'till you make it!
	json.NewEncoder(w).Encode(Domain{ID: uuid.New(), Name: domainName, Provider: domainProvider})
}

type ImagePreset struct {
	ID          uuid.UUID
	DisplayName string
	Image       string
}

var imagePresets = map[string]ImagePreset{
	"a8ef3baa-6fbd-4a9b-ad94-247c9273c57d": ImagePreset{
		ID:          uuid.Must(uuid.Parse("a8ef3baa-6fbd-4a9b-ad94-247c9273c57d")),
		DisplayName: "Wordpress 4.8.1 running PHP 5.6 on Apache",
		Image:       "wordpress:4.8.1-php5.6-apache",
	},
}

func ListContainerPresets(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(imagePresets)
}
