package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/dgrijalva/jwt-go"

	"github.com/auth0/go-jwt-middleware"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
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
	// s.Methods("OPTIONS").HandlerFunc(CorsHandler)
	s.Methods("GET").HandlerFunc(Login)
	s.Methods("POST").HandlerFunc(LoginHandler)
	r.HandleFunc("/logout", Logout)

	// Client endpoints
	s = r.PathPrefix("/clients").Subrouter()
	// s.Methods("OPTIONS").HandlerFunc(CorsHandler)
	s.Methods("GET").HandlerFunc(GetClients)
	s.Methods("POST").HandlerFunc(PostClient)
	s.Methods("DELETE").HandlerFunc(DeleteClient)
	s.Methods("PUT").HandlerFunc(PutClient)

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

	var apiPort int64 = 8080
	if value, exists := os.LookupEnv("PORT"); exists {
		apiPort, _ = strconv.ParseInt(value, 10, 0)
	}
	log.Println(fmt.Sprintf("Starting server on port %d", apiPort))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", apiPort), handlers.CORS()(r)))
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
