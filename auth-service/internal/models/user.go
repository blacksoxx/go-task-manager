package models

import (
	"time"
	"github.com/golang-jwt/jwt/v4"

)

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Password  string    `json:"-"` // Hide from JSON
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignupRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type AuthResponse struct {
	User  *User  `json:"user"`
	Token string `json:"token"`
}

type Claims struct {
    UserID string `json:"user_id"`
    Email  string `json:"email"`
    jwt.StandardClaims
}