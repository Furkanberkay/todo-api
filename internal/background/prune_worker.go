package background

import (
	"context"
	"log/slog"
	"time"
)

type Worker struct {
	Pruner *TodoPruner
	Logger *slog.Logger
}

func NewPruneWorker(pruner *TodoPruner, logger *slog.Logger) *Worker {
	return &Worker{
		Pruner: pruner,
		Logger: logger,
	}
}

func (w *Worker) StartPruneJob(ctx context.Context, retentionDays int, batchSize int) {
	ticker := time.NewTicker(time.Second * 10)
	for {
		select {
		case <-ctx.Done():
			w.Logger.Error("context is done")
		case <-ticker.C:
			total, err := w.Pruner.Prune(ctx, retentionDays, batchSize)
			if err != nil {
				w.Logger.Error(err.Error())
			}

			if total < batchSize {
				return
			}
		}

	}
}
