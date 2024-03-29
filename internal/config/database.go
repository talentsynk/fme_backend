package config

import (
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"log"
)

  var DB *gorm.DB

  func ConnectToDb() {
	  var err error
	  DB, err = gorm.Open(mysql.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})
	
	  if err != nil {
		  log.Fatal("Error to connect to database")
	  }
  }