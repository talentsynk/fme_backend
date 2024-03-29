package utilities

import (
	"fmt"
	"math/rand"
	"time"
)


func GenerateOtp() string {
    // Seed the random number generator with the current time
    rand.Seed(time.Now().UnixNano())

    // Generate a random 4-digit number
    otp := rand.Intn(1000000)
    otpString := fmt.Sprintf("%05d", otp)

    return otpString
}
