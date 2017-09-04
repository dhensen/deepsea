package models_test

import (
	"local/deepsea/api/models"
	"testing"

	"github.com/gorilla/schema"
)

func TestFillDomainModel(t *testing.T) {
	values := map[string][]string{
		"name":     {"Dino"},
		"provider": {"local"},
	}

	domain := new(models.Domain)
	decoder := schema.NewDecoder()
	decoder.Decode(domain, values)

	assertEqual(t, "Dino", domain.Name)
	assertEqual(t, "local", domain.Provider)
}

func assertEqual(t *testing.T, expected interface{}, actual interface{}) {
	if expected != actual {
		t.Errorf("Expected %v, but got %v", expected, actual)
	}
}
