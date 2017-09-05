package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	c "local/deepsea/api/controllers"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", Index)
	r.HandleFunc("/health", c.Health)

	dir := "static"
	// This will serve files under /static/<filename>
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(dir))))

	// Login
	s := r.PathPrefix("/login").Subrouter()
	s.Methods("POST").HandlerFunc(c.LoginHandler)
	r.HandleFunc("/logout", c.Logout)

	// authentication handler checks for a JWT token
	auth := c.JWTAuth

	// Client endpoints
	s = r.PathPrefix("/clients").Subrouter()
	s.Methods("GET").HandlerFunc(auth(c.GetClients))
	s.Methods("POST").HandlerFunc(auth(c.PostClient))
	s.Methods("DELETE").HandlerFunc(auth(c.DeleteClient))
	s.Methods("PUT").HandlerFunc(auth(c.PutClient))

	// Domain endpoints
	s = r.PathPrefix("/domains").Subrouter()
	s.Methods("GET").HandlerFunc(auth(c.ListDomains))
	s.Methods("POST").HandlerFunc(auth(c.BuyDomain))

	// List container presets
	r.HandleFunc("/container-presets", auth(c.ListContainerPresets)).Methods("GET")

	// Container endpoints
	s = r.PathPrefix("/containers").Subrouter()
	s.Methods("POST").HandlerFunc(auth(c.AddContainer))
	s.Methods("GET").HandlerFunc(auth(c.ListContainers))

	// Backups
	s = r.PathPrefix("/backups/{id:[0-9]+}").Subrouter()
	s.Methods("GET").HandlerFunc(auth(c.ListBackups))
	s.Methods("POST").HandlerFunc(auth(c.CreateBackup))

	var apiPort int64 = 8080
	if value, exists := os.LookupEnv("PORT"); exists {
		apiPort, _ = strconv.ParseInt(value, 10, 0)
	}
	log.Println(fmt.Sprintf("Starting server on port %d", apiPort))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", apiPort), handlers.CORS()(r)))
}

// Home "/" handler
func Index(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
