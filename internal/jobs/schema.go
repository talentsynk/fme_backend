package jobs

import (
	"time"
)

var CreateJobSchema struct {
    JobTitle         string 
	Location         string
	JobType          string
	Budget           float64
	Category         string 
	Description      string
	Requirement      string
	Responsibilities string
}

type GetJobByIdsSchema struct{
	Id uint
	JobTitle string
	EmployerId uint
	EmployerFirstName string
	EmployerLastName string
	Skills string
	CreatedAt time.Time
	Location string
	JobType string
	Requirements string
	Responsibilities string
	Description string
	Status string
	HiringStatus bool
}

type GetJobSchema struct {
	Id uint
	JobTitle string
	Description string
	Amount string
	JobType string
	Status string

}


type GetSavedJobSchema struct{
	Id	uint
	Name 	string
	Description 	string
	Amount 		float64
	JobType 	string
	Location 	string
}

type GetAppliedJobSchema struct {
	Id	uint
	Name 	string
	Description 	string
	ApplicationStatus 	string

}

var HireArtisanSchema struct {
	JobId 		uint
	ArtisanId 	uint
}


type GetApplicantSchema struct {
	ApplicationId uint
	ArtisanId uint
	BusinessName string
	BusinessDescription string
	JobApplicationDate  time.Time
	AverageRating float64
	FirstName string
	LastName string
	ApplicationStatus string


}

var CompleteJobSchema struct {
	Rating uint
	Description string
}

type JobFilterSchema struct {
	Status string		`form:"status"`
	MinBudget float64	`form:"min_budget"`
	MaxBudget float64	`form:"max_budget"`
	JobType string		`form:"job_type"`
	DaysAgo uint			`form:"days_ago"`
}





