// Package usecase orchestrates high-level operations by coordinating the services.
package usecase

import (
	"context"
	"log"
	"sch/internal/domain/task/service"
	"time"
)

// TaskUseCase defines use case methods related to task processing.
type TaskUseCase interface {
	Execute(ctx context.Context)
}

// TaskUseCaseImpl is an implementation of TaskUseCase.
type TaskUseCaseImpl struct {
	taskService service.TaskService
}

// NewTaskUseCase returns a new TaskUseCase instance.
func NewTaskUseCase(taskService service.TaskService) TaskUseCase {
	return &TaskUseCaseImpl{taskService: taskService}
}

// Execute starts the task processing workflow and logs timing information.
func (uc *TaskUseCaseImpl) Execute(ctx context.Context) {
	start := time.Now()
	log.Println("Task UseCase Execution Started")
	if err := uc.taskService.ProcessTasks(ctx); err != nil {
		log.Printf("Error executing task use case: %v", err)
	}
	duration := time.Since(start)
	log.Printf("Task UseCase Execution Completed in %v", duration)
}
