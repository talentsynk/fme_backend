package student

import (
	"time"

	"gorm.io/gorm"
)

type Student struct {
	// complete it
	gorm.Model
	Firstname     string    `gorm:"type:varchar(255); not null"`
	Lastname      string     `gorm:"type:varchar(255); not null"`
	DOB           time.Time  `gorm:"type:date; not null"`
	StateOfOrigin string     `gorm:"type:varchar(255); not null"`
	StateOfResidence string   `gorm:"type:varchar(255); not null"`
	Gender          string      `gorm:"type:varchar(255)"`
	GraduationStatus bool 
	Fmestudent        bool
	UserID uint
	MdaID uint
	StcID uint
}