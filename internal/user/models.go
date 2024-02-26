package myuser

import (
    "gorm.io/gorm"
    "time"
)

type User struct {
    gorm.Model
    FirstName      string         `gorm:"not null"`
    LastName       string         `gorm:"not null"`
    PhoneNumber    string         `gorm:"unique;not null"`
    Email          string         `gorm:"unique;not null"`
    Password       string         `gorm:"not null"`
    Address        string         `gorm:"not null"`
    OTP            string         
    OTPExpiresAt   time.Time      
    
}
