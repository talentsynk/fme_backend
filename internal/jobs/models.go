package jobs

import (
   "gorm.io/gorm"
)


type Job  struct {
	gorm.Model
    JobTitle          string   `gorm:"type:varchar(255);not null"`
	Location          string   `gorm:"type:varchar(255);not null"`
	Budget            string   `gorm:"type:varchar(255);not null"`
	JobType           string   `gorm:"type:varchar(255);not null"`
	Category          string   `gorm:"type:varchar(255);not null"`
	Description       string   `gorm:"not null"`
	Requirement       string   `gorm:"not null"`
	Responsibilities  string   `gorm:"not null"`
	Status            bool     `gorm:"not null"`
	EmployerID        uint
    
}

type JobApplication struct {
	gorm.Model
    JobID      uint   `gorm:"not null"`
    StudentID uint   `gorm:"not null"`
}


type SaveJob struct {
	gorm.Model
    JobID      uint   `gorm:"not null"`
    StudentID uint   `gorm:"not null"`
}

type JobRecommendation struct {
	 gorm.Model
	 StudentID uint   `gorm:"not null"`
	 RecommendationText   string   `gorm:"type:varchar(255);not null"`
}

type CompletedJobs struct{
	gorm.Model
	JobID      uint   `gorm:"not null"`
}