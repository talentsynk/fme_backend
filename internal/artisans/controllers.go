package artisans

import (
	"encoding/csv"
	"fme_backend/internal/config"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)


func GetArtisanProfile(c *gin.Context) {
	idStr:= c.Param("id")
		if idStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "path parameter not provided",
			})
			return
		}

    id, err := strconv.Atoi(idStr)

    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "message": "path parameter invalid",
        })
        return
    }

	var artisan struct {
		ID                uint
		BusinessName      string
		BusinessDescription string
		Skill             string
		AverageRating     float64
		FirstName string
		LastName  string
	}
	
	err = config.DB.Table("artisans").
		Select("artisans.id, artisans.business_name, artisans.business_description, artisans.first_name AS first_name, artisans.last_name AS last_name, artisans.skill, COALESCE(AVG(job_application_ratings.rating), 0) as average_rating").
		Joins("LEFT JOIN job_applications ON job_applications.artisan_id = artisans.id").
		Joins("LEFT JOIN job_application_ratings ON job_application_ratings.job_application_id = job_applications.id").
		Where("artisans.id = ?", id).
		Group("artisans.id").
		First(&artisan).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error trying to fetch data",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"artisan":artisan})


}


func GetAllArtisans(c * gin.Context) {
	var artisans []struct {
		ID                uint
		BusinessName      string
		BusinessDescription string
		Skill             string
		AverageRating     float64
		FirstName 			string
		LastName 			string
	}

	var minAverageRating,maxAverageRating float64
	minAverageRating = 0
	maxAverageRating = 5


	var queryParams ArtisanFilterSchema
	if err := c.ShouldBindQuery(&queryParams); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
	if queryParams.MinRating >= 0 && queryParams.MinRating <= 5 {
		minAverageRating = queryParams.MinRating
	}
	
	if queryParams.MaxRating > 0 && queryParams.MaxRating <= 5 {
		maxAverageRating = queryParams.MaxRating
	}
	query := config.DB.Table("artisans").
		Select("artisans.id, artisans.business_name, artisans.business_description, artisans.first_name, artisans.last_name, artisans.skill, COALESCE(AVG(job_application_ratings.rating), 0) as average_rating").
		Joins("LEFT JOIN job_applications ON job_applications.artisan_id = artisans.id").
		Joins("LEFT JOIN job_application_ratings ON job_application_ratings.job_application_id = job_applications.id").
		Group("artisans.id").
		Having("COALESCE(AVG(job_application_ratings.rating), 0) >= ? AND COALESCE(AVG(job_application_ratings.rating), 0) <= ?", minAverageRating,maxAverageRating)

	if queryParams.RatingSort {
		query.Order("average_rating DESC")
	}

	err := query.Scan(&artisans).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message":err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"artisans":artisans})
}

func GetArtisanJobStat(c *gin.Context){
	// get the artisan id 

	artisanIdStr, exists := c.Get("artisanID")
	if !exists{
		c.JSON(http.StatusUnauthorized, gin.H{"message":"problem with the authorization token"})
	     return
	}
	artisanID, ok := artisanIdStr.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized,gin.H{"message":"Problem with the authorization token"})
		return
	}

	// needed variables
	var totalAppliedJobs int64
	var totalJobRecommendation int64
	var totalCompletedJobs int64

   // Count the total number of jobs applied for
	if err := config.DB.Table("job_applications").Where("artisan_id = ?", artisanID).Count(&totalAppliedJobs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error", "details": err.Error()})
	    return
	}

	// count the total job recommendations 
	if err := config.DB.Table("job_application_ratings").Joins("JOIN job_applications ON job_application_ratings.job_application_id = job_applications.id").Where("job_applications.artisan_id = ?", artisanID).Count(&totalJobRecommendation).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error", "details": err.Error()})
	    return
	}

	// get completed jobs details 
	if err := config.DB.Table("job_applications").Joins("JOIN jobs ON job_applications.job_id = jobs.id").Where("job_applications.artisan_id = ? AND jobs.status = ?", artisanID,"completed").Count(&totalCompletedJobs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error", "details": err.Error()})
	    return
	}



	c.JSON(http.StatusOK, gin.H{
		"total_applied_jobs":totalAppliedJobs,
		"total_job_recommendations":totalJobRecommendation,
		"total_jobs_completed":totalCompletedJobs,
	})


}

