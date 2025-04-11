// Package config provides configuration loading from environment variables.
package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds the configuration settings for the application.
type Config struct {
	Scheduler SchedulerConfig
	DB        DatabaseConfig
}

// SchedulerConfig holds scheduler-specific configurations.
type SchedulerConfig struct {
	EngineName string
	Interval   time.Duration
}

// DatabaseConfig holds database connection settings.
type DatabaseConfig struct {
	DSN string
}

// AppConfig is the global configuration instance.
var AppConfig Config

// LoadConfig loads configuration values from environment variables.
func LoadConfig() error {
	if err := godotenv.Load(); err != nil {
		return err
	}

	// Load scheduler settings.
	AppConfig.Scheduler.EngineName = getEnv("SCHEDULER_ENGINE_NAME", "MyScheduler")
	intervalSec := getEnvAsInt("SCHEDULER_INTERVAL_SECONDS", 10)
	AppConfig.Scheduler.Interval = time.Duration(intervalSec) * time.Second

	// Load database settings (for demonstration purposes only).
	AppConfig.DB.DSN = getEnv("DATABASE_DSN", "user:pass@tcp(localhost:3306)/dbname")

	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	log.Printf("Environment variable %s not set, using default: %s", key, defaultValue)
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	log.Printf("Environment variable %s not set properly, using default: %d", key, defaultValue)
	return defaultValue
}
