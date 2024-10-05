package etp

import (
	"log/slog"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all the configuration variables from .env
type Config struct {
	DatabaseURL      string
	ServerPort       int
	JWTAccessSecret  string
	JWTRefreshSecret string
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) LoadConfig() error {
	slog.Info("Loading .env file")
	if err := godotenv.Load(); err != nil {
		slog.Error("Error loading .env file: ", "error", err)
		return err
	}

	c.DatabaseURL = getEnv("DATABASE_URL", "postgres://postgres:password@localhost:5432/etp?sslmode=disable")
	c.ServerPort = getEnvAsInt("SERVER_PORT", 1234)
	c.JWTAccessSecret = getEnv("JWT_ACCESS_SECRET", "secret")
	c.JWTRefreshSecret = getEnv("JWT_REFRESH_SECRET", "secret")

	return nil
}

// Helper function to get environment variables as string with a fallback
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// Helper function to get environment variables as an integer with a fallback
func getEnvAsInt(name string, fallback int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return fallback
}

//
// // Helper function to get environment variables as a boolean with a fallback
// func getEnvAsBool(name string, fallback bool) bool {
// 	valueStr := getEnv(name, "")
// 	if value, err := strconv.ParseBool(valueStr); err == nil {
// 		return value
// 	}
// 	return fallback
// }
//
// // Helper function to get environment variables as a time.Duration with a fallback
// func getEnvAsDuration(name string, fallback string) time.Duration {
// 	valueStr := getEnv(name, fallback)
// 	duration, err := time.ParseDuration(valueStr)
// 	if err != nil {
// 		return time.Duration(30) * time.Second // Default to 30 seconds
// 	}
// 	return duration
// }
