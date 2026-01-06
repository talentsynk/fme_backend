package config

import (
	"os"

	"github.com/joho/godotenv"
)

var _ = godotenv.Load("sandbox.env")

type Config struct {
	DatabaseURL string
	HashSecret  string
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
	return "Pass123*"
}

func GetFmeEmail() string {
	return "fme.testing@gmail.com"
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
