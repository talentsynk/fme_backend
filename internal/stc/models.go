package stc

import "gorm.io/gorm"

type Stc struct {
	//complete it
	gorm.Model
	Ownership                    string  `gorm:"unique;not null"`
	CentreCode                   string  `gorm:"unique;not null"`
	Name                         string  `gorm:"unique;not null"`
	LocalGovernment              string  `gorm:"unique;not null"`
	State                        string  `gorm:"unique;not null"`
	isOperational                bool   
	CertificateOfOperationURL    string `gorm:"unique;not null"`
	Fmestc                       bool
	MdaID                        uint
	UserID                       uint	
}