package utilities

import (
    "crypto/rand"
    "encoding/base64"
    "strings"
)


func GenerateRandomPassword() string {
    // Generate a random byte slice
    randomBytes := make([]byte, 8)
    _, err := rand.Read(randomBytes)
    if err != nil {
        // Handle error
        return ""
    }

    // Encode the byte slice to base64 string
    randomPassword := base64.URLEncoding.EncodeToString(randomBytes)

    // Remove any non-alphanumeric characters and convert to lowercase
    randomPassword = strings.ToLower(randomPassword)
    randomPassword = strings.Map(func(r rune) rune {
        if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
            return r
        }
        return -1
    }, randomPassword)

    // Trim the password to 8 characters
    if len(randomPassword) > 8 {
        randomPassword = randomPassword[:8]
    }

    return randomPassword
}

// Function to create user with automatic password generation
// func CreateUserWithAutomaticPassword(tx *gorm.DB, email, phoneNumber string) (bool, string, uint) {
//     // Generate random password
//     randomPassword := generateRandomPassword()

//     // Check if password generation failed
//     if randomPassword == "" {
//         return false, "Failed to generate password", 0
//     }

//     // Use CreateMdaUser function to create user with generated password
//     return CreateMdaUser(tx, phoneNumber, email, randomPassword)
// }



