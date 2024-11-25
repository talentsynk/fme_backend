package employer

import (

	"gorm.io/gorm"
	myuser "fme_backend/internal/user"

)




type Employer  struct {
    gorm.Model

    FirstName   string `gorm:"type:varchar(255);not null"`
	LastName    string `gorm:"type:varchar(255);not null"`
	PhoneNumber string `gorm:"unique;not null"`
	NIN         string `gorm:"unique;not null"`
	State       string `gorm:"type:varchar(255);not null"`
	LGA         string `gorm:"type:varchar(255);not null"`
	IsCompany 	bool	`gorm:"default:false"`
	CompanyName	string	`gorm:"type:varchar(255)"`
	CompanyCAC	string	`gorm:"type:varchar(255)"`
	UserID      uint
	User     myuser.User 	`gorm:"foreignKey:UserID;references:ID"`
}

