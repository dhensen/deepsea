package models

import (
	"time"
)

type User struct {
	ID         uint       `json:"id" gorm:"primary_key"`
	UUID       string     `json:"uuid"`
	Name       string     `json:"name"`
	Password   string     `json:"password"`
	AccessType string     `json:"access_type"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at" sql:"index"`
}

type Company struct {
	ID        uint       `json:"id" gorm:"primary_key"`
	UUID      string     `json:"uuid"`
	Clients   []Client   `json:"clients" gorm:"many2many:company_clients;"`
	Name      string     `json:"name"`
	KvkNumber string     `json:"kvk_number"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" sql:"index"`
}

type Client struct {
	ID           uint       `json:"id" gorm:"primary_key"`
	UUID         string     `json:"uuid"`
	FirstName    string     `json:"first_name"`
	LastName     string     `json:"last_name"`
	EmailAddress string     `json:"email_address" gorm:"not null;unique"`
	City         string     `json:"city"`
	Postcode     string     `json:"postcode"`
	Address      string     `json:"address"`
	PhoneNumber  string     `json:"phone_number"`
	Companies    []Company  `json:"companies" gorm:"many2many:company_clients;"`
	Domains      Domains    `json:"domains"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at" sql:"index"`
}

type Domain struct {
	ID        uint       `json:"id" gorm:"primary_key"`
	UUID      string     `json:"uuid"`
	Name      string     `json:"name"`
	Provider  string     `json:"provider"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" sql:"index"`
}

type Domains []Domain

type DomainProvider struct {
	ID        uint       `json:"id" gorm:"primary_key"`
	UUID      string     `json:"uuid"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" sql:"index"`
}
