package auth

import "github.com/golang-jwt/jwt/v5"

type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Surname  string `json:"surname" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email,max=100"`
	Password string `json:"password" validate:"required,min=3,max=64"`
	Username string `json:"username" validate:"required,alphanum,min=3,max=30"`
}

type RegisterResponse struct {
	ID       uint   `json:"ID"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email,max=100"`
	Password string `json:"password" validate:"required,min=3,max=64"`
}

type MyCustomClaim struct {
	UserID uint   `json:"userID"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}
