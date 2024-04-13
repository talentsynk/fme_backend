package main

import (
	myuser "fme_backend/internal/user"
	"fme_backend/internal/config"
	"fme_backend/internal/course"
	"fme_backend/internal/mdas"
	"fme_backend/internal/stc"
	"fme_backend/internal/student"
	
)

func init() {
	config.ConnectToDb()
}

func main() {
	config.DB.AutoMigrate(
		// Add all the migration structs here
		&myuser.User{},
		&mda.Mda{},
		&student.Student{},
		&course.Category{},
		&course.Course{},
		&stc.Stc{},
		
	)
}