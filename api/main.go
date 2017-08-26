package main

import (
	"flag"
	"html/template"
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
	t, _ := template.ParseFiles("templates/index.html")
	t.Execute(w, nil)
}

// HasValidJwtToken returns true if the request has a valid JWT token, false otherwise
// Code based on jwtmiddleware.CheckJWT implementation
func HasValidJwtToken(r *http.Request) bool {
	m := GetJWTMiddleware()
	token, err := m.Options.Extractor(r)
	if err != nil {
		return false
	}
	parsedToken, err := jwt.Parse(token, m.Options.ValidationKeyGetter)
	if err != nil {
		return false
	}

	if m.Options.SigningMethod != nil && m.Options.SigningMethod.Alg() != parsedToken.Header["alg"] {
		return false
	}

	return parsedToken.Valid
}

// Home "/login" handler
func Login(w http.ResponseWriter, r *http.Request) {
	if HasValidJwtToken(r) {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	t, _ := template.ParseFiles("templates/login.html")
	t.Execute(w, nil)
}
