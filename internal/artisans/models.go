package artisans

import "gorm.io/gorm"




type Artisans struct {
	gorm.Model
	FirstName string
	LastName string
	StateOfResidence string
	LGA string
	BusinessName string
	BusinessDescription string
	StudentID uint 
	UserID uint
	Skill string


}