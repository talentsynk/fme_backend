package mda

import (
	"gorm.io/gorm"
)

type Mda struct {
	gorm.Model
	RegisterName       string `gorm:"not null"`
	Address            string `gorm:"not null"`
	StateOfOperation   string `gorm:"not null"`
	UserID             uint
}