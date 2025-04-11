// Package model defines the core business entities.
package model

// Task represents a simple domain entity for scheduled tasks.
type Task struct {
	ID   int
	Name string
	Data string
}
