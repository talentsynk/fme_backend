package student

import (
	"time"

	"gorm.io/gorm"
)

type Student struct {
	// complete it
	gorm.Model
	Firstname        string     `gorm:"type:varchar(255); not null"`
	Lastname         string     `gorm:"type:varchar(255); not null"`
	DOB              time.Time  `gorm:"type:date; not null"`
	StateOfOrigin    string     `gorm:"type:varchar(255)"`
	StateOfResidence string  `gorm:"type:varchar(255); not null"`
	Gender           string   `gorm:"type:varchar(255)"`
	SID				 string   `gorm:"type:varchar(255)"`
	NsqLevel		 string   `gorm:"type:varchar(255)"`
	Address		     string       `gorm:"type:varchar(255)"`
	GraduationStatus bool 
	PhoneNumber	     string 		  `gorm:"unique"`
	Fmestudent       bool
	UserID           uint
	MdaID            uint
	StcID            uint
}