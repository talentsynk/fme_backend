package utilities

import "regexp"

func VeriryNINFormat(nin string) bool {
	re := regexp.MustCompile(`^\d{11}$`)
	return re.MatchString(nin)
}