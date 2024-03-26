package myuser

import (
    "gorm.io/gorm"
    "time"
)

type User struct {
    gorm.Model
    
    PhoneNumber    string         `gorm:"unique;not null"`
    Email          string         `gorm:"unique;not null"`
    Password       string         `gorm:"not null"`
    OTP            string         
    OTPExpiresAt   time.Time   
    IsMda bool `gorm:"not null"`
    IsStc bool  `gorm:"not null"`
    IsFme bool `gorm:"not null"`
    IsStudent bool `gorm:"not null"`
    IsAdmin bool `gorm:"not null"`
    ActivityStatus string `gorm:"not null"`    
}
