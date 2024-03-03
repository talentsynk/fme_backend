package utilities

import (
	"math/rand"
	"time"
)


func GenerateOtp() int {
    // Seed the random number generator with the current time
    rand.Seed(time.Now().UnixNano())

    // Generate a random 5-digit number
    return rand.Intn(90000) + 10000
}
