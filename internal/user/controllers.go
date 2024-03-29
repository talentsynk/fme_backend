package myuser

import (
	//Inbuilt packages
	"fmt"
	"net/http"
	"strconv"
	"time"

	//Project packages
	"fme_backend/internal/config"
	"fme_backend/internal/utilities"

	//External packages
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

// Create FME User
func CreateFmeUser(c *gin.Context) {
	//Read the request body and binds it to the schema variable
	fmt.Println("controller started")
	if c.Bind(&UserCreateSchema) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read request body",
		})
		return
	}
	

	// Check that the phone number is of the right syntax
	if !utilities.IsNigerianPhoneNumber(UserCreateSchema.PhoneNumber){
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read request body",
		})
		return
	}

	// Hash the password

	hash, err := bcrypt.GenerateFromPassword([]byte(UserCreateSchema.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	// setup the user create instance
	user:= User{
		
		PhoneNumber: UserCreateSchema.PhoneNumber,
		Email: UserCreateSchema.Email,
		Password: string(hash),
		OTPExpiresAt: time.Now(),
		Role:1,
		IsActive: true,
	}

	result := config.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	c.JSON(200, gin.H{"message": "User created successfully"})
}



//Suspend the user
func SuspendUser(c *gin.Context) {
		// Receive the part parameter.
		id:= c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "path parameter not provided",
			})
			return
		}

		// get the instance  
		var instance User
		instance_result := config.DB.First(&instance, id)
		if instance_result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "instance does not exist",
			})
			return
		}

		// get the user and confirm permission
		userId ,userexists := c.Get("userId")
		if !userexists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Problem with the authorization token",
			})
		}

		var user User
		user_result := config.DB.First(&user, userId)
		if user_result.Error != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Problem with authorization token",
			})
			return
		}


		if !CanSuspendActivate(&user,&instance) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Only fmes can suspend",
			})
			return
		}

		// suspend the user
		instance.IsActive = false
		result:= config.DB.Save(&instance)
		if result.Error !=nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Unable to update the user record",
			})
			return
		}

		c.JSON(200, gin.H{"message": "User suspended successfully"})
	}


// Activate User
func ActivateUser(c *gin.Context) {
		// Receive the part parameter.
	id:= c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "path parameter not provided",
		})
		return
	}

	// get the instance  
	var instance User
	instance_result := config.DB.First(&instance, id)
	if instance_result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "instance does not exist",
		})
		return
	}

	// get the user and confirm permission
	userId ,userexists := c.Get("userId")
	if !userexists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Problem with the aithorization token",
		})
	}

	var user User
	user_result := config.DB.First(&user, userId)
	if user_result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Problem with authorization token",
		})
		return
	}

	if !CanSuspendActivate(&user,&instance) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Only fmes can suspend",
		})
		return
	}

	// activate the user
	instance.IsActive = true
	result:= config.DB.Save(&instance)
	if result.Error !=nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Unable to update the user record",
		})
		return
	}

	c.JSON(200, gin.H{"message": "User activated successfully"})
	}
	


// Login 
func Login(c *gin.Context) {
	var user User
	// receive the request body
	if c.Bind(&LoginSchema) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Incorrect request schema",
		})
		return
	}

	// get the user record
	config.DB.Where("email= ?", LoginSchema.Email).First(&user)
	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Incorrect email details",
		})
		return
	}

	// check if the user is active
	if !user.IsActive {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Suspended User",
		})
		return
	}

	// compare the passwords
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(LoginSchema.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Incorrect login details",
		})
		return
	}

	// generate the jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.Email,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(config.GetHashSecret()))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Unable to authenticate this user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"jwt":        tokenString,
		"token_type": "Bearer",
		"expires_in": 86400, //time in seconds
		"message": "succesful login",
		"role":user.Role,
	})
}



// Request Otp 
func RequestOtp(c *gin.Context) {
		// Get User data from authentication token
		userId ,_ := c.Get("userId")

		var user User
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
