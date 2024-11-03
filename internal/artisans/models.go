package artisans

import (
	myuser "fme_backend/internal/user"

	"gorm.io/gorm"
)




type Artisans struct {
	gorm.Model
	FirstName string
	LastName string
	StateOfResidence string
	LGA string
	BusinessName string
	BusinessDescription string
	UserID uint
	Skill string
	User myuser.User		`gorm:"foreignKey:UserID;references:ID"`

}