package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"testing"
	"time"
	config2 "todoApp3/config"
	"todoApp3/internal/domain"

	"github.com/lmittmann/tint"
	"github.com/stretchr/testify/assert"
)

type MockRepository struct {
	ShouldReturnError bool
	ExistingUser      *domain.User
}

func (m *MockRepository) RegisterUser(ctx context.Context, auth *domain.User) error {
	if m.ShouldReturnError {
		return errors.New("mock db error")
	}
	return nil
}

func (m *MockRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	if m.ShouldReturnError {
		return nil, errors.New("mock db error")
	}
	if m.ExistingUser != nil {
		return m.ExistingUser, nil
	}
	return nil, domain.ErrUserNotFound
}

func (m *MockRepository) GetUserByUserID(ctx context.Context, id uint) (*domain.User, error) {
	return nil, nil
}
func (m *MockRepository) SaveRefreshToken(ctx context.Context, token *domain.RefreshToken) error {
	return nil
}
func (m *MockRepository) GetRefreshToken(ctx context.Context, token *string) (*domain.RefreshToken, error) {
	return nil, nil
}
func (m *MockRepository) RotateRefreshToken(ctx context.Context, oldID uint, newRec *domain.RefreshToken) error {
	return nil
}

func TestRegisterService(t *testing.T) {
	slogHandler := tint.NewHandler(os.Stdout, &tint.Options{
		Level:      slog.LevelDebug,
		TimeFormat: time.TimeOnly,
		AddSource:  true,
	})
	logger := slog.New(slogHandler)

	cfg := &config2.Config{
		JwtExpirationMinutes: 15,
		SecretKey:            "testsecret",
	}

	t.Run("register", func(t *testing.T) {

		input := &CreateUserInput{
			Name:     "Berkay",
			Surname:  "Ozcan",
			Email:    "berkay@test.com",
			Username: "berkayo",
			Password: "12345",
		}

		repo := MockRepository{ExistingUser: nil, ShouldReturnError: false}
		ch := make(chan domain.EmailJob)
		s := NewService(&repo, cfg, logger, ch)
		user, err := s.RegisterUser(context.Background(), input)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "Berkay", user.Name)
		fmt.Println(user.HashedPassword)
		assert.NotEqual(t, input.Password, user.HashedPassword)

	})

}
