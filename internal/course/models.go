package course

import (

	"gorm.io/gorm"
)

type Course struct {

	gorm.Model
	Name        string `gorm:"unique;not null"`
	Description string `gorm:"not null"`
	CategoryID  uint
	Category    Category `gorm:"foreignKey:CategoryID;references:ID"`
	
}

type Category struct {
	gorm.Model
	Name        string `gorm:"unique;not null"`
	Description string `gorm:"not null"`
	Courses []Course
}

