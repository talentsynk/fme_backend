package myuser

import (
	//Inbuilt packages
	"net/http"
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
		IsMda: false,
		IsStc: false,
		IsFme: true,
		IsStudent: false,
		IsAdmin: false,
		ActivityStatus: "active"
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
		userId ,userexists := c.Get("userID")
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
  
// 		var user User
// 		result := config.DB.First(&user, userId)
// 		if result.Error != nil {
// 			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
// 			return
//     }
	
// 		user.ActivityStatus = "active"
//      	result = config.DB.Save(&user)
// 		if result.Error != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
// 			return
// 		}

	// get the instance  
// 	var instance User
// 	instance_result := config.DB.First(&instance, id)
// 	if instance_result.Error != nil {
// 		c.JSON(http.StatusNotFound, gin.H{
// 			"message": "instance does not exist",
// 		})
// 		return
// 	}

// 	// get the user and confirm permission
// 	userId ,userexists := c.Get("userID")
// 	if !userexists {
// 		c.JSON(http.StatusUnauthorized, gin.H{
// 			"message": "Problem with the aithorization token",
// 		})
// 		return
// 	}

// 	var user User
// 	user_result := config.DB.First(&user, userId)
// 	if user_result.Error != nil {
// 		c.JSON(http.StatusUnauthorized, gin.H{
// 			"message": "Problem with authorization token",
// 		})
// 		return
// 	}


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
	// bind request schema 
	if c.Bind(&RequestOtpSchema) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Incorrect request schema",
		})
		return
	}


	// check if email is registered and fetch user data
	var user User
	config.DB.First(&user, "email = ?", RequestOtpSchema.Email)
	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Email is not in the record",
		})
		return
	}

	if !user.IsActive {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "This account has been suspended",
		})
		return
	}

	// generate otp and expiry time 
	user.OTP = utilities.GenerateOtp()
	user.OTPExpiresAt = time.Now().Add(time.Minute*3)
	result:= config.DB.Save(&user)
	if result.Error !=nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Unable to update the user record",
		})
		return
	}


	// send otp -remember this must be changed to mail
	c.JSON(http.StatusOK, gin.H{
		"message": "Otp generated succesfully",
		"otp": user.OTP,
	})
}

// change password
func VerifyOtp(c *gin.Context) {
	// bind the schema
	if c.Bind(&VerifyOtpSchema) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Incorrect request schema",
		})
		return
	}

	// get the user and verify otp
	var user User
	config.DB.First(&user, "email = ?", VerifyOtpSchema.Email)
	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Email is not in the record",
		})
		return
	}

	if user.OTP != VerifyOtpSchema.Otp {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Incorrect Otp",
		})
		return
	}

	user.OTPVerified = true
	result:= config.DB.Save(&user)
	if result.Error !=nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Unable to update the user record",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Otp verified succesfully",
	})

}


func ChangePassword(c *gin.Context) {
	// bind the schema
	if c.Bind(&ChangePasswordSchema) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Incorrect request schema",
		})
		return
	}

	// get the user and verify otp
	var user User
	config.DB.First(&user, "email = ?", ChangePasswordSchema.Email)
	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Email is not in the record",
		})
		return
	}

	//check otp verification
	if !user.OTPVerified || time.Now().After(user.OTPExpiresAt) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Bad otp",
		})
		return
	}


	// change and hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(ChangePasswordSchema.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	user.Password = string(hash)
	user.OTPVerified = false
	user.OTPExpiresAt = time.Now()

	result:= config.DB.Save(&user)
	if result.Error !=nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Unable to update the user record",
		})
		return
	}

	// return success
	c.JSON(http.StatusOK, gin.H{
		"message": "Password changed succesfully",
	})
}
