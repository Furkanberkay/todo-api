package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log/slog"
	"strings"
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
	emailCh    chan<- domain.EmailJob
}

func NewService(repository domain.AuthRepository, config *config.Config, logger *slog.Logger, emailCh chan<- domain.EmailJob) *Service {
	return &Service{
		repository: repository,
		config:     config,
		logger:     logger,
		emailCh:    emailCh,
	}
}

func (s *Service) Register(ctx context.Context, createUserInput *CreateUserInput) (*domain.User, error) {

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

	msg := fmt.Sprintf("Welcome %s %s", user.Name, user.Surname)

	select {
	case s.emailCh <- domain.EmailJob{
		Email:      user.Email,
		Name:       user.Name,
		Surname:    user.Surname,
		EnqueuedAt: time.Now(),
		Message:    msg,
	}:
	case <-ctx.Done():
		s.logger.Warn("User registered but welcome email skipped due to context cancellation",
			"email", user.Email,
			"error", ctx.Err())

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
		Role:   user.Role,
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

	if strings.TrimSpace(refreshToken) == "" {
		return nil, domain.ErrUnAuthorized
	}

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
		return nil, serviceErr
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

func (s *Service) Delete(ctx context.Context, role string, id uint) error {
	if role != "admin" {
		return domain.ErrForbidden
	}

	if err := s.repository.DeleteUser(ctx, id); err != nil {
		return err
	}
	return nil
}
