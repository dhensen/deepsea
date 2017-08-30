package main

import uuid "github.com/google/uuid"

type User struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Password   string    `json:"password"`
	AccessType string    `json:"access_type"`
}

type Domain struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Provider string    `json:"provider"`
}

type Domains []Domain

type DomainProvider struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
