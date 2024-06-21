package config

import (
	// "log"
	"os"
    // "github.com/joho/godotenv"
)

// type Config struct {
// 	// Configuration Variables
// 	DatabaseURL string
// 	HashSecret  string
// }

// var AppConfig *Config

// func init() {
// 	//   Load .env file
// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Fatal("Error loading .env file")
// 	}

// 	//   Initializing AppConfig with values from .env file
// 	AppConfig = &Config{
// 		DatabaseURL: os.Getenv("DATABASE_URL"),
// 		HashSecret:  os.Getenv("HASH_SECRET"),
// 	}
// }

func GetDatabaseURL() string {
	return os.Getenv("DATABASE_URL")
}

func GetHashSecret() string {
	return os.Getenv("HASH_SECRET")
}

func GetResendSecret() string {
    return os.Getenv("RESEND_SECRET")
}
