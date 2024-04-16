package config

import (	
	//   "gorm.io/driver/postgres"
	"gorm.io/driver/mysql"
	   "gorm.io/gorm"

	"log"
)

var DB *gorm.DB


func ConnectToDb() {
	var err error

	 DB, err = gorm.Open(mysql.Open(GetDatabaseURL()), &gorm.Config{}) 
	//   DB, err = gorm.Open(postgres.Open(GetDatabaseURL()), &gorm.Config{}) 


	if err != nil {
		log.Fatal("Error to connect to database")
	}
}

