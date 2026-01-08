package background

import (
	"context"
	"log/slog"
	"time"
)

type TodoPruneRepo interface {
	Delete(ctx context.Context, cutoff time.Time, limit int) (int, error)
}

type TodoPruner struct {
	repo   TodoPruneRepo
	logger *slog.Logger
}

func NewTodoPruner(repo TodoPruneRepo, logger *slog.Logger) *TodoPruner {
	return &TodoPruner{repo: repo, logger: logger}
}

func (p *TodoPruner) Prune(ctx context.Context, retentionDays int, batchSize int) (int, error) {
	totalDeleted := 0
	cutoff := time.Now().AddDate(0, 0, -retentionDays)

	for {
		if ctx.Err() != nil {
			return totalDeleted, ctx.Err()
		}

		queryCtx, queryCancel := context.WithTimeout(ctx, 1*time.Second)

		count, err := p.repo.Delete(queryCtx, cutoff, batchSize)

		queryCancel()

		if err != nil {
			return totalDeleted, err
		}

		totalDeleted += count
		if count < batchSize {
			break
		}

		select {
		case <-ctx.Done():
			return totalDeleted, ctx.Err()
		case <-time.After(100 * time.Millisecond):
		}
	}
	return totalDeleted, nil
}
