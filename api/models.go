package main

import uuid "github.com/google/uuid"

type Organisation struct {
	ID   uuid.UUID
	Name string
}

type User struct {
	ID   uuid.UUID
	Name string
}

type Domain struct {
	ID       uuid.UUID
	Name     string
	Provider string
}

type Domains []Domain

type DomainProvider struct {
	ID   uuid.UUID
	Name string
}
