package stc

import "gorm.io/gorm"

type Stc struct {
	//complete it
	gorm.Model
	UserID uint
	MdaID uint
}