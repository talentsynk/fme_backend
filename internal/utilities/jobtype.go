package utilities

import (
	"strings"
)

const (
	PartTime       string = "part-time"
	FullTime       string = "full-time"
)

type JobType struct{}

func ValidateJobType(jobType string)(string, bool){
	jobType = strings.ToLower(jobType)
	switch jobType{
	case PartTime:
		return PartTime, true
	case FullTime:
		return FullTime, true
	default:
		return "", false
	}
}