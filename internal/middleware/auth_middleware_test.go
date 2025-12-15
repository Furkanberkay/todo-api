package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"todoApp3/internal/domain"

	"github.com/labstack/echo/v4"
)

type MockVerifier struct {
	ShouldFail bool
}

func (m *MockVerifier) Verify(ctx context.Context, tokenString string) (*domain.MyCustomClaim, error) {
	if m.ShouldFail {
		return nil, domain.ErrInvalidToken
	}

	return &domain.MyCustomClaim{UserID: 1}, nil
}

func TestMockMiddleware(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer mock-token")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		newHandler := func(e echo.Context) error {
			if e.Get("user") == nil {
				return e.JSON(http.StatusInternalServerError, map[string]string{
					"error": "",
				})
			}

			return e.JSON(http.StatusOK, map[string]string{
				"message": "success",
			})

		}

		middleware := AuthenticationMiddleware{Verify: &MockVerifier{ShouldFail: false}}
		chain := middleware.Authenticate(newHandler)
		err := chain(c)

		if err != nil {
			t.Errorf("unexpected error %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Errorf("ERROR! Expected 200 (OK), Received %d", rec.Code)
		}
	})

	t.Run("fail", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer naber")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := func(c echo.Context) error {
			t.Error("ERROR : Next Handler ran while the token was invalid! There's a security vulnerability!")
			return nil
		}

		middleware := AuthenticationMiddleware{Verify: &MockVerifier{ShouldFail: true}}

		chain := middleware.Authenticate(handler)
		_ = chain(c)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("HATA! Beklenen 401 (Unauthorized), Gelen Kod: %d", rec.Code)
		}

	})
}
