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

	var artisan Artisans
	err =  config.DB.First(&artisan, id).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error trying to fetch data",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"artisan":artisan})


}


func GetAllArtisans(c * gin.Context) {
	var artisans []Artisans
	err := config.DB.Find(&artisans)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error trying to fetch data",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"artisans":artisans})
}