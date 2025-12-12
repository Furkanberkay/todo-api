package auth

import "time"

type CreateUserInput struct {
	Name     string
	Surname  string
	Username string
	Email    string
	Password string
}

type LoginInput struct {
	Email    string
	Password string
}

type LoginOutput struct {
	AccessToken            string
	ExpiresIn              int
	RefreshToken           string
	RefreshTokenExpiresDay time.Time
}
