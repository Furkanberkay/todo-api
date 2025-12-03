package auth

import (
	"context"
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
	result := r.db.WithContext(ctx).Create(&user)

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
