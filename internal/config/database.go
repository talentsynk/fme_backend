package config

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"

	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDb() {
	var err error
	DB, err = gorm.Open(postgres.Open(GetDatabaseURL()), &gorm.Config{})
	if err != nil {
		log.Fatal("Error to connect to database")
	}
	fmt.Println("successfully connected to db.")
}
