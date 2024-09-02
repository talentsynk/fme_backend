package artisans

import (
	"fme_backend/internal/config"
	"net/http"
	"strconv"

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
	}
	
	err = config.DB.Table("artisans").
		Select("artisans.id, artisans.business_name, artisans.business_description, artisans.skill, COALESCE(AVG(job_application_ratings.rating), 0) as average_rating").
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
	
	err := config.DB.Table("artisans").
		Select("artisans.id, artisans.business_name, artisans.business_description, artisans.first_name, artisans.last_name, artisans.skill, COALESCE(AVG(job_application_ratings.rating), 0) as average_rating").
		Joins("LEFT JOIN job_applications ON job_applications.artisan_id = artisans.id").
		Joins("LEFT JOIN job_application_ratings ON job_application_ratings.job_application_id = job_applications.id").
		Group("artisans.id").
		Scan(&artisans).Error
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