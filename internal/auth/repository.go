package auth

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"time"
	"todoApp3/internal/domain"

	"gorm.io/gorm"
)

type Repository struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewRepository(db *gorm.DB, logger *slog.Logger) domain.AuthRepository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

func (r *Repository) RegisterUser(ctx context.Context, user *domain.User) error {
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

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
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

func (r *Repository) GetRefreshToken(ctx context.Context, oldToken *string) (*domain.RefreshToken, error) {

	refreshToken := new(domain.RefreshToken)
	if result := r.db.WithContext(ctx).First(refreshToken, "token_hash = ?", oldToken); result.Error != nil {

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		r.logger.Error("refreshToken Error",
			slog.String("error", result.Error.Error()))

		return nil, domain.InternalError
	}

	return refreshToken, nil

}

func (r *Repository) RotateRefreshToken(ctx context.Context, oldTokenID uint, newTokenModel *domain.RefreshToken) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		if err := tx.Model(&domain.RefreshToken{}).Where("id=?", oldTokenID).Update("revoked", &now).Error; err != nil {
			r.logger.Error("failed to revoke old token",
				slog.String("error", err.Error()))
			return domain.InternalError
		}

		if err := tx.Model(&domain.RefreshToken{}).Create(newTokenModel).Error; err != nil {
			r.logger.Error("failed to create new token within tx",
				slog.String("error", err.Error()))
			return domain.InternalError
		}

		return nil
	})
}

func (r *Repository) GetUserByUserID(ctx context.Context, userID uint) (*domain.User, error) {
	user := new(domain.User)

	if result := r.db.WithContext(ctx).First(user, "id = ?", userID); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}

		r.logger.Error("failed to find user",
			slog.Int("user_id", int(userID)),
			slog.String("error", result.Error.Error()),
		)
		return nil, domain.InternalError
	}
	return user, nil
}

func (r *Repository) SaveRefreshToken(ctx context.Context, refreshTokenModel *domain.RefreshToken) error {
	if err := r.db.WithContext(ctx).Model(&domain.RefreshToken{}).Create(refreshTokenModel).Error; err != nil {
		r.logger.Error("refreshToken create error",
			slog.String("component", "gormRepository/saveRefreshToken"),
			slog.String("error", err.Error()))

		return domain.InternalError
	}
	return nil
}

func (r *Repository) DeleteUser(ctx context.Context, id uint) error {

	res := r.db.WithContext(ctx).Delete(&domain.User{}, id)

	if res.Error != nil {
		r.logger.Error("database error during user deletion",
			slog.String("component", "authRepository"),
			slog.Int("user_id", int(id)),
			slog.String("error", res.Error.Error()))

		return domain.InternalError
	}

	if res.RowsAffected == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}
