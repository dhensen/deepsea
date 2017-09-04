package models

import (
	uuid "github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type User struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Password   string    `json:"password"`
	AccessType string    `json:"access_type"`
}

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
	Domains      Domains
}

type Domain struct {
	gorm.Model
	UUID     uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Provider string    `json:"provider"`
}

type Domains []Domain

type DomainProvider struct {
	UUID uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
