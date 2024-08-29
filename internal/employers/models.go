package employer


import (
	
	"gorm.io/gorm"
)




type Employer  struct {
    gorm.Model

    FirstName   string `gorm:"type:varchar(255);not null"`
	LastName    string `gorm:"type:varchar(255);not null"`
	PhoneNumber string `gorm:"unique;not null"`
	NIN         string `gorm:"unique;not null"`
	State       string `gorm:"type:varchar(255);not null"`
	LGA         string `gorm:"type:varchar(255);not null"`
	UserID      uint
}