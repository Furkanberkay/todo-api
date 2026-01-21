package middleware

import (
	"net/http"
	"todoApp3/internal/domain"
	"todoApp3/internal/httpx"

	"github.com/labstack/echo/v4"
)

type AuthenticationMiddleware struct {
	Verify httpx.TokenVerifier
}

func NewAuthenticationMiddleware(verify httpx.TokenVerifier) *AuthenticationMiddleware {
	return &AuthenticationMiddleware{
		Verify: verify,
	}
}

func (a *AuthenticationMiddleware) Authenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := c.Request().Header.Get("Authorization")
		claim, err := a.Verify.Verify(c.Request().Context(), tokenString)

		if err != nil {
			return c.JSON(http.StatusUnauthorized, httpx.ResponseError{Message: domain.ErrUnAuthorized.Error()})
		}

		c.Set("user", claim)
		return next(c)

	}

}
