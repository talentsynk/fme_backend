package jobs

import (
	"time"
)

var CreateJobSchema struct {
    JobTitle         string 
	Location         string
	JobType          string
	Budget           string
	Category         string 
	Description      string
	Requirement      string
	Responsibilities string
}

type GetAllJobsSchema struct{
	Id               uint
	FirstName        string 
	LastName         string 
    Location         string
    Description      string 
    JobType          string 
	Budget           string 
    JobTitle         string 
    Requirement      string
	Responsibilities string
    CreatedAt        time.Time
	Status           bool 
	EmployerID       uint  
 
}

type GetJobSchema struct {
	Id                 uint
	FirstName          string 
	LastName           string 
    Location           string
    Description        string 
	Budget             string
    JobType            string  
    JobTitle           string 
    Requirement        string
	Responsibilities   string
    CreatedAt          time.Time
	Status             bool 
	EmployerID         uint  

}




