package mda

import (
	"gorm.io/gorm"
)

type Mda struct {
	gorm.Model
	RegisterName       string `gorm:"unique;not null"`
	Address            string `gorm:"unique;not null"`
	StateOfOperation   string `gorm:"not null"`
	UserID             uint
}