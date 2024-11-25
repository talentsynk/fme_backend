package config

import (
	    //   "gorm.io/driver/postgres"
	     "log"
          "gorm.io/driver/mysql"
	    "gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDb() {
	var err error
	        DB, err = gorm.Open(mysql.Open(GetDatabaseURL()), &gorm.Config{})
	        //  DB, err = gorm.Open(postgres.Open("postgresql://fme_db_backend_qom2_user:V3FxYabrqqMjMDQPklrxIeiwdtHeLBa7@dpg-cr5mnb2j1k6c739762hg-a.oregon-postgres.render.com/fme_db_backend_qom2"), &gorm.Config{})
         	if err != nil {
	     	log.Fatal("Error to connect to database")
    	}
}
