// main.go initializes the configuration, repository, service, use case, and scheduler.
// It registers the task use case in the scheduler and waits for a graceful shutdown.
package main

import (
	"context"
	"log"
	"sch/config"
	"sch/internal/domain/task/repository"
	"sch/internal/domain/task/service"
	"sch/scheduler"
	"sch/usecase"
)

func main() {
	// Load application configuration.
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize repository, service, and use case.
	taskRepo := repository.NewInMemoryTaskRepository()
	taskService := service.NewTaskService(taskRepo)
	taskUseCase := usecase.NewTaskUseCase(taskService)

	// Initialize the scheduler.
	sch := scheduler.NewScheduler()

	// Register a scheduler task that executes the task use case.
	sch.RegisterTask(config.AppConfig.Scheduler.EngineName, config.AppConfig.Scheduler.Interval, func(ctx context.Context) {
		taskUseCase.Execute(ctx)
	})

	// Start the scheduler.
	sch.Start()
	log.Println("Scheduler started.")

	// Keep the application running until an interrupt is received.
	sch.WaitForShutdown()

	log.Println("Application stopped gracefully.")
}
