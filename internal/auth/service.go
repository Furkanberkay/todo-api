package auth

import (
	"context"
	"todoApp3/internal/domain"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repository domain.AuthRepository
}

func NewService(repository domain.AuthRepository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) RegisterUser(ctx context.Context, createUserInput *CreateUserInput) (*domain.User, error) {
	password := createUserInput.Password
	hashedPassword, hashError := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if hashError != nil {
		return nil, domain.InternalError
	}
	passwordString := string(hashedPassword)
	user := MapCreateUserInputToUser(createUserInput, passwordString)

	if err := s.repository.RegisterUser(ctx, user); err != nil {
		return nil, err
	}

	return user, nil

}
