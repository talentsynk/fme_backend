package mda

import (
	"gorm.io/gorm"
)

type Mda struct {
	gorm.Model

	ID         uint
	Name       string `gorm:"unique;not null"`
	AgencyCode string `gorm:"unique;not null"`
	IsActive   bool   `gorm:"unique;not null"`
	UserID     uint
}