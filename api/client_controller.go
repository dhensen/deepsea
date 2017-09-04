package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

func init() {
	db = getDB()
	db.AutoMigrate(&Client{}, &Company{})
}

func getDB() *gorm.DB {
	db, err := gorm.Open("mysql", "root:deepsea@tcp(localhost:3307)/test_deepsea?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	return db
}

func GetClients(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	var clients []Client
	db.Find(&clients)

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

	var company Company
	db.FirstOrCreate(&company, Company{
		Name:      companyName,
		KvkNumber: kvkNumber,
	})

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
	// Prerequisites:
	// 1. all resources used by clients must be deleted
	// 2. all outstanding invoices must be payed by client
	// 3. contract duration is expired OR contract is paid off
	// When all prerequisites are met, soft delete the client
	json.NewEncoder(w).Encode("not implemened yet")
}
