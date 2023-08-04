// Package model contains models of using entities
package model

import "github.com/google/uuid"

// User contains an info about the user and will be written in a users table
type User struct {
	ID           uuid.UUID `json:"-"`
	Login        string    `json:"username" validate:"required,min=3,max=15"`
	Password     []byte    `json:"password" validate:"required,min=5,max=15"`
	RefreshToken string    `json:"-" `
}
