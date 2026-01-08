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

	p.logger.Info("prune_cycle_started", "cutoff", cutoff, "batch_size", batchSize)

	for {
		if ctx.Err() != nil {
			p.logger.Warn("prune_cycle_interrupted", "reason", ctx.Err().Error(), "total_deleted", totalDeleted)
			return totalDeleted, ctx.Err()
		}

		count, err := p.repo.Delete(ctx, cutoff, batchSize)
		if err != nil {
			p.logger.Error("prune_repo_delete_failed", "error", err.Error(), "total_deleted", totalDeleted)
			return totalDeleted, err
		}

		totalDeleted += count

		p.logger.Debug("prune_batch_processed", "deleted_count", count, "total_deleted", totalDeleted)
		if count < batchSize {
			p.logger.Info("prune_cycle_completed", "total_deleted", totalDeleted)
			break
		}

		select {
		case <-ctx.Done():
			p.logger.Warn("prune_cycle_cancelled_during_backoff", "total_deleted", totalDeleted)
			return totalDeleted, ctx.Err()
		case <-time.After(100 * time.Millisecond):
		}
	}
	return totalDeleted, nil
}
