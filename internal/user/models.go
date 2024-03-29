package myuser

import (
	"fme_backend/internal/mda"
	"fme_backend/internal/stc"
	"fme_backend/internal/student"
	"time"

	"gorm.io/gorm"
)

type User struct {
    gorm.Model
    // think about about using a single role column and the effect on the
    
    PhoneNumber    string         `gorm:"unique;not null"`
    Email          string         `gorm:"unique;not null"`
    Password       string         `gorm:"not null"`
    OTP            string         
    OTPExpiresAt   time.Time 
    OTPVerified bool 
    Role int `gorm:"not null"`
    IsActive bool `gorm:"not null"`
    Mdas mda.Mda
    Stcs stc.Stc
    Students student.Student
    

    
}
