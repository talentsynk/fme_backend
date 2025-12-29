package myuser

import (
	// "fme_backend/internal/interfaces"
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email         string `gorm:"unique;not null"`
	Password      string `gorm:"not null"`
	OTP           string
	OTPExpiresAt  time.Time
	OTPVerified   bool
	Role          int  `gorm:"not null"`
	IsActive      bool `gorm:"not null"`
	SuspendReason string
}
