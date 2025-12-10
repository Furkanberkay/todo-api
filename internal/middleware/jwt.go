package middleware

import (
	"net/http"
	"strings"
	"todoApp3/config"
	"todoApp3/internal/domain"

	"github.com/golang-jwt/jwt/v5"
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
		authHeader := c.Request().Header.Get("Authorization")

		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "missing authorization header",
			})
		}

		if len(authHeader) < 7 || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "invalid token format",
			})
		}

		tokenString := authHeader[7:]

		myClaim := domain.MyCustomClaim{}

		token, err := jwt.ParseWithClaims(tokenString, &myClaim, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, echo.NewHTTPError(http.StatusUnauthorized, map[string]string{
					"error": "unexpected signing method",
				})
			}

			return []byte(m.config.SecretKey), nil
		})

		if err != nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "invalid or expired token",
			})
		}

		c.Set("user", &myClaim)

		return next(c)
	}
}
