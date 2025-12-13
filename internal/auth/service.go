package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log/slog"
	"time"
	"todoApp3/config"
	"todoApp3/internal/domain"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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

func (s *Service) Login(ctx context.Context, loginInput *LoginInput) (*LoginOutput, error) {

	email := loginInput.Email

	user, err := s.repository.GetUserByEmail(ctx, email)

	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	password := loginInput.Password

	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password)); err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	accessToken, accessTokenErr := s.generateAccessToken(user)

	if accessTokenErr != nil {
		return nil, accessTokenErr
	}

	refreshToken, refreshTokenModel := s.createRefreshTokenModel(user.ID)

	if err := s.repository.SaveRefreshToken(ctx, refreshTokenModel); err != nil {
		return nil, err
	}
	expirationDuration := time.Duration(s.config.JwtExpirationMinutes) * (time.Minute)

	output := LoginOutput{
		RefreshToken:     refreshToken,
		AccessToken:      accessToken,
		ExpiresIn:        int(expirationDuration.Seconds()),
		RefreshExpiresAt: refreshTokenModel.ExpiresAt,
	}

	return &output, nil

}

func (s *Service) generateAccessToken(user *domain.User) (string, error) {
	expirationDuration := time.Duration(s.config.JwtExpirationMinutes) * (time.Minute)

	myClaim := domain.MyCustomClaim{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expirationDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    s.config.AppName,
		},
	}

	tokenString := jwt.NewWithClaims(jwt.SigningMethodHS256, myClaim)
	token, err := tokenString.SignedString([]byte(s.config.SecretKey))

	if err != nil {
		s.logger.Error("failed to sign access token",
			slog.String("error", err.Error()),
			slog.String("user_id", fmt.Sprint(user.ID)),
		)
		return "", err
	}

	return token, nil
}

func (s *Service) createRefreshTokenModel(userID uint) (string, *domain.RefreshToken) {
	refreshToken := uuid.New().String()

	hashBytes := sha256.Sum256([]byte(refreshToken))
	tokenHashHex := hex.EncodeToString(hashBytes[:])

	refreshTokenStruct := domain.RefreshToken{
		UserID:    userID,
		TokenHash: tokenHashHex,
		ExpiresAt: time.Now().Add(time.Hour * 24 * time.Duration(s.config.RefReshExpDays)),
	}

	return refreshToken, &refreshTokenStruct
}

func (s *Service) RefreshTokens(ctx context.Context, refreshToken string) (*LoginOutput, error) {
	hashBytes := sha256.Sum256([]byte(refreshToken))
	tokenHashHex := hex.EncodeToString(hashBytes[:])

	oldTokenModel, err := s.repository.GetRefreshToken(ctx, &tokenHashHex)

	if err != nil {
		return nil, domain.ErrUnAuthorized
	}

	if time.Now().After(oldTokenModel.ExpiresAt) {
		return nil, domain.ErrUnAuthorized
	}

	if oldTokenModel.Revoked != nil {
		return nil, domain.ErrUnAuthorized
	}

	user, serviceErr := s.repository.GetUserByUserID(ctx, oldTokenModel.UserID)

	if serviceErr != nil {
		return nil, err
	}

	accessToken, accessTokenErr := s.generateAccessToken(user)

	if accessTokenErr != nil {
		return nil, accessTokenErr
	}

	newRefreshToken, newRefreshTokenModel := s.createRefreshTokenModel(user.ID)

	if err := s.repository.RotateRefreshToken(ctx, oldTokenModel.ID, newRefreshTokenModel); err != nil {
		return nil, err
	}

	expirationDuration := time.Duration(s.config.JwtExpirationMinutes) * (time.Minute)

	output := LoginOutput{
		RefreshToken:     newRefreshToken,
		AccessToken:      accessToken,
		ExpiresIn:        int(expirationDuration.Seconds()),
		RefreshExpiresAt: newRefreshTokenModel.ExpiresAt,
	}

	return &output, nil

}
