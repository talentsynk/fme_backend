package config

import (
	          "log"
	        "os"
              "github.com/joho/godotenv"
)

type Config struct {
	// Configuration Variables
	DatabaseURL string
	HashSecret  string
}

var AppConfig *Config

func init() {
	//   Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	//   Initializing AppConfig with values from .env file
	AppConfig = &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		HashSecret:  os.Getenv("HASH_SECRET"),
	}
}

func GetDatabaseURL() string {
	return os.Getenv("DATABASE_URL")
}

func GetHashSecret() string {
	return os.Getenv("HASH_SECRET")
}

func GetResendSecret() string {
    return os.Getenv("RESEND_SECRET")
}

func GetEnvType() string {
	return os.Getenv("ENV_TYPE")
}

func GetFmePassWord() string {
	return os.Getenv("FME_PASSWORD")
}

func GetFmeEmail() string {
	return os.Getenv("FME_EMAIL")
}

func GetSendMailAcctToken() string {
	return os.Getenv("PM_ACCT_TOKEN")
}

func GetSendMailServerToken() string {
	return os.Getenv("PM_SERVER_TOKEN")
}

func GetHomeMail() string {
	return os.Getenv("HOME_MAIL")
}