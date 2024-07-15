package utilities

import (
	"strings"
)

const (
	OnHire       string = "on-hire"
	ContractJob  string = "contract-job"
)

type JobType struct{}

func ValidateJobType(jobType string)(string, bool){
	jobType = strings.ToLower(jobType)
	switch jobType{
	case OnHire:
		return OnHire, true
	case ContractJob:
		return ContractJob, true
	default:
		return "", false
	}
}