package mda

import (
	"gorm.io/gorm"
)

type Mda struct {
	gorm.Model
	RegisterName       string `gorm:"unique;not null"`
	Email              string `gorm:"unique;not null"`
	Address            string `gorm:"unique;not null"`
	StateOfOperation   string `gorm:"unique;not null"`
	IsActive           bool   `gorm:"unique;not null"`
	UserID             uint
}