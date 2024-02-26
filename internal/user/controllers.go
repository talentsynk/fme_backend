package myuser

import (
	"fme_backend/internal/config"
	"fme_backend/internal/utilities"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)



// @Summary Create a new user
// @Description Create a new user with the provided details
// @Tags users
// @Accept json
// @Produce json
// @Param body body UserCreateSchema true "User details"
// @Success 200 {object} User
// @Failure 400 {object} ErrorResponse
// @Router /users [post]
func CreateFme(c * gin.Context) {
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
		FirstName: UserCreateSchema.FirstName,
		LastName: UserCreateSchema.PhoneNumber,
		PhoneNumber: UserCreateSchema.PhoneNumber,
		Email: UserCreateSchema.Email,
		Password: string(hash),

	}

	result := config.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create user",
		})
		return
	}
}

