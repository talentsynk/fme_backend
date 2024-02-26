package myuser

import (
	"fmt"
	"net/http"
	"time"

	"fme_backend/internal/config"
	"fme_backend/internal/utilities"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)


func CreateFme(c *gin.Context) {
	//Read the request body and binds it to the schema variable
	fmt.Println("coltroller started")
	if c.Bind(&UserCreateSchema) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read request body",
		})
		return
	}
	fmt.Println("binded succesfully")

	// Check that the phone number is of the right syntax
	if !utilities.IsNigerianPhoneNumber(UserCreateSchema.PhoneNumber){
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read request body",
		})
		return
	}

	fmt.Println(UserCreateSchema.PhoneNumber)


	// Hash the password

	hash, err := bcrypt.GenerateFromPassword([]byte(UserCreateSchema.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	fmt.Println("hash succesful")
	// setup the user create instance

	user:= User{
		FirstName: UserCreateSchema.FirstName,
		LastName: UserCreateSchema.LastName,
		PhoneNumber: UserCreateSchema.PhoneNumber,
		Email: UserCreateSchema.Email,
		Password: string(hash),
		OTPExpiresAt: time.Now(),
	}

	fmt.Println(user)

	result := config.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	c.JSON(200, gin.H{"message": "User created successfully"})
}

