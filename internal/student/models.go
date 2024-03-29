package student

import "gorm.io/gorm"

type Student struct {
	// complete it
	gorm.Model
	UserID uint
	MdaID uint
}