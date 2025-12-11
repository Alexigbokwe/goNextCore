package scheduler

import (
	"github.com/Alexigbokwe/gonext-framework/core/logger"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

// Scheduler manages recurring jobs
type Scheduler struct {
	cron *cron.Cron
}

// NewScheduler creates a new scheduler instance
func NewScheduler() *Scheduler {
	// WithSeconds allows 6-part cron spec (second, minute, hour, dom, month, dow)
	c := cron.New(cron.WithSeconds())
	return &Scheduler{
		cron: c,
	}
}

// Add registers a job with a cron spec
// Spec format: "Second Minute Hour Dom Month Dow"
// Example: "0 * * * * *" (Every minute)
func (s *Scheduler) Add(spec string, cmd func()) (cron.EntryID, error) {
	id, err := s.cron.AddFunc(spec, func() {
		// Recover from panic to prevent crashing scheduler
		defer func() {
			if r := recover(); r != nil {
				logger.Log.Error("Panic in scheduled job", zap.Any("reason", r))
			}
		}()
		cmd()
	})

	if err != nil {
		logger.Log.Error("Failed to add scheduled job", zap.String("spec", spec), zap.Error(err))
		return 0, err
	}

	logger.Log.Info("Scheduled job added", zap.String("spec", spec), zap.Int("id", int(id)))
	return id, nil
}

// Start begins the scheduler
func (s *Scheduler) Start() {
	s.cron.Start()
	logger.Log.Info("Scheduler started")
}

// Stop stops the scheduler
func (s *Scheduler) Stop() {
	s.cron.Stop()
	logger.Log.Info("Scheduler stopped")
}
