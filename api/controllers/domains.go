package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"local/deepsea/api/models"

	uuid "github.com/google/uuid"
)

func ListDomains(w http.ResponseWriter, r *http.Request) {
	domains := models.Domains{
		models.Domain{Name: "dinohensen.nl", Provider: "transip"},
		models.Domain{Name: "ebrandlocal.com", Provider: "transip"},
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
	json.NewEncoder(w).Encode(models.Domain{UUID: uuid.New().String(), Name: domainName, Provider: domainProvider})
}
