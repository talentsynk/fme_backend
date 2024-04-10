package utilities

func ValidateGender(gender string) (string, bool){
	switch gender {
	case "male":
		return "male", true
	case "female":
		return "female", true
	default:
		return "", false
	}
}