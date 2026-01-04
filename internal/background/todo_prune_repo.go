package background

import (
	"context"
	"time"

	"todoApp3/internal/domain"

	"gorm.io/gorm"
)

type TodoPruneRepoDB struct {
	db *gorm.DB
}

func NewTodoPruneRepoDB(db *gorm.DB) *TodoPruneRepoDB {
	return &TodoPruneRepoDB{db: db}
}

func (r *TodoPruneRepoDB) DeleteCompletedBefore(ctx context.Context, cutoff time.Time, limit int) (int, error) {
	if limit <= 0 {
		limit = 1000
	}

	res := r.db.WithContext(ctx).
		Unscoped().
		Where("completed = ? AND created_at < ?", true, cutoff).
		Limit(limit).
		Delete(&domain.Todo{})

	if res.Error != nil {
		return 0, res.Error
	}

	return int(res.RowsAffected), nil
}
