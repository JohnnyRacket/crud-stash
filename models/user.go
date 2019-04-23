package models

import (
	"github.com/google/uuid"
)

//user represents a user in the system, can own apps
type User struct {
	Id   uuid.UUID
	Name string
	Apps []string
}
