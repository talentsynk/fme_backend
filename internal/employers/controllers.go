package employer

import (
	"fme_backend/internal/config"
	myuser "fme_backend/internal/user"
	"fme_backend/internal/utilities"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)


func CreateEmployer(c *gin.Context) {
    // Bind and validate request body
    if err := c.ShouldBindJSON(&CreateEmployerSchema); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "message": "Invalid request body",
        })
        return
    }

    // State validation
    var State string
    if CreateEmployerSchema.State != "" {
        var result bool
        State, result = utilities.ValidateState(CreateEmployerSchema.State)
        if !result {
            c.JSON(http.StatusBadRequest, gin.H{
                "message": "Incorrect state of origin",
            })
            return
        }
    }

    // Phone number validation
    if !utilities.IsNigerianPhoneNumber(CreateEmployerSchema.PhoneNumber) {
        c.JSON(http.StatusBadRequest, gin.H{
            "message": "Incorrect phone number",
        })
        return
    }

    // Encrypt password before storing
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(CreateEmployerSchema.Password), 10)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "message": "Failed to encrypt password",
        })
        return
    }

    // Start database transaction
    tx := config.DB.Begin()

    // Create employer user
    result, message, newUserID := myuser.CreateEmployerUser(tx, CreateEmployerSchema.Email, string(hashedPassword))
    if !result {
        tx.Rollback()
        c.JSON(http.StatusBadRequest, gin.H{
            "message": message,
        })
        return
    }

    // Create employer record
    employer := Employer{
        FirstName:   CreateEmployerSchema.FirstName,
        LastName:    CreateEmployerSchema.LastName,
        Email:       CreateEmployerSchema.Email,
        PhoneNumber: CreateEmployerSchema.PhoneNumber,
        NIN:         CreateEmployerSchema.NIN,
        Password:    string(hashedPassword), // Store the hashed password
        State:       State,
        LGA:         CreateEmployerSchema.LGA,
        UserID:      newUserID,
    }

    if err := tx.Create(&employer).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusBadRequest, gin.H{
            "message": "Failed to create Employer",
        })
        return
    }

    // Commit transaction
    tx.Commit()

    // Send success response
    c.JSON(http.StatusOK, gin.H{
        "message": "Employer created successfully",
    })

    fmt.Println(employer)
}


func GetEmployer(c *gin.Context) {
	employerIDstr, exists := c.Get("employerID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized user",
		})
		return
	}

	employerID, ok := employerIDstr.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized user failed to convert to uint"})
		return
	}

	var employer GetEmployerSchema

	result := config.DB.Table("employers").
		Select("id, first_name, last_name, email, phone_number, nin, state, lga, user_id").
		First(&employer, employerID)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Employer not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"employer": employer})
}


func GetEmployerByID(c *gin.Context){
	employerIDParam := c.Param("id")
     
	employerID, err := strconv.Atoi(employerIDParam)
     if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":"Invalid employer ID",
		})
		return 
	 }

	 var employer GetEmployerSchema
      result := config.DB.Table("employers").
	  Select("id, first_name, last_name, email, phone_number, nin, state, lga, user_id").
	  Where("id = ?", employerID).
	  First(&employer)

  if result.Error != nil{
	c.JSON(http.StatusNotFound, gin.H{
		"error":"Employer not found",
	})

	return 
  }

  c.JSON(http.StatusOK, gin.H{"employer":employer})

}
