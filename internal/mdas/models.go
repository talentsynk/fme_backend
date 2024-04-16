package mda

import (
	"gorm.io/gorm"
)

type Mda struct {
	gorm.Model
	RegisterName       string  `gorm:"not null"`
	Address            string  `gorm:"not null"`
	StateOfOperation   string  `gorm:"not null"`
	Email              string  `gorm:"unique;not null"`
	UserID             uint
	IsActive           bool    `gorm:"not null"`
}

