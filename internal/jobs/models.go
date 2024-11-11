package jobs

import (
	"gorm.io/gorm"
)


type Job  struct {
	gorm.Model
    JobTitle          string   `gorm:"type:varchar(255);not null"`
	Location          string   `gorm:"type:varchar(255);not null"`
	Budget            float64   
	JobType           string   `gorm:"type:varchar(255);not null"`
	Category          string   `gorm:"type:varchar(255);not null"`
	Description       string   `gorm:"not null"`
	Requirement       string   `gorm:"not null"`
	Responsibilities  string   `gorm:"not null"`
	HiringStatus      bool 
	Skills            string
	Status	          string
	EmployerID        uint 
}

type JobApplication struct {
	gorm.Model
    JobID             uint   
    ArtisanID         uint  
	Ratings           string
	ApplicationStatus string 
}

type JobApplicationRating struct {
	gorm.Model
	Rating            uint
	JobApplicationID  uint
	Description	      string
}


type SaveJob struct {
	gorm.Model
    JobID      uint   
    ArtisanID  uint 
}

type EmployerJobRating struct {
	gorm.Model
    JobID       uint   
	Ratings     uint
	Description string 
}

