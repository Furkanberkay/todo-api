package todo

import (
	"context"
	"log/slog"
	"todoApp3/internal/domain"

	"gorm.io/gorm"
)

type Repository struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewRepository(db *gorm.DB, logger *slog.Logger) domain.TodoRepository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

func (r *Repository) GetTodos(ctx context.Context, page int, limit int, userID uint) ([]domain.Todo, int, error) {
	var todos []domain.Todo
	var totalCount int64

	if err := r.db.WithContext(ctx).Model(&domain.Todo{}).Where("user_id=?", userID).Count(&totalCount).Error; err != nil {
		r.logger.Error("database error during todo count",
			slog.String("component", "TodoRepository"),
			slog.String("error", err.Error()),
		)
		return nil, 0, domain.InternalError
	}

	offset := (page - 1) * limit

	result := r.db.WithContext(ctx).Where("user_id=?", userID).Offset(offset).Limit(limit).Find(&todos)

	if result.Error != nil {
		r.logger.Error("database error during todo list fetch",
			slog.String("component", "TodoRepository"),
			slog.String("error", result.Error.Error()),
		)
		return nil, 0, domain.InternalError
	}
	return todos, int(totalCount), nil

}
