package stc

import "gorm.io/gorm"

type Stc struct {
	//complete it
	gorm.Model
	
	Name                         string  `gorm:"not null"`
	Email                        string  `gorm:"unique;not null"`
	Address                      string  `gorm:"not null"`
	State                        string  `gorm:"not null"`
	isOperational                bool   
	CertificateOfOperationURL    string `gorm:"not null"`
	Fmestc                       bool
	MdaID                        uint
	UserID                       uint	
}
