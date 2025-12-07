package middleware

import (
	"todoApp3/config"

	"github.com/labstack/echo/v4"
)

type AuthMiddleware struct {
	config *config.Config
}

func NewAuthMiddleware(config *config.Config) *AuthMiddleware {
	return &AuthMiddleware{
		config: config,
	}
}

func (m *AuthMiddleware) ValidateJwt(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return next(c)
	}
}
