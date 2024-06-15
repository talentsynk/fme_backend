package main

import (
	"fme_backend/internal/config"
	"fme_backend/internal/course"
	employer "fme_backend/internal/employers"
	"fme_backend/internal/mdas"
	"fme_backend/internal/stc"
	"fme_backend/internal/student"
	myuser "fme_backend/internal/user"
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
		&course.StudentCourse{},
		&course.MdaCourse{},
		&course.StcCourse{},
		&employer.Employer{},
	)
}