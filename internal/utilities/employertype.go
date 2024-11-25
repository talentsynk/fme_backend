package utilities

import (
	"strings"
)

const (
	Individual       string = "individual"
	Organisation       string = "Organisation"
)

type EmployerType struct{}

func ValidateEmployerType(employerType string)(string, bool){
	employerType = strings.ToLower(employerType)
	switch employerType{
	case Individual:
		return Individual, true
	case Organisation:
		return FullTime, true
	default:
		return "", false
	}
}