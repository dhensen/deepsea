package main

import (
	"flag"
	"fmt"
	"html"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"

	"github.com/auth0/go-jwt-middleware"
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
	s := r.PathPrefix("/login").Subrouter()
	s.Methods("GET").HandlerFunc(Login)
	s.Methods("POST").HandlerFunc(LoginHandler)
	r.HandleFunc("/logout", Logout)

	// Domain endpoints
	s = r.PathPrefix("/domains").Subrouter()
	s.Methods("GET").HandlerFunc(Authenticated(ListDomains))
	s.Methods("POST").HandlerFunc(Authenticated(BuyDomain))

	// List container presets
	r.HandleFunc("/container-presets", Authenticated(ListContainerPresets)).Methods("GET")

	// Container endpoints
	s = r.PathPrefix("/containers").Subrouter()
	s.Methods("POST").HandlerFunc(Authenticated(AddContainer))
	s.Methods("GET").HandlerFunc(Authenticated(ListContainers))

	// Backups
	s = r.PathPrefix("/backups/{id:[0-9]+}").Subrouter()
	s.Methods("GET").HandlerFunc(Authenticated(ListBackups))
	s.Methods("POST").HandlerFunc(Authenticated(CreateBackup))

	log.Println("Starting server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func GetJWTMiddleware() *jwtmiddleware.JWTMiddleware {
	return jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
		Extractor:     CookieExtractor(jwtCookieKey),
	})
}

// Home "/" handler
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}
