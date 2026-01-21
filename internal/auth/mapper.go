package auth

import "todoApp3/internal/domain"

func (r *RegisterRequest) ToRegisterInput() *CreateUserInput {
	return &CreateUserInput{
		Username: r.Username,
		Surname:  r.Surname,
		Name:     r.Name,
		Email:    r.Email,
		Password: r.Password,
	}
}

func NewRegisterResponse(u *domain.User) *RegisterResponse {
	return &RegisterResponse{
		ID:       u.ID,
		Name:     u.Name,
		Surname:  u.Surname,
		Username: u.Username,
		Email:    u.Email,
	}
}

func (dto *LoginRequest) ToLoginInput() *LoginInput {
	return &LoginInput{
		Email:    dto.Email,
		Password: dto.Password,
	}
}

func MapCreateUserInputToUser(input *CreateUserInput, hashedPassword string) *domain.User {
	return &domain.User{
		Name:           input.Name,
		Surname:        input.Surname,
		Username:       input.Username,
		Email:          input.Email,
		Role:           "user",
		HashedPassword: hashedPassword,
	}
}
