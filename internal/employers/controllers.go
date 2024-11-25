package employer

import (
	"fme_backend/internal/config"
	myuser "fme_backend/internal/user"
	"fme_backend/internal/utilities"
	"net/http"
	"strconv"
	"time"

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

	// verify if the user is a company and verify their set their acct details
	

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
	if CreateEmployerSchema.IsCompany {
		employer.IsCompany = true
		employer.CompanyName = CreateEmployerSchema.CompanyName
		employer.CompanyCAC = CreateEmployerSchema.CompanyCAC
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
	Select("id, first_name, last_name, phone_number, nin, state, lga, user_id").
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

func GetEmployerRating(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == ""{
		c.JSON(http.StatusBadRequest, gin.H{"error":"job Id is required"})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message":"Invalid job id provided"})
	     return
	}

	// get the filters 
	var queryParams RatingFilterSchema
	if err := c.ShouldBindQuery(&queryParams); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

	var ratings []struct {
		CreatedAt     time.Time
		Rating        uint
		Description   string
		ArtisanID     uint
		BusinessName  string
	}
	
	query := config.DB.Table("employer_job_ratings").
		Joins("JOIN job_applications ON job_applications.id = employer_job_ratings.job_id").
		Joins("JOIN artisans ON artisans.id = job_applications.artisan_id").
		Joins("JOIN jobs ON jobs.id = job_applications.job_id").
		Select(`employer_job_ratings.created_at AS created_at,
				employer_job_ratings.ratings AS rating,
				employer_job_ratings.description AS description,
				artisans.id AS artisan_id,
				artisans.business_name AS business_name`).
		Where("jobs.employer_id = ?", id)

	if queryParams.DaysAgo != 0 {
		someDaysAgo := time.Now().AddDate(0,0,-int(queryParams.DaysAgo))
		query.Where("employer_job_ratings.created_at >= ?",someDaysAgo)
	}
	if queryParams.MaxRating != 0 {
		query.Where("employer_job_ratings.ratings <= ?",queryParams.MaxRating)
	}

	if queryParams.MinRating != 0 {
		query.Where("employer_job_ratings.ratings >= ?",queryParams.MinRating)
	}


	err = query.Scan(&ratings).Error	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message":"error querying the database"})
	     return
	}

	c.JSON(http.StatusOK,gin.H{
		"ratings":ratings,
	})
	
}

func GetAllEmployerJobs(c *gin.Context) {
	// get employer id
	idStr := c.Param("id")
	if idStr == ""{
		c.JSON(http.StatusBadRequest, gin.H{"error":"job Id is required"})
		return
	}

	employerID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message":"Invalid job id provided"})
	     return
	}

	var queryParams JobFilterSchema
	if err := c.ShouldBindQuery(&queryParams); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

	// query the db for the data using filters and employer id
	var jobs []GetJobSchema
	result := config.DB.Table("jobs").
					 Select("jobs.id AS id, jobs.job_title AS job_title, jobs.description AS description, jobs.budget AS amount, jobs.job_type AS job_type, jobs.status AS status").
					 Where("employer_id = ?", employerID)
	
	// handle filters
	if queryParams.MaxBudget != 0 {
		result.Where("budget <= ?",queryParams.MaxBudget)
	}

	if queryParams.MinBudget != 0 {
		result.Where("budget >= ?",queryParams.MinBudget)
	}

	if queryParams.JobType == "full-time" || queryParams.JobType == "part-time" {
		result.Where("job_type = ?",queryParams.JobType)
	}

	if queryParams.Status == "open" || queryParams.Status == "ongoing" || queryParams.Status == "completed" {
		result.Where("status = ?",queryParams.Status)

	}

	if queryParams.DaysAgo != 0 {
		someDaysAgo := time.Now().AddDate(0,0,-int(queryParams.DaysAgo))
		result.Where("created_at >= ?",someDaysAgo)
	}

	err = result.Order("jobs.created_at DESC").Find(&jobs).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "problem querying database",
		})
		return
	}

	//return the data 
	c.JSON(http.StatusOK,gin.H{
		"jobs":jobs,
	})
}

func GetSimilarEmployerDetails(c * gin.Context) {
	// get the employer details
	idStr := c.Param("id")
	if idStr == ""{
		c.JSON(http.StatusBadRequest, gin.H{"error":"job Id is required"})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message":"Invalid job id provided"})
	     return
	}

	// get the employer state or
	var employerState string 
	if err = config.DB.Table("employers").Select("state").Where("id = ?",id).Scan(&employerState).Error;
			err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message":"error reading database"})
	     return
			}

	// get all employers that are in that area 
	var employers []GetEmployerSchema

	if result := config.DB.Table("employers").
	Select("id, first_name, last_name, phone_number, nin, state, lga, user_id").
	Where("state = ? AND id <> ?",employerState,id).
	Find(&employers); result.Error != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()});
	     return
	}

	// send the success message
	c.JSON(http.StatusOK, gin.H{"employers":employers}) 

}



// // 
// {
//     "name_of_organisation": "Company Ltd",
//     "email": "company@example.com",
//     "phone_number": "234...",
//     "cac": "CAC123456",
//     "state": "Lagos",
//     "lga": "Ikeja",
//     "password": "securepassword"
// }


// // 

// {
//     "first_name": "John",
//     "last_name": "Doe",
//     "email": "john@example.com",
//     "phone_number": "234...",
//     "nin": "1234567890",
//     "state": "Lagos",
//     "lga": "Ikeja",
//     "password": "securepassword"
// }