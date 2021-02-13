package scheduler

import (
	"context"
	"github.com/ashtanko/octo-server/store"
	"github.com/sirupsen/logrus"
	"time"
)

func Init(ctx context.Context, store store.Store) {
	cfg := NewConfig()
	go scheduler(ctx, cfg, store)
}

func scheduler(ctx context.Context, cfg *Config, store store.Store) {
	ticker := time.NewTicker(cfg.JobRunInterval)
	for {
		select {
		case <-ticker.C:
			err := store.Job().CancelTransactionAndCorrectBalance(cfg.JobRunInterval, 0)
			if err != nil {
				logrus.Error(err)
			}
		case <-ctx.Done():
			ticker.Stop()
			logrus.Info("Job finished")
			return
		}
	}
}
