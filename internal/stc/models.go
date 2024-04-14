package stc

import "gorm.io/gorm"

type Stc struct {
	//complete it
	gorm.Model
	Ownership                    string  `gorm:"unique;not null"`
	CentreCode                   string  `gorm:"unique;not null"`
	Name                         string  `gorm:"not null"`
	LocalGovernment              string  `gorm:"not null"`
	State                        string  `gorm:"not null"`
	isOperational                bool   
	CertificateOfOperationURL    string `gorm:"unique;not null"`
	Fmestc                       bool
	MdaID                        uint
	UserID                       uint	
}