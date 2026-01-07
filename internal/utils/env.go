package utils

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"sync"
)

type Env struct {
	PORT        		 string
	STATIC_DIR			 string
	KEA_API_IP   		 string
	KEA_API_URL      string
	KEA_API_USERNAME string
	KEA_API_PASSWORD string
	KEA_DB_HOST      string
	KEA_DB_PORT      int
	KEA_DB_USER      string
	KEA_DB_PASSWORD  string
	KEA_DB_NAME      string
}

var envOnce sync.Once
var env Env

// LoadEnvFile reads a .env file and sets the environment variables
func LoadEnvFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Ignore empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Split the line into key-value pairs
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue // Skip invalid lines
		}

		// Trim spaces from key and value
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Set the environment variable
		if err := os.Setenv(key, value); err != nil {
			return err
		}
	}

	return scanner.Err()
}

// getEnv retrieves an environment variable, returning a default if empty
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// LoadEnv initializes the environment variables
func LoadEnv() {
	envOnce.Do(func() {
		if loadErr := LoadEnvFile(".env"); loadErr != nil {
			if !os.IsNotExist(loadErr) {
				Fatal("Error loading .env file: %v", loadErr)
			}
		}

		// Load environment variables with defaults where necessary
		env.PORT = getEnv("PORT", "8080")
		env.STATIC_DIR = getEnv("STATIC_DIR", "static")
		
		env.KEA_API_IP = os.Getenv("KEA_API_IP")
		env.KEA_API_URL = os.Getenv("KEA_API_URL")
		env.KEA_API_USERNAME = os.Getenv("KEA_API_USERNAME")
		env.KEA_API_PASSWORD = os.Getenv("KEA_API_PASSWORD")

		env.KEA_DB_HOST = getEnv("KEA_DB_HOST", env.KEA_API_IP)
		mySQLPortStr := getEnv("KEA_DB_PORT", "3306")
		mySQLPort, err := strconv.Atoi(mySQLPortStr)
		if err != nil {
			Fatal("Invalid KEA_DB_PORT: %s. Must be an integer.", mySQLPortStr)
		}
		env.KEA_DB_PORT = mySQLPort
		env.KEA_DB_USER = os.Getenv("KEA_DB_USER")
		env.KEA_DB_PASSWORD = os.Getenv("KEA_DB_PASSWORD")
		env.KEA_DB_NAME = os.Getenv("KEA_DB_NAME")
	})
}

func ValidateEnv(e Env) {
  requiredEnvVars := []struct {
    name  string
    value interface{}
  }{
    {"HTTP_PORT", e.PORT},
		{"KEA_API_IP", e.KEA_API_IP},
    {"KEA_API_USERNAME", e.KEA_API_USERNAME},
    {"KEA_API_PASSWORD", e.KEA_API_PASSWORD},
    {"KEA_DB_HOST", e.KEA_DB_HOST},
    {"KEA_DB_PORT", e.KEA_DB_PORT},
    {"KEA_DB_USER", e.KEA_DB_USER},
    {"KEA_DB_PASSWORD", e.KEA_DB_PASSWORD},
    {"KEA_DB_NAME", e.KEA_DB_NAME},
  }

  for _, envVar := range requiredEnvVars {
		switch v := envVar.value.(type) {
		case string:
			if v == "" {
				Fatal("Missing or empty environment variable: %s", envVar.name)
			} else {
				Debug("%s = %s", envVar.name, v)
			}
		case int:
			if v == 0 {
				Fatal("Invalid or missing environment variable: %s must be a valid integer", envVar.name)
			} else {
				Debug("%s = %d", envVar.name, v)
			}
		default:
			Fatal("Unexpected type for environment variable: %s", envVar.name)
		}
	}
}

// GetEnv returns the initialized Env struct
func GetEnv() Env {
	return env
}
