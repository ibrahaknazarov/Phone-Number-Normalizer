// Package config handles configuration loading from environment variables.
package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config represents the application configuration.
type Config struct {
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
}

// Load reads configuration from environment variables with sensible defaults.
func Load() *Config {
	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnvInt("DB_PORT", 5432),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "phone"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
}

// PSQLInfo returns a PostgreSQL connection string for the admin database.
func (c *Config) PSQLInfo() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBSSLMode)
}

// PSQLInfoWithDB returns a PostgreSQL connection string with a specific database.
func (c *Config) PSQLInfoWithDB(dbName string) string {
	return fmt.Sprintf("%s dbname=%s", c.PSQLInfo(), dbName)
}

// getEnv retrieves an environment variable with a default fallback.
func getEnv(key, defaultVal string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultVal
}

// getEnvInt retrieves an integer environment variable with a default fallback.
func getEnvInt(key string, defaultVal int) int {
	valStr := os.Getenv(key)
	if val, err := strconv.Atoi(valStr); err == nil {
		return val
	}
	return defaultVal
}
