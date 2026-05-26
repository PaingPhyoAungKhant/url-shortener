// Package config provides functionality to load and manage application configuration from environment variables.
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds the configruation for the application
type Config struct {
	App    AppConfig
	Server ServerConfig
}

type AppConfig struct {
	Env       string
	LogLevel  string
	LogFormat string
}

type ServerConfig struct {
	Port            int
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
}

// Load reads the configuration from the environment variables and returns a Config struct.
func Load() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	config := &Config{
		App: AppConfig{
			Env:       getEnv("APP_ENV", "development"),
			LogLevel:  getEnv("APP_LOGLEVEL", "INFO"),
			LogFormat: getEnv("APP_LOGFORMAT", "json"),
		},
		Server: ServerConfig{
			Port:            getEnvAsInt("SERVER_PORT", 8081),
			ReadTimeout:     getEnvAsDuration("SERVER_READ_TIMEOUT", 15*time.Second),
			WriteTimeout:    getEnvAsDuration("SERVER_WRITE_TIMEOUT", 15*time.Second),
			IdleTimeout:     getEnvAsDuration("SERVER_IDLE_TIMEOUT", 60*time.Second),
			ShutdownTimeout: getEnvAsDuration("SERVER_SHUTDOWN_TIMEOUT", 10*time.Second),
		},
	}
	return config, nil
}

func (c *Config) Validate() error {
	switch c.App.Env {
	case "development", "production", "test":
	default:
		return fmt.Errorf("app.env %q: want development|production|test", c.App.Env)
	}

	if c.Server.Port < 1 || c.Server.Port > 65535 {
		return fmt.Errorf("server.port %d: want 1-65535", c.Server.Port)
	}

	switch c.App.LogLevel {
	case "debug", "info", "warn", "error":
	default:
		return fmt.Errorf("app.loglevel %q: want debug|info|warn|error", c.App.LogLevel)
	}

	switch c.App.LogFormat {
	case "json", "text":
	default:
		return fmt.Errorf("app.logformat %q: want json|text", c.App.LogFormat)
	}

	return nil
}

// getEnv retrieves the value of the environment variable name by the key.
// if the variable is not set, it return the defaule value provided as the second argument.
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getEnv retrieves the value of the environment variable name by the key and returns it as an integer.
// if the variable is not set, it return the defaule value provided as the second argument.
func getEnvAsInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intValue
}

// getEnv retrieves the value of the environment variable name by the key and returns it as time.Duration.
// if the variable is not set, it return the defaule value provided as the second argument.
func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	timeValue, err := time.ParseDuration(value)
	if err != nil {
		return defaultValue
	}
	return timeValue
}
