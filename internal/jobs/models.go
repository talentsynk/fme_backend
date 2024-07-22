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
	Status            bool    
	EmployerID        uint 

}

type JobApplication struct {
	gorm.Model
    JobID      uint   
    StudentID  uint   
}


type SaveJob struct {
	gorm.Model
    JobID      uint   
    StudentID  uint 
}

type JobRecommendation struct {
	 gorm.Model
	 StudentID            uint     
	 RecommendationText   string   `gorm:"type:varchar(255);not null"`
     
	}

type CompletedJobs struct{
	 gorm.Model
	 JobID          uint
	 StudentID      uint  
	 EmployerID     uint
	 RecommendationText   string   `gorm:"type:varchar(255);not null"`
}


type ArtisanEmployed struct{
	gorm.Model
	StudentID  uint
	EmployerID uint
	RecommendationText   string   `gorm:"type:varchar(255);not null"` 

}