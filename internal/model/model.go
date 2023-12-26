// Package model contains models of using entities
package model

import "github.com/google/uuid"

// User contains an info about the user and will be written in a users table
type User struct {
	ID           uuid.UUID `json:"-"`
	Login        string    `json:"login" validate:"required,min=5,max=20"`
	Password     []byte    `json:"password" validate:"required,min=8"`
}
