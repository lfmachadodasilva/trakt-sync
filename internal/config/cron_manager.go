package config

import (
	"context"
	"sync"

	"github.com/robfig/cron/v3"
)

type CronManager struct {
	cron      *cron.Cron
	entryID   cron.EntryID
	frequency string
	lock      sync.Mutex
}

// NewCronManager creates a new instance of CronManager
func NewCronManager() *CronManager {
	return &CronManager{
		cron: cron.New(),
	}
}

// Start starts the cron job with the given frequency and job function
func (cm *CronManager) Start(ctx context.Context, frequency string, job func()) error {
	cm.lock.Lock()
	defer cm.lock.Unlock()

	// Stop the existing cron job if running
	if cm.entryID != 0 {
		cm.cron.Remove(cm.entryID)
	}

	// Add the new cron job
	entryID, err := cm.cron.AddFunc(frequency, job)
	if err != nil {
		return err
	}

	cm.entryID = entryID
	cm.frequency = frequency
	cm.cron.Start()

	return nil
}

// Stop stops the currently running cron job
func (cm *CronManager) Stop() {
	cm.lock.Lock()
	defer cm.lock.Unlock()

	if cm.entryID != 0 {
		cm.cron.Remove(cm.entryID)
		cm.entryID = 0
	}
	cm.cron.Stop()
}

// UpdateFrequency updates the cron job frequency dynamically
func (cm *CronManager) UpdateFrequency(ctx *context.Context, newFrequency string, job func()) error {
	if cm.frequency == newFrequency {
		return nil // No change needed
	}
	return cm.Start(*ctx, newFrequency, job)
}
