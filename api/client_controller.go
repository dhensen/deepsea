package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Company struct {
	gorm.Model
	Clients   []Client `gorm:"many2many:company_clients;"`
	Name      string
	KvkNumber string
}

type Client struct {
	gorm.Model
	FirstName    string
	LastName     string
	EmailAddress string `gorm:"not null;unique"`
	City         string
	Postcode     string
	Address      string
	PhoneNumber  string
	Companies    []Company `gorm:"many2many:company_clients;"`
}

var db *gorm.DB

func init() {
	db = getDB()
	db.AutoMigrate(&Client{}, &Company{})
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/client", ClientsHandler())
	log.Fatal(http.ListenAndServe(":8081", mux))
}

func getDB() *gorm.DB {
	db, err := gorm.Open("mysql", "root:deepsea@tcp(localhost:3307)/test_deepsea?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	return db
}

func ClientsHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodHead:
			fallthrough
		case http.MethodGet:
			GetClients(w, r)
		case http.MethodPost:
			PostClient(w, r)
		case http.MethodDelete:
			DeleteClient(w, r)
		case http.MethodPut:
			PutClient(w, r)
		default:
			MethodNotAllowed(w, r)
		}
	})
}

func MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func GetClients(w http.ResponseWriter, r *http.Request) {
	var clients []Client
	db.Find(&clients)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(clients)
}

func PostClient(w http.ResponseWriter, r *http.Request) {
	firstName := r.PostFormValue("firstName")
	lastName := r.PostFormValue("lastName")
	emailAddress := r.PostFormValue("emailAddress")
	city := r.PostFormValue("city")
	postcode := r.PostFormValue("postcode")
	address := r.PostFormValue("address")
	phoneNumber := r.PostFormValue("phoneNumber")

	companyName := r.PostFormValue("companyName")
	kvkNumber := r.PostFormValue("kvkNumber")

	company := Company{
		Name:      companyName,
		KvkNumber: kvkNumber,
	}

	client := Client{
		FirstName:    firstName,
		LastName:     lastName,
		EmailAddress: emailAddress,
		City:         city,
		Postcode:     postcode,
		Address:      address,
		PhoneNumber:  phoneNumber,
		Companies:    []Company{company},
	}

	db.Create(&client)

	errors := db.GetErrors()
	if len(errors) > 0 {
		log.Println(errors)
		w.WriteHeader(http.StatusBadRequest)
		// TODO: return a 409 Conflict if the user already exists
		// w.WriteHeader(http.StatusConflict)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(client)
}

func PutClient(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("not implemened yet")
}

func DeleteClient(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("not implemened yet")
}
