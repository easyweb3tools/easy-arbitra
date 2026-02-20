package worker

import (
	"context"
	"time"

	"go.uber.org/zap"
)

type Syncer interface {
	Name() string
	RunOnce(ctx context.Context) error
}

type ScheduledSyncer struct {
	Syncer   Syncer
	Interval time.Duration
}

type Manager struct {
	logger *zap.Logger
	jobs   []ScheduledSyncer
}

func NewManager(logger *zap.Logger, jobs ...ScheduledSyncer) *Manager {
	return &Manager{logger: logger, jobs: jobs}
}

func (m *Manager) Start(ctx context.Context, runOnStartup bool) {
	for _, job := range m.jobs {
		if runOnStartup {
			go m.runOnce(ctx, job.Syncer)
		}
		go m.runTicker(ctx, job)
	}
}

func (m *Manager) runTicker(ctx context.Context, job ScheduledSyncer) {
	if job.Interval <= 0 {
		job.Interval = 10 * time.Minute
	}
	ticker := time.NewTicker(job.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			m.runOnce(ctx, job.Syncer)
		}
	}
}

func (m *Manager) runOnce(ctx context.Context, syncer Syncer) {
	start := time.Now()
	err := syncer.RunOnce(ctx)
	if err != nil {
		m.logger.Warn("syncer failed", zap.String("syncer", syncer.Name()), zap.Error(err), zap.Duration("duration", time.Since(start)))
		return
	}
	m.logger.Info("syncer completed", zap.String("syncer", syncer.Name()), zap.Duration("duration", time.Since(start)))
}
