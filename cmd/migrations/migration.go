package cmd

import (
	"fme_backend/config"
	models  "fme_backend/user/models"
)

func init() {
	config.ConnectToDb()
}

func main() {
	config.DB.AutoMigrate(
		// Add all the migration structs here
		&models.User{},
	)
}