package controller

import (
	//Inbuilt packages
	"fmt"
	"net/http"
	"strconv"
	"time"
    model "fme_backend/user/models"
    
	//Project packages
	"fme_backend/config"
	"fme_backend/utilities"

	//External packages
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

// Create FME User
func CreateFmeUser(c *gin.Context) {
	//Read the request body and binds it to the schema variable
	fmt.Println("controller started")
	if c.Bind(&model.UserCreateSchema) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read request body",
		})
		return
	}
	

	// Check that the phone number is of the right syntax
	if !utilities.IsNigerianPhoneNumber(model.UserCreateSchema.PhoneNumber){
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read request body",
		})
		return
	}

	fmt.Println(model.UserCreateSchema.PhoneNumber)


	// Hash the password

	hash, err := bcrypt.GenerateFromPassword([]byte(model.UserCreateSchema.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	fmt.Println("hash succesful")
	// setup the user create instance

	user:= model.User{
		
		PhoneNumber: model.UserCreateSchema.PhoneNumber,
		Email: model.UserCreateSchema.Email,
		Password: string(hash),
		OTPExpiresAt: time.Now(),
		IsMda: false,
		IsStc: false,
		IsFme: true,
		IsStudent: false,
		IsAdmin: false,
		ActivityStatus: "active",


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



// Deactivate User
func DeactivateUser(c *gin.Context) {
	// Get User data from authentication token
	userId ,_ := c.Get("userId")

	var user model.User
	result := config.DB.First(&user, userId)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
}
	
	user.ActivityStatus = "inactive"

	result = config.DB.Save(&user)
    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
        return
    }

	c.JSON(200, gin.H{"message": "User Deactivated successfully"})
	
}



//Suspend the user
func SuspendUser(c *gin.Context) {
		// Get User data from authentication token
		userId ,_ := c.Get("userId")

		var user model.User
		result := config.DB.First(&user, userId)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
    }
	
		user.ActivityStatus = "suspended"

		result = config.DB.Save(&user)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
			return
		}

		c.JSON(200, gin.H{"message": "User Suspended successfully"})
	}



// Activate User
	func ActivateUser(c *gin.Context) {
		// Get User data from authentication token
		userId ,_ := c.Get("userId")

		var user model.User
		result := config.DB.First(&user, userId)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
    }
	
		user.ActivityStatus = "active"

		result = config.DB.Save(&user)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
			return
		}

		c.JSON(200, gin.H{"message": "User Suspended successfully"})
	}
	


// Login 
func Login(c *gin.Context) {
	var user model.User

	if c.Bind(&model.LoginSchema) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read request body",
		})
		return
	}

	config.DB.First(&user, "email = ?", model.LoginSchema.Email)
	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Incorrect email",
		})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(model.LoginSchema.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Incorrect Password",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.Email,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(config.GetHashSecret()))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"jwt":        tokenString,
		"token_type": "Bearer",
		"expires_in": 86400, //time in seconds
	})
}



// Request Otp 
func RequestOtp(c *gin.Context) {
		// Get User data from authentication token
		userId ,_ := c.Get("userId")

		var user model.User
		result := config.DB.First(&user, userId)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
    }
	// set the otp and expiry time 
	user.OTPExpiresAt = time.Now().Add(time.Minute *15)
	user.OTP =  strconv.Itoa(utilities.GenerateOtp())

	//Save the user
	result = config.DB.Save(&user)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
			return
		}
	
		c.JSON(200, gin.H{"message": "User otp request successful"})
	
	// mail the otp to the user - the whole process must be a transaction
}



// change password
func ChangePassword(c *gin.Context) {
	// make sure that the otp matches and has not expired

	// hash and store the new password
}
