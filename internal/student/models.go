package student

import (
	"fme_backend/internal/course"
	myuser "fme_backend/internal/user"
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
	PhoneNumber	     string 		`gorm:"unique"`
	Fmestudent        bool
	IsDisabled 	      bool
	DisabilityName    string			`gorm:"type:varchar(255)"`
	GraduationDate    time.Time  `gorm:"type:date"`
	UserID            uint
	MdaID             uint
	StcID             uint
	NationalIdentityNumber string	`gorm:"not null; default:12345678912"`
	LocalGovernment	string	`gorm:"not null; default:'kosofe'"`
	User     myuser.User 	`gorm:"foreignKey:UserID;references:ID"`

}

type StudentCourse struct{
	gorm.Model
	CourseID uint	`gorm:"not null"`
	StudentID uint	`gorm:"not null"`
	IsCertified	bool
	Course    course.Course `gorm:"foreignKey:CourseID;references:ID"`
	Student Student		`gorm:"foreignKey:StudentID;references:ID"`
}