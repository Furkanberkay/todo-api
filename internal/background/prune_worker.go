package background

import (
	"context"
	"log/slog"
	"time"
)

type Worker struct {
	pruner *TodoPruner
	logger *slog.Logger
}

func NewPruneWorker(pruner *TodoPruner, logger *slog.Logger) *Worker {
	return &Worker{
		pruner: pruner,
		logger: logger,
	}
}

func (w *Worker) StartPruneJob(ctx context.Context, retentionDays int, batchSize int) {
	ticker := time.NewTicker(time.Second * 10)
	for {
		select {
		case <-ctx.Done():
			w.logger.Error("context is done")
		case <-ticker.C:
			total, err := w.pruner.Prune(ctx, retentionDays, batchSize)
			if err != nil {
				w.logger.Error(err.Error())
			}

			if total < batchSize {
				return
			}
		}

	}
}
