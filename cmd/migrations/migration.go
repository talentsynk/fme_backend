package main

import (
	"fme_backend/internal/config"
	"fme_backend/internal/user"
	"fme_backend/internal/mdas"
)

func init() {
	config.ConnectToDb()
}

func main() {
	config.DB.AutoMigrate(
		// Add all the migration structs here
		&myuser.User{},
		&mda.Mda{},
		
	)
}