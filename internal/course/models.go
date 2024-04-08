package course

import "gorm.io/gorm"

type Course struct {

	gorm.Model
	Name string `gorm:"unique;not null"`
	Description string `gorm:"not null"`
	SectorID uint


}

type Sector struct {
	gorm.Model
	Name string `gorm:"unique;not null"`
	Description string `gorm:"not null"`
	Courses Course
}