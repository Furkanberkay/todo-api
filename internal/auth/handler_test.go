package auth

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
	"todoApp3/internal/domain"

	validator2 "github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/lmittmann/tint"
	"github.com/stretchr/testify/assert"
)

type MockAuthService struct {
	ShouldError bool
}

func (s *MockAuthService) RegisterUser(ctx context.Context, input *CreateUserInput) (*domain.User, error) {
	if s.ShouldError {
		return nil, domain.InternalError
	}

	return &domain.User{
		Name:           input.Name,
		Email:          input.Email,
		Username:       input.Username,
		Surname:        input.Surname,
		HashedPassword: input.Password,
	}, nil
}

func (s *MockAuthService) Login(ctx context.Context, loginInput *LoginInput) (*LoginOutput, error) {
	if s.ShouldError {
		return nil, domain.ErrUnAuthorized
	}

	output := LoginOutput{
		AccessToken:      "asdf",
		RefreshToken:     "2232",
		RefreshExpiresAt: time.Now(),
		ExpiresIn:        2,
	}

	return &output, nil
}

func (s *MockAuthService) RefreshTokens(ctx context.Context, refreshToken string) (*LoginOutput, error) {
	return nil, nil
}

func TestCreateUser(t *testing.T) {

	t.Run("create", func(t *testing.T) {
		e := echo.New()
		jsonBody := `{"name":"Berkay", "email":"berkay@test.com", "password":"verysecretpassword", "username":"furkano","surname":"ozcan"}`

		req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)

		mockService := MockAuthService{ShouldError: false}
		val := validator2.New()

		slogHandler := tint.NewHandler(os.Stdout, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.TimeOnly,
			AddSource:  true,
		})

		slogLogger := slog.New(slogHandler)

		h := NewHandler(&mockService, val, slogLogger)

		err := h.RegisterUser(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Contains(t, rec.Body.String(), "Berkay")
	})

	t.Run("login", func(t *testing.T) {
		e := echo.New()
		jsonBody := `{"email": "berkay@test.com", "password":"ozcan"}`
		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		val := validator2.New()

		slogHandler := tint.NewHandler(os.Stdout, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.TimeOnly,
			AddSource:  true,
		})

		slogLoger := slog.New(slogHandler)

		h := NewHandler(&MockAuthService{ShouldError: false}, val, slogLoger)

		err := h.Login(c)
		fmt.Println(rec.Body.String())

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "refreshToken")

	})
}
