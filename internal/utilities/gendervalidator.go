package utilities

import "strings"

func ValidateGender(gender string) (string, bool){
	gender = strings.ToLower(gender)
	switch gender {
	case "male":
		return "male", true
	case "female":
		return "female", true
	default:
		return "", false
	}
}