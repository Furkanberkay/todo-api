package auth

import "todoApp3/internal/domain"

func MapCreateUserInputToUser(userRegisterInput *CreateUserInput, hashedPassword string) *domain.User {
	return &domain.User{
		Name:         userRegisterInput.Name,
		Surname:      userRegisterInput.Surname,
		Username:     userRegisterInput.Username,
		Email:        userRegisterInput.Email,
		PasswordHash: hashedPassword,
	}
}
