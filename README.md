# Task Scheduler Application

This project is a simple task scheduling application built with Go. It demonstrates how to structure a Go project using
domain-driven design principles, handle concurrency with goroutines and semaphores, and schedule periodic work using
the [gocron](https://github.com/go-co-op/gocron) library. The application simulates processing a set of tasks stored in
an in-memory repository.

## Table of Contents

- [Overview](#overview)
- [Project Structure](#project-structure)
- [Components](#components)
    - [Model](#model)
    - [Repository](#repository)
    - [Service](#service)
    - [Use Case](#use-case)
    - [Scheduler](#scheduler)
    - [Tool](#tool)
- [Configuration](#configuration)
- [Installation](#installation)
- [Running the Application](#running-the-application)
- [Graceful Shutdown](#graceful-shutdown)
- [License](#license)

## Overview

The application orchestrates a simple workflow for processing tasks:

1. **Task Definition**: The domain entity `Task` is defined in the `model` package.
2. **Task Storage**: An in-memory repository (`InMemoryTaskRepository`) holds the tasks and provides methods to retrieve
   pending tasks and mark them as processed.
3. **Task Processing**: The `TaskService` (in the `service` package) retrieves tasks from the repository and processes
   them concurrently. Concurrency is controlled by a semaphore (implemented with a Go channel) to limit the number of
   parallel workers.
4. **Use Case Execution**: The `TaskUseCase` orchestrates the high-level operation by calling the task service and
   logging execution details.
5. **Scheduling**: The scheduler (in the `scheduler` package) abstracts over the gocron library to periodically execute
   the task use case. It supports starting the schedule asynchronously and gracefully shutting down on interrupt
   signals.
6. **Utility Functions**: The `tool` package provides helper functions such as calculating the maximum parallelism based
   on the system’s available processors and simulating work by pausing execution.

## Project Structure

```
.
├── config           # Configuration loading (includes .env files)
├── internal
│   └── domain
│       └── task
│           ├── model       # Defines Task structure (entity)
│           ├── repository  # In-memory repository for tasks
│           └── service     # Business logic for task processing
├── usecase          # Orchestration of task processing workflow
├── scheduler        # Task scheduling abstraction over gocron
├── tool             # Utility functions (e.g., parallelism and work simulation)
└── cmd
    └── app
        └── main.go  # Application entry point: configuration, setup, and scheduler start
```

## Components

### Model

- **Location**: `internal/domain/task/model`
- **Purpose**: Defines the `Task` struct which represents the core domain entity with attributes like `ID`, `Name`, and
  `Data`.

### Repository

- **Location**: `internal/domain/task/repository`
- **Purpose**: Implements data access logic. The `InMemoryTaskRepository` holds a slice of tasks in memory, provides a
  method to retrieve pending tasks, and marks tasks as processed (by removing them).

### Service

- **Location**: `internal/domain/task/service`
- **Purpose**: Contains business logic for processing tasks. The `TaskService` retrieves pending tasks concurrently
  using a semaphore to limit the number of parallel workers. It calls `tool.SimulateWork` to mimic doing some actual
  processing.

### Use Case

- **Location**: `usecase`
- **Purpose**: Acts as an orchestration layer that initiates task processing via the TaskService and logs the execution
  time.

### Scheduler

- **Location**: `scheduler`
- **Purpose**: Provides an abstraction over the gocron library for scheduling tasks. It allows registering a periodic
  function (task use case) and handles graceful shutdown when an interrupt signal is detected.

### Tool

- **Location**: `tool`
- **Purpose**: Offers utility functions including:
    - `GetMaxParallelism`: Determines the level of concurrency based on system capabilities.
    - `SimulateWork`: Simulates work by pausing execution (using `time.Sleep`).

## Configuration

The application uses a configuration file typically managed via environment variables. You need to create a `.env` file
based on the provided `.env.example` file. This configuration might include details such as the scheduler engine name
and the interval at which the tasks should run.

### .env.example

```env
# Scheduler configuration
SCHEDULER_ENGINE_NAME=default
SCHEDULER_INTERVAL_SECONDS=10
```

- **SCHEDULER_ENGINE_NAME**: A descriptive name for the scheduler.
- **SCHEDULER_INTERVAL_SECONDS**: The interval (e.g., "10" for 10 seconds) at which the task use case should be
  executed.

Ensure that you copy this file to a new file named `.env` and adjust the settings as necessary for your environment.

## Installation

1. **Prerequisites**:
    - [Go](https://golang.org/dl/) (version 1.16 or later recommended)
    - Git

2. **Clone the Repository**:

   ```bash
   git clone https://github.com/christiandoxa/sch.git
   cd sch
   ```

3. **Download Dependencies**:

   Use Go modules to manage dependencies:

   ```bash
   go mod tidy
   ```

4. **Configure Environment Variables**:

   Copy the `.env.example` file to `.env` and adjust variables as needed:

   ```bash
   cp .env.example .env
   ```

## Running the Application

To run the application, execute the following command from the project root:

```bash
go run cmd/app/main.go
```

During runtime, the application will:

- Load the configuration from the `.env` file.
- Initialize the in-memory task repository with sample tasks.
- Set up the task service and use case to process the tasks.
- Register a scheduler task with the specified interval.
- Start the scheduler, which periodically executes the task processing workflow.
- Log details of task processing, including acquisition and release of semaphore slots for concurrency control.

## Graceful Shutdown

The scheduler listens for OS interrupt signals (e.g., Ctrl+C or SIGTERM). When such a signal is received, the scheduler
gracefully stops all running jobs and logs a shutdown message:

```plaintext
Gracefully shutting down scheduler...
Application stopped gracefully.
```

This ensures that any tasks in progress can complete or exit safely before the application terminates.

## License

This project is distributed under the MIT License. See the [LICENSE](LICENSE) file for more details.