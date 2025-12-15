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
	Success bool
}

func (m *MockVerifier) Verify(ctx context.Context, tokenString string) (*domain.MyCustomClaim, error) {
	myClaim := domain.MyCustomClaim{
		UserID: 1,
		Email:  "berkay123@hotmail.com",
	}

	m.Success = true
	return &myClaim, nil
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

		middleware := AuthenticationMiddleware{Verify: &MockVerifier{}}
		chain := middleware.Authenticate(newHandler)
		err := chain(c)

		if err != nil {
			t.Errorf("unexpected error %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Errorf("ERROR! Expected 200 (OK), Received %d", rec.Code)
		}
	})
}
