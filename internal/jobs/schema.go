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
	Id            uint
	FirstName     string
	LastName      string 
    CreatedAt     time.Time
    JobType       string 
    Location      string 
    Budget        string 
    Description   string 
    EmployerID    uint  
}

type GetJobSchema struct {
	Id               uint
	FirstName        string 
	LastName         string 
    Location         string
    Description      string 
    JobType          string  
    JobTitle           string 
    Requirement      string
	Responsibilities string
    CreatedAt        time.Time
	EmployerID       uint  

}

