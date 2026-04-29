package application

import (
	"time"

	"github.com/techies/orders-api/helpers"
)

// Config holds all application configuration
type Config struct {
	Port            int
	JWTSecret       string
	DatabaseURL     string
	IdleTimeout     time.Duration
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
	RateLimit       int
}

// DefaultConfig returns the default configuration loaded from environment
func DefaultConfig() *Config {
	return &Config{
		Port:            helpers.GetEnvInt("PORT", 3000),
		JWTSecret:       helpers.GetEnv("JWT_SECRET", "wF7zR9p8Yq2Vt6kLmH3uN1bX5eG0aQjZ"),
		DatabaseURL:     "./orders.db?parseTime=true",
		IdleTimeout:     helpers.GetEnvDuration("IDLE_TIMEOUT", time.Minute),
		ReadTimeout:     helpers.GetEnvDuration("READ_TIMEOUT", 10*time.Second),
		WriteTimeout:    helpers.GetEnvDuration("WRITE_TIMEOUT", 30*time.Second),
		ShutdownTimeout: helpers.GetEnvDuration("SHUTDOWN_TIMEOUT", 5*time.Second),
		RateLimit:       100,
	}
}
