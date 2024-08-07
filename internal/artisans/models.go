package artisans

import "gorm.io/gorm"




type Artisans struct {
	gorm.Model


	StudentID uint 
	UserID uint
	Skill string


}