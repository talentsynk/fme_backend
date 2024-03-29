package main

import (
	"fme_backend/internal/config"
	"fme_backend/internal/mda"
	"fme_backend/internal/user"
)

func init() {
	config.ConnectToDb()
}

func main() {
	config.DB.AutoMigrate(
		// Add all the migration structs here
		&myuser.User{},&mda.Mda{},
	)
}