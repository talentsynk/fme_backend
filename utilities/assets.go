package utilities

import "regexp"



func IsNigerianPhoneNumber(phoneNumber string) bool {
    // Define a regular expression pattern for Nigerian phone numbers
    // Nigerian phone numbers start with either 080, 081, 090, 070, or 091 followed by 8 digits
    pattern := `^(080|081|090|070|091)\d{8}$`

    // Compile the regular expression pattern
    regex := regexp.MustCompile(pattern)

    // Use the MatchString method to check if the phoneNumber matches the pattern
    return regex.MatchString(phoneNumber)
}