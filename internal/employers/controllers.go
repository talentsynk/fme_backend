package employer

import (
	"fme_backend/internal/config"
	 myuser "fme_backend/internal/user"
	"fme_backend/internal/utilities"
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
        PhoneNumber: CreateEmployerSchema.PhoneNumber,
        NIN:         CreateEmployerSchema.NIN,
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

}



func GetAllEmployer(c *gin.Context){
	var employers []GetEmployerSchema

	if result := config.DB.Table("employers").
	Select("id, first_name, last_name, email, phone_number, nin, state, lga, user_id").
	Find(&employers); result.Error != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()});
	     return
	}

     c.JSON(http.StatusOK, gin.H{"employers":employers}) 
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

	var result struct {
		ID          uint
		FirstName   string
		LastName    string
		Email       string
		PhoneNumber string
	}
	
	err := config.DB.Table("employers").
		Select("employers.id, employers.first_name, employers.last_name, users.email, employers.phone_number").
		Joins("JOIN users ON users.id = employers.user_id").
		Where("employers.id = ?", employerID).
		First(&result).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Employer not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"employer": result})
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

	 var result struct {
		ID          uint
		FirstName   string
		LastName    string
		Email       string
		PhoneNumber string
	}
	
	err = config.DB.Table("employers").
		Select("employers.id, employers.first_name, employers.last_name, users.email, employers.phone_number").
		Joins("JOIN users ON users.id = employers.user_id").
		Where("employers.id = ?", employerID).
		First(&result).Error

	if err != nil{
		c.JSON(http.StatusNotFound, gin.H{
			"error":"Employer not found",
		})

	return 
  }

  c.JSON(http.StatusOK, gin.H{"employer":result})

}

func EmployerDashboard(c *gin.Context){
	employerIDstr, exist := c.Get("employerID")
	if !exist{
		c.JSON(http.StatusUnauthorized, gin.H{
			"message":"Problem with the authorization token",
		})

		return
	}
	employerID, ok := employerIDstr.(uint)
	if !ok{
		c.JSON(http.StatusUnauthorized, gin.H{
			"message":"Problem with the authorization token",
		})
		return
	}

	var totalJobPosted int64
	var totalArtisanEmployed int64
	var totalCompletedJobs int64

		// Count the total number of jobs posted by the employer
		err := config.DB.Table("jobs").
				Where("employer_id = ?", employerID).
				Count(&totalJobPosted).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error in total jobs",})
			return
		}
	
		// Count the total completed jobs
		err = config.DB.Table("jobs").
		Where("employer_id = ? AND status = ?", employerID, "completed").
				Count(&totalCompletedJobs).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error in completed jobs",})
			return
		}


		// Count the total number of artisans employed
		err = config.DB.Table("job_applications").
						Joins("JOIN jobs ON jobs.id = job_applications.job_id").
						Where("jobs.employer_id = ? AND job_applications.application_status = ?", employerID, "hired").
						Select("COUNT(DISTINCT job_applications.artisan_id)").
						Count(&totalArtisanEmployed).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error in total employed artisans",})
			return
		}


		c.JSON(http.StatusOK, gin.H{
			"total_job_posted":        totalJobPosted,
			"total_artisan_employed":  totalArtisanEmployed,
			"total_jobs_completed":    totalCompletedJobs,
		})
}

func GetEmployerProfileStats(c *gin.Context) {
	// get the employer id string
	employerIDstr, exist := c.Get("employerID")
	if !exist{
		c.JSON(http.StatusUnauthorized, gin.H{
			"message":"Problem with the authorization token",
		})

		return
	}
	employerID, ok := employerIDstr.(uint)
	if !ok{
		c.JSON(http.StatusUnauthorized, gin.H{
			"message":"Problem with the authorization token",
		})
		return
	}

	// needed variables
	var totalJobPosted int64
	var averageRating float64
	var totalJobCompleted int64
	var totalReccomendations int64

	// get the total jobs posted
	err := config.DB.Table("jobs").
				Where("employer_id = ?", employerID).
				Count(&totalJobPosted).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error in total jobs",})
			return
		}

	// get the average rating 
	// get the rating 
	err = config.DB.Table("employers").
		Select("COALESCE(AVG(employer_job_ratings.ratings), 0) as average_rating").
		Joins("LEFT JOIN jobs ON jobs.employer_id = employers.id").
		Joins("LEFT JOIN employer_job_ratings ON employer_job_ratings.job_id = jobs.id").
		Where("employers.id = ?", employerID).
		Group("employers.id").
		Scan(&averageRating).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error", "details": err.Error()})
	    return
	}

	//total completed jobs
	err = config.DB.Table("jobs").
		Where("employer_id = ? AND status = ?", employerID, "completed").
				Count(&totalJobCompleted).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error in completed jobs",})
			return
		}

	// total reccomendations
	if err := config.DB.Table("employer_job_ratings").
						Joins("JOIN jobs ON jobs.id = employer_job_ratings.job_id").
						Where("jobs.id = ?", employerID).
						Count(&totalReccomendations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error", "details": err.Error()})
	    return
	}

	// return success
	c.JSON(http.StatusOK, gin.H{
		"rating":averageRating,
		"total_recommendations":totalReccomendations,
		"total_jobs_completed":totalJobCompleted,
		"total_jobs_posted":totalJobPosted,
	})
}


func GetEmployerProfileStatsByArtisan(c *gin.Context) {
	// get the employer id
	employerIDParam := c.Param("id")
     
	employerID, err := strconv.Atoi(employerIDParam)
     if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":"Invalid employer ID",
		})
		return 
	 }

	// needed variables
	var totalJobPosted int64
	var averageRating float64
	var totalJobCompleted int64
	var totalReccomendations int64

	// get the total jobs posted
	err = config.DB.Table("jobs").
				Where("employer_id = ?", employerID).
				Count(&totalJobPosted).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error in total jobs",})
			return
		}

	// get the average rating 
	// get the rating 
	err = config.DB.Table("employers").
		Select("COALESCE(AVG(employer_job_ratings.ratings), 0) as average_rating").
		Joins("LEFT JOIN jobs ON jobs.employer_id = employers.id").
		Joins("LEFT JOIN employer_job_ratings ON employer_job_ratings.job_id = jobs.id").
		Where("employers.id = ?", employerID).
		Group("employers.id").
		Scan(&averageRating).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error", "details": err.Error()})
	    return
	}

	//total completed jobs
	err = config.DB.Table("jobs").
		Where("employer_id = ? AND status = ?", employerID, "completed").
				Count(&totalJobCompleted).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error in completed jobs",})
			return
		}

	// total reccomendations
	if err := config.DB.Table("employer_job_ratings").
						Joins("JOIN jobs ON jobs.id = employer_job_ratings.job_id").
						Where("jobs.id = ?", employerID).
						Count(&totalReccomendations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error", "details": err.Error()})
	    return
	}

	// return success
	c.JSON(http.StatusOK, gin.H{
		"rating":averageRating,
		"total_recommendations":totalReccomendations,
		"total_jobs_completed":totalJobCompleted,
		"total_jobs_posted":totalJobPosted,
	})
}