func GetArtisanProfileStat(c *gin.Context) {
	// get the artisan id 

	artisanIdStr, exists := c.Get("artisanID")
	if !exists{
		c.JSON(http.StatusUnauthorized, gin.H{"message":"problem with the authorization token"})
	     return
	}
	artisanID, ok := artisanIdStr.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized,gin.H{"message":"Problem with the authorization token"})
		return
	}

	// needed variables
	var projectsCompleted int64
	var averageRating float64
	var totalReccomendations int64
	
	// get the rating 
	err := config.DB.Table("artisans").
		Select("COALESCE(AVG(job_application_ratings.rating), 0) as average_rating").
		Joins("LEFT JOIN job_applications ON job_applications.artisan_id = artisans.id").
		Joins("LEFT JOIN job_application_ratings ON job_application_ratings.job_application_id = job_applications.id").
		Where("artisans.id = ?", artisanID).
		Group("artisans.id").
		Scan(&averageRating).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error", "details": err.Error()})
	    return
	}

	// get the projects completed 
	if err := config.DB.Table("job_applications").
						Joins("JOIN jobs ON job_applications.job_id = jobs.id").
						Where("job_applications.artisan_id = ? AND jobs.status = ?", artisanID,"completed").
						Count(&projectsCompleted).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error", "details": err.Error()})
	    return
	}


	// get the total reccomendations 
	if err := config.DB.Table("job_application_ratings").
						Joins("JOIN job_applications ON job_application_ratings.job_application_id = job_applications.id").
						Where("job_applications.artisan_id = ?", artisanID).
						Count(&totalReccomendations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error", "details": err.Error()})
	    return
	}

	c.JSON(http.StatusOK, gin.H{
		"rating":averageRating,
		"total_job_recommendations":totalReccomendations,
		"total_jobs_completed":projectsCompleted,
	})

}

func GetArtisanProfileStatByEmployer(c *gin.Context) {
	// get the artisan profile 
	idStr:= c.Param("id")
		if idStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "path parameter not provided",
			})
			return
		}

    artisanID, err := strconv.Atoi(idStr)

    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "message": "path parameter invalid",
        })
        return
    }

	// needed variables
	var projectsCompleted int64
	var averageRating float64
	var totalReccomendations int64
	
	// get the rating 
	err = config.DB.Table("artisans").
		Select("COALESCE(AVG(job_application_ratings.rating), 0) as average_rating").
		Joins("LEFT JOIN job_applications ON job_applications.artisan_id = artisans.id").
		Joins("LEFT JOIN job_application_ratings ON job_application_ratings.job_application_id = job_applications.id").
		Where("artisans.id = ?", artisanID).
		Group("artisans.id").
		Scan(&averageRating).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error", "details": err.Error()})
	    return
	}

	// get the projects completed 
	if err := config.DB.Table("job_applications").
						Joins("JOIN jobs ON job_applications.job_id = jobs.id").
						Where("job_applications.artisan_id = ? AND jobs.status = ?", artisanID,"completed").
						Count(&projectsCompleted).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error", "details": err.Error()})
	    return
	}


	// get the total reccomendations 
	if err := config.DB.Table("job_application_ratings").
						Joins("JOIN job_applications ON job_application_ratings.job_application_id = job_applications.id").
						Where("job_applications.artisan_id = ?", artisanID).
						Count(&totalReccomendations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error", "details": err.Error()})
	    return
	}

	c.JSON(http.StatusOK, gin.H{
		"rating":averageRating,
		"total_job_recommendations":totalReccomendations,
		"total_jobs_completed":projectsCompleted,
	})
}

