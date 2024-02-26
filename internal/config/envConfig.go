package config

import (
    "github.com/joho/godotenv"
    "log"
    "os"
)

type Config struct {
    // Configuration Variables 
    DatabaseURL string
    
}

var AppConfig *Config

func init() {
    // Load .env file
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    // Initializing AppConfig with values from .env file
    AppConfig = &Config{
        DatabaseURL: os.Getenv("DATABASE_URL"),
        
    }
}

// GetDatabaseURL returns the configured database URL
func GetDatabaseURL() string {
    return AppConfig.DatabaseURL
}

// Create functions for other env variables
 