package auth

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"todoApp3/internal/domain"

	"gorm.io/gorm"
)

type GormRepository struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewRepository(db *gorm.DB, logger *slog.Logger) domain.AuthRepository {
	return &GormRepository{
		db:     db,
		logger: logger,
	}
}

func (r *GormRepository) RegisterUser(ctx context.Context, user *domain.User) error {
	result := r.db.WithContext(ctx).Create(user)

	if result.Error != nil {
		errStr := result.Error.Error()

		if strings.Contains(errStr, "UNIQUE constraint failed") || strings.Contains(errStr, "Duplicate Entry") {
			return domain.ErrUserAlreadyExists
		}

		r.logger.Error("database error during user creation",
			slog.String("component", "Auth Repository"),
			slog.String("error", errStr),
		)
		return domain.InternalError
	}
	return nil
}

func (r *GormRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	cleanEmail := strings.TrimSpace(email)

	if cleanEmail == "" {
		return nil, domain.InvalidInput
	}

	user := new(domain.User)

	if result := r.db.WithContext(ctx).First(user, "email = ?", cleanEmail); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}

		r.logger.Error("failed to find user",
			slog.String("email", cleanEmail),
			slog.String("error", result.Error.Error()),
		)
		return nil, domain.InternalError
	}
	return user, nil

}

func (r *GormRepository) GetRefreshToken(ctx context.Context, oldToken *string) (*domain.RefreshToken, error) {

	refreshToken := new(domain.RefreshToken)
	if err := r.db.WithContext(ctx).First(refreshToken, "WHERE token_hash = ?", oldToken); err != nil {
		if errors.Is(err.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		r.logger.Error("refreshToken Error",
			slog.String("error", err.Error.Error()))

		return nil, domain.InternalError
	}

	return refreshToken, nil

}

func (r *GormRepository) SaveRefreshToken(ctx context.Context, token *domain.RefreshToken) error {
	if err := r.db.WithContext(ctx).Model(&domain.RefreshToken{}).Create(token); err != nil {
		r.logger.Error("database error during refreshToken creation",
			slog.String("component", "Auth Repository"),
			slog.String("error", err.Error.Error()),
		)
		return domain.InternalError
	}
	return nil
}
