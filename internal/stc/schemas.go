package stc 

import (
	"time"
)

var StcCreateSchema  struct {
	Name                         string 
	Email                        string 
	Address                      string  
	State                        string
	PhoneNumber                  string
}

type GetAllStcSchema struct{
	Id          uint
	StateOfOperation string
	CreatedAt    time.Time
	Name        string
	Address      string
	IsActive     bool   `json:"is_active"`
	StudentCount uint    `json:"student_count"`
	CourseCount     uint
	UserId          uint 
	Email       string
	CertifiedStudentCount  uint
	NonCertifedStudentCount uint
}


type GetStcSchema struct{
	Id          uint
	Name        string
	Address      string
	IsActive     bool   `json:"is_active"`
	StudentCount uint    `json:"student_count"`
	CertifiedStudentCount   uint
	NonCertifedStudentCount     uint
	CourseCount     uint
	UserId          uint
}