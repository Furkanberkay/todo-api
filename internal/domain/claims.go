package domain

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
)

type MyCustomClaim struct {
	UserID uint   `json:"userID"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

type TokenVerify interface {
	Verify(ctx context.Context, tokenString string) (*MyCustomClaim, error)
}
