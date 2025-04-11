// Package scheduler provides a simplified abstraction over gocron
// to register and manage periodic tasks.
package scheduler

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-co-op/gocron/v2"
)

// TaskFunc defines the function signature for scheduler tasks.
type TaskFunc func(ctx context.Context)

// Scheduler wraps the gocron.Scheduler.
type Scheduler struct {
	cron gocron.Scheduler
}

// NewScheduler creates a new Scheduler instance using the local timezone.
func NewScheduler() *Scheduler {
	s, err := gocron.NewScheduler(gocron.WithLimitConcurrentJobs(1, gocron.LimitModeReschedule))
	if err != nil {
		log.Fatal(err)
	}
	return &Scheduler{cron: s}
}

// RegisterTask registers a new task with a given descriptive name and interval.
// The task function receives a background context for cancellation or timeout.
func (s *Scheduler) RegisterTask(name string, interval time.Duration, task TaskFunc) {
	_, err := s.cron.NewJob(
		gocron.DurationJob(interval),
		gocron.NewTask(func() {
			ctx := context.Background()
			log.Printf("Task [%s] started", name)
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Task [%s] encountered panic: %v", name, r)
				}
			}()
			task(ctx)
			log.Printf("Task [%s] completed", name)
		}),
		gocron.WithSingletonMode(gocron.LimitModeReschedule),
	)
	if err != nil {
		log.Printf("Failed to register task [%s]: %v", name, err)
	}
}

// Start begins executing the scheduled tasks asynchronously.
func (s *Scheduler) Start() {
	s.cron.Start()
}

// Stop halts the scheduler.
func (s *Scheduler) Stop() {
	err := s.cron.StopJobs()
	if err != nil {
		log.Printf("Failed to stop scheduler: %v", err)
	}
}

// WaitForShutdown listens for OS interrupt signals (Ctrl+C, SIGTERM)
// and gracefully stops the scheduler.
func (s *Scheduler) WaitForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Gracefully shutting down scheduler...")
	s.Stop()
}
