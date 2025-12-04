package auth

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
	"todoApp3/config"
	"todoApp3/internal/domain"
)

type Service struct {
	repository domain.AuthRepository
	config     *config.Config
	logger     *slog.Logger
}

func NewService(repository domain.AuthRepository, config *config.Config, logger *slog.Logger) *Service {
	return &Service{
		repository: repository,
		config:     config,
		logger:     logger,
	}
}

func (s *Service) RegisterUser(ctx context.Context, createUserInput *CreateUserInput) (*domain.User, error) {

	if _, err := s.repository.GetUserByEmail(ctx, createUserInput.Email); err == nil {
		return nil, domain.ErrUserAlreadyExists
	}
	password := createUserInput.Password
	hashedPassword, hashError := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if hashError != nil {
		s.logger.Error("failed to hash password", slog.String("error", hashError.Error()))
		return nil, domain.InternalError
	}
	passwordString := string(hashedPassword)
	user := MapCreateUserInputToUser(createUserInput, passwordString)

	if err := s.repository.RegisterUser(ctx, user); err != nil {
		return nil, err
	}

	return user, nil

}

func (s *Service) Login(ctx context.Context, loginInput LoginInput) (string, error) {

	email := loginInput.Email

	user, err := s.repository.GetUserByEmail(ctx, email)

	if err != nil {
		return "", err
	}

	password := loginInput.Password

	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password)); err != nil {
		return "", domain.ErrInvalidCredentials
	}

	claims := MyCustomClaim{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "FurkanBerkayOzcan",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, signedErr := token.SignedString([]byte(s.config.SecretKey))

	if signedErr != nil {
		s.logger.Error("token error",
			slog.String("error", signedErr.Error()))

		return "", domain.InternalError
	}

	return tokenString, nil

}
