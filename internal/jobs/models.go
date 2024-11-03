package jobs

import (
	"fme_backend/internal/artisans"
	employer "fme_backend/internal/employers"

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
	HiringStatus            bool 
	Skills    string
	Status	string
	EmployerID        uint 
	Emoloyer employer.Employer			`gorm:"foreignKey:EmployerID;references:ID"`
}

type JobApplication struct {
	gorm.Model
    JobID      uint   
    ArtisanID  uint  
	ApplicationStatus string
	Job	Job `gorm:"foreignKey:JobID;references:ID"` 
	Artisan artisans.Artisans		`gorm:"foreignKey:ArtisanID;references:ID"`
}

type JobApplicationRating struct {
	gorm.Model
	Rating uint
	JobApplicationID  uint
	Description	string
	JobApplication JobApplication `gorm:"foreignKey:JobApplicationID;references:ID"`
}


type SaveJob struct {
	gorm.Model
    JobID      uint   
    ArtisanID  uint 
	Job	Job `gorm:"foreignKey:JobID;references:ID"` 
	Artisan artisans.Artisans		`gorm:"foreignKey:ArtisanID;references:ID"`
}

type EmployerJobRating struct {
	gorm.Model
    JobID      uint   
	Ratings uint
	Description string 
	Job	Job `gorm:"foreignKey:JobID;references:ID"` 
}

