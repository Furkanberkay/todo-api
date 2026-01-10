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
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			w.logger.Info("prune_worker_shutdown", "reason", ctx.Err())
			return
		case <-ticker.C:
			w.logger.Info("prune_job_triggered")
			jobCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			total, err := w.pruner.Prune(jobCtx, retentionDays, batchSize)
			cancel()
			if err != nil {
				w.logger.Error(err.Error())
				continue
			}

			if total > 0 {
				w.logger.Info("prune_job_success", "total_deleted", total)
			} else {
				w.logger.Info("prune_job_no_op", "msg", "no_records_found")
			}
		}

	}
}
