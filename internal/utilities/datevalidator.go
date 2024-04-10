package utilities

import "time"


func ParseDoB(dobStr string) (time.Time, error) {
	// date format MM/DD/YYYY
	format := "01/02/2006"
  
	// Parse the DoB string into a time.Time instance
	t, err := time.Parse(format, dobStr)
  
	// Check for parsing errors
	if err != nil {
	  return time.Time{}, err
	}
  
	// Return the parsed time.Time instance
	return t, nil
  }