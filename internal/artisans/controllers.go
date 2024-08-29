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