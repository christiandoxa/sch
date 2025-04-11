// Package service contains business logic for processing tasks.
package service

import (
	"context"
	"log"
	"sch/internal/domain/task/model"
	"sch/internal/domain/task/repository"
	"sch/tool"
	"sync"
)

// TaskService defines methods to process tasks.
type TaskService interface {
	ProcessTasks(ctx context.Context) error
}

// TaskServiceImpl implements TaskService.
type TaskServiceImpl struct {
	repo repository.TaskRepository
}

// NewTaskService returns a new instance of TaskService.
func NewTaskService(repo repository.TaskRepository) TaskService {
	return &TaskServiceImpl{repo: repo}
}

// ProcessTasks retrieves pending tasks from the repository and processes them concurrently.
// It uses a semaphore (channel) to limit the number of parallel workers.
func (s *TaskServiceImpl) ProcessTasks(ctx context.Context) error {
	tasks, err := s.repo.GetPendingTasks(ctx)
	if err != nil {
		log.Printf("Error fetching tasks: %v", err)
		return err
	}
	log.Printf("Found %d tasks to process.", len(tasks))

	var wg sync.WaitGroup
	maxWorkers := tool.GetMaxParallelism()
	sem := make(chan struct{}, maxWorkers)

	for _, task := range tasks {
		wg.Add(1)
		sem <- struct{}{} // Acquire a slot.
		go func(t model.Task) {
			defer wg.Done()
			defer func() { <-sem }() // Release the slot.

			// Process the task (insert business logic here).
			log.Printf("Processing task ID: %d, Name: %s", t.ID, t.Name)
			// Simulate work (e.g., call external API, perform computation, etc.).
			tool.SimulateWork(2) // Simulate 2 seconds of work.

			// Mark the task as processed in the repository.
			if err := s.repo.MarkTaskProcessed(ctx, t.ID); err != nil {
				log.Printf("Error marking task %d as processed: %v", t.ID, err)
			}
		}(task)
	}
	wg.Wait()
	log.Println("All tasks processed.")
	return nil
}
