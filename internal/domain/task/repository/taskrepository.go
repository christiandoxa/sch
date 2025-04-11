// Package repository implements data access logic.
package repository

import (
	"context"
	"log"
	"sch/internal/domain/task/model"
	"sync"
)

// TaskRepository defines repository operations for Task entities.
type TaskRepository interface {
	GetPendingTasks(ctx context.Context) ([]model.Task, error)
	MarkTaskProcessed(ctx context.Context, taskID int) error
}

// InMemoryTaskRepository is an in-memory implementation of TaskRepository.
type InMemoryTaskRepository struct {
	tasks []model.Task
	mu    sync.Mutex
}

// NewInMemoryTaskRepository creates a new in-memory repository instance.
func NewInMemoryTaskRepository() *InMemoryTaskRepository {
	// Initialize with sample tasks.
	tasks := []model.Task{
		{ID: 1, Name: "Task1", Data: "Data for Task 1"},
		{ID: 2, Name: "Task2", Data: "Data for Task 2"},
		{ID: 3, Name: "Task3", Data: "Data for Task 3"},
	}
	return &InMemoryTaskRepository{tasks: tasks}
}

// GetPendingTasks returns a list of pending tasks.
func (repo *InMemoryTaskRepository) GetPendingTasks(_ context.Context) ([]model.Task, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	// For demonstration, return all tasks.
	return repo.tasks, nil
}

// MarkTaskProcessed marks a task as processed by removing it from the store.
func (repo *InMemoryTaskRepository) MarkTaskProcessed(_ context.Context, taskID int) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	for i, t := range repo.tasks {
		if t.ID == taskID {
			repo.tasks = append(repo.tasks[:i], repo.tasks[i+1:]...)
			log.Printf("Task with ID %d marked as processed.", taskID)
			return nil
		}
	}
	return nil
}
