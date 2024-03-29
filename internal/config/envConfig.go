package config

import (
    "github.com/joho/godotenv"
    "log"
    "os"
)

type Config struct {
    // Configuration Variables 
    DatabaseURL string
    HashSecret string
    
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
        HashSecret: os.Getenv("HASH_SECRET"),
        
    }
}

// GetDatabaseURL returns the configured database URL
func GetDatabaseURL() string {
    return AppConfig.DatabaseURL
}

func GetHashSecret() string {
    return AppConfig.HashSecret
}

// Create functions for other env variables
 