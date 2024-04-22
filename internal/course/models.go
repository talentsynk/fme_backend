package course

import "gorm.io/gorm"

type Course struct {

	gorm.Model
	Name        string `gorm:"unique;not null"`
	Description string `gorm:"not null"`
	CategoryID  uint
}

type Category struct {
	gorm.Model
	Name        string `gorm:"unique;not null"`
	Description string `gorm:"not null"`
	Courses     Course
}

type StudentCourse struct{
	gorm.Model
	CourseID uint	`gorm:"not null"`
	StudentID uint	`gorm:"not null"`
	IsCertified	bool
}

type MdaCourse struct{
	gorm.Model
	CourseID uint	`gorm:"not null"`
	MdaID uint	`gorm:"not null"`
	// IsCertified	bool
}

type StcCourse struct{
	gorm.Model
	CourseID uint	`gorm:"not null"`
	StcID uint	`gorm:"not null"`
	// IsCertified	bool
}