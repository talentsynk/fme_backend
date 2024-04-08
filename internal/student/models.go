package student

import (
	"time"

	"gorm.io/gorm"
)

type Student struct {
	// complete it
	gorm.Model
	Firstname     string `gorm:"type:varchar(255); not null"`
	Lastname      string `gorm:"type:varchar(255); not null"`
	DOB           time.Time  `gorm:"type:date"`
	StateOfOrigin string `gorm:"type:varchar(255)"`
	StateOfResidence string `gorm:"type:varchar(255)"`
	Gender        string `gorm:"type:varchar(255)"`
	GraduationStatus bool 
	Fmestudent bool
	UserID uint
	MdaID uint
}