package background

import (
	"context"
	"log/slog"
	"time"
)

type TodoPruneRepo interface {
	DeleteCompletedBefore(ctx context.Context, cutoff time.Time, limit int) (int, error)
}

type TodoPruner struct {
	repo   TodoPruneRepo
	logger *slog.Logger
}

func NewTodoPruner(repo TodoPruneRepo, logger *slog.Logger) *TodoPruner {
	return &TodoPruner{repo: repo, logger: logger}
}

func (p *TodoPruner) Prune(ctx context.Context, retentionDays int, batchSize int) (int, error) {
	if retentionDays <= 0 {
		retentionDays = 30
	}
	if batchSize <= 0 {
		batchSize = 1000
	}

	cutoff := time.Now().AddDate(0, 0, -retentionDays)

	p.logger.Info("prune started", "retention_days", retentionDays, "batch_size", batchSize)

	total := 0
	for {
		if err := ctx.Err(); err != nil {
			p.logger.Warn("prune stopped", "reason", err.Error(), "deleted_count", total)
			return total, err
		}

		deleted, err := p.repo.DeleteCompletedBefore(ctx, cutoff, batchSize)
		if err != nil {
			p.logger.Error("prune failed", "err", err.Error(), "deleted_count", total)
			return total, err
		}

		total += deleted
		if deleted < batchSize {
			break
		}

		select {
		case <-ctx.Done():
			p.logger.Warn("prune stopped", "reason", ctx.Err().Error(), "deleted_count", total)
			return total, ctx.Err()
		case <-time.After(100 * time.Millisecond):
		}
	}

	p.logger.Info("prune finished", "deleted_count", total)
	return total, nil
}
