package stc

import "gorm.io/gorm"

type Stc struct {
	//complete it
	gorm.Model
	ID uint
	Ownership                    string  `gorm:"unique;not null"`
	CentreCode                   string  `gorm:"unique;not null"`
	Name                         string  `gorm:"unique;not null"`
	LocalGovernment              string  `gorm:"unique;not null"`
	State                        string  `gorm:"unique;not null"`
	isOperational                bool    `gorm:"unique;not null"`
	CertificateOfOperationURL    string `gorm:"unique;not null"`
	MdaID                        uint
	UserID                       uint	
}