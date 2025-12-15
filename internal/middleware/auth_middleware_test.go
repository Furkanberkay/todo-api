package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"todoApp3/internal/domain"

	"github.com/labstack/echo/v4"
)

func TestAuthenticationMiddleware(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderAuthorization, "Bearer mock-token")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		middleware := NewAuthenticationMiddleware(&MockVerifier{})

		nextHandler := func(c echo.Context) error {
			if c.Get("user") == nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "User context'e eklenmedi!")
			}

			return c.String(http.StatusOK, "Giris Basarili")
		}

		handler := middleware.Authenticate(nextHandler)
		_ = handler(c)

		if rec.Code != http.StatusOK {
			t.Errorf("HATA! Beklenen 200 (OK), Gelen %d", rec.Code)
		}
	})

}

type MockVerifier struct{}

func (m *MockVerifier) Verify(ctx context.Context, tokenString string) (*domain.MyCustomClaim, error) {

	myclaim := domain.MyCustomClaim{
		UserID: 1,
		Email:  "berkay123@gmail.com",
	}

	return &myclaim, nil
}
