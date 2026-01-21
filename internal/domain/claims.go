package domain

import (
	"github.com/golang-jwt/jwt/v5"
)

type MyCustomClaim struct {
	UserID uint   `json:"userID"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}
