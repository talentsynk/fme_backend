package mda

import "gorm.io/gorm"


type Mda struct {
	//Complete it
	gorm.Model
	Name string `gorm:"unique;not null"`
	AgencyCode string `gorm:"unique;not null"`
	UserID uint 

}