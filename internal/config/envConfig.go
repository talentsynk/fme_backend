package config

import "os"

func GetDatabaseURL() string {
	return os.Getenv("DATABASE_URL")
}

func GetHashSecret() string {
	return os.Getenv("HASH_SECRET")
}
