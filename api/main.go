package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"local/deepsea/api/controllers"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", Index)
	r.HandleFunc("/health", controllers.Health)

	dir := "static"
	// This will serve files under /static/<filename>
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(dir))))

	// Login
	s := r.PathPrefix("/login").Subrouter()
	s.Methods("POST").HandlerFunc(controllers.LoginHandler)
	r.HandleFunc("/logout", controllers.Logout)

	// Client endpoints
	s = r.PathPrefix("/clients").Subrouter()
	s.Methods("GET").HandlerFunc(controllers.GetClients)
	s.Methods("POST").HandlerFunc(controllers.PostClient)
	s.Methods("DELETE").HandlerFunc(controllers.DeleteClient)
	s.Methods("PUT").HandlerFunc(controllers.PutClient)

	// Domain endpoints
	s = r.PathPrefix("/domains").Subrouter()
	s.Methods("GET").HandlerFunc(controllers.Authenticated(controllers.ListDomains))
	s.Methods("POST").HandlerFunc(controllers.Authenticated(controllers.BuyDomain))

	// List container presets
	r.HandleFunc("/container-presets", controllers.Authenticated(controllers.ListContainerPresets)).Methods("GET")

	// Container endpoints
	s = r.PathPrefix("/containers").Subrouter()
	s.Methods("POST").HandlerFunc(controllers.Authenticated(controllers.AddContainer))
	s.Methods("GET").HandlerFunc(controllers.Authenticated(controllers.ListContainers))

	// Backups
	s = r.PathPrefix("/backups/{id:[0-9]+}").Subrouter()
	s.Methods("GET").HandlerFunc(controllers.Authenticated(controllers.ListBackups))
	s.Methods("POST").HandlerFunc(controllers.Authenticated(controllers.CreateBackup))

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