func GetArtisanRating(c *gin.Context) {
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

	var ratings []struct {
		CreatedAt    time.Time
		Rating       uint
		Description  string
		EmployerID   uint
		FirstName    string
		LastName     string
	}
		// get the filters 
		var queryParams RatingFilterSchema
		if err := c.ShouldBindQuery(&queryParams); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	
	query := config.DB.Table("job_application_ratings").
		Joins("JOIN job_applications ON job_applications.id = job_application_ratings.job_application_id").
		Joins("JOIN jobs ON jobs.id = job_applications.job_id").
		Joins("JOIN employers ON employers.id = jobs.employer_id").
		Select(`job_application_ratings.created_at AS created_at,
				job_application_ratings.rating AS rating,
				job_application_ratings.description AS description,
				jobs.employer_id AS employer_id,
				employers.first_name AS first_name,
				employers.last_name AS last_name`).
		Where("job_applications.artisan_id = ?", id)

		if queryParams.DaysAgo != 0 {
			someDaysAgo := time.Now().AddDate(0,0,-int(queryParams.DaysAgo))
			query.Where("job_application_ratings.created_at >= ?",someDaysAgo)
		}
		if queryParams.MaxRating != 0 {
			query.Where("job_application_ratings.ratings <= ?",queryParams.MaxRating)
		}
	
		if queryParams.MinRating != 0 {
			query.Where("job_application_ratings.ratings >= ?",queryParams.MinRating)
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


func GetMyDetails(c *gin.Context) {
	artisanIdStr, exists := c.Get("artisanID")
	if !exists{
		c.JSON(http.StatusUnauthorized, gin.H{"message":"problem with the authorization token"})
	     return
	}
	artisanID, ok := artisanIdStr.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized,gin.H{"message":"Problem with the authorization token"})
		return
	}

	var result struct {
		Email       string
		FirstName   string
		LastName    string
	}
	
	err := config.DB.Table("artisans").
		Select("users.email, artisans.first_name, artisans.last_name").
		Joins("JOIN users ON users.id = artisans.user_id").
		Where("artisans.id = ?", artisanID).
		Find(&result).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message":"error querying the database"})
	     return
	}

	c.JSON(http.StatusOK,gin.H{
		"artisan":result,
	})


}


func GetContactDetails(c *gin.Context) {
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

	var result struct {
		Email       string
		PhoneNumber string
	}
	
	err = config.DB.Table("artisans").
		Select("users.email, students.phone_number").
		Joins("JOIN students ON students.user_id = artisans.user_id").
		Joins("JOIN users ON users.id = artisans.user_id").
		Where("artisans.id = ?", id).
		Find(&result).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message":"error querying the database"})
	     return
	}

	c.JSON(http.StatusOK,gin.H{
		"contact":result,
	})

}

func DownloadArtisanData(c *gin.Context) {
	var artisans []ArtisanDataDownload
	// get the student data 
	err := 	config.DB.Table("artisans").
						Select(`
						artisans.first_name AS first_name,
						artisans.last_name AS last_name,
						artisans.id AS artisan_id,
						jobs.job_title AS job_title,
						jobs.budget AS budget,
						jobs.job_type AS job_type,
						jobs.location AS job_location,
						jobs.description AS job_description,
						job_applications.application_status AS  application_status,
						jobs.status AS job_status,
						job_application_ratings.rating AS artisan_rating,
						job_application_ratings.description AS artisan_rating_description
						`).
						Joins("JOIN job_applications on job_applications.artisan_id = artisans.id").
						Joins("JOIN jobs ON job_applications.job_id = jobs.id").
						Joins("JOIN job_application_ratings on job_application_ratings.job_application_id = job_applications.id").
						Scan(&artisans).
						Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message":"error querying the database"})
	     return
	}

	// parse it into csv format 
	csvData := [][]string{
        {"FirstName", "LastName", "ArtisanID", "JobTitle", "Budget", "JobType","JobLocation", "JobDescription", "ApplicationStatus", "JobStatus","ArtisanPerformanceRating","ArtisanRatingDescription"},
    }

	for _,artisan := range artisans {
		csvData = append(csvData,[]string{
			artisan.FirstName,
			artisan.LastName,
			artisan.ArtisanID,
			artisan.JobTitle,
			strconv.FormatFloat(float64(artisan.Budget),'f',4,64),
			artisan.JobType,
			artisan.JobLocation,
			artisan.JobDescription,
			artisan.ApplicationStatus,
			artisan.JobStatus,
			artisan.ArtisanRating,
			artisan.ArtisanRatingDescription,
		}) 
	}

	c.Header("Content-Disposition", "attachment; filename=marketplace.csv")
    c.Header("Content-Type", "text/csv")

	//send it out
	w := csv.NewWriter(c.Writer)
    defer w.Flush()

    for _, record := range csvData {
        if err := w.Write(record); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"message": "error writing csv"})
            return
        }
    }

}


func DownloadSingleArtisanData(c *gin.Context) {
	//get the artisan ID


	//get the student data 

	// parse it into the csv

	// send it out 
}