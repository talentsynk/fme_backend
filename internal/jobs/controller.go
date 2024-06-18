package jobs


import (
	 "fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"fme_backend/internal/config"
)

func CreateJob(c *gin.Context) {
	fmt.Println("controller started")

	// Extract employerID from the context
	employerIDstr, exist := c.Get("employerID")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Problem with the authorization token",
		})
		return
	}

	employerID, ok := employerIDstr.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Problem with the authorization token",
		})
		return
	}


	if err := c.BindJSON(&CreateJobSchema); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read request body",
		})
		return
	}
	tx := config.DB.Begin() 

	job := Job{
		JobTitle:    CreateJobSchema.JobTitle,
		ArtisanType:  CreateJobSchema.ArtisanType,
		Location:    CreateJobSchema.Location,
		Role:         CreateJobSchema.Role,
		Budget:       CreateJobSchema.Budget,
		Time:         CreateJobSchema.Time,
		Category:     CreateJobSchema.Category,
		Description:  CreateJobSchema.Description,
		Requirement: CreateJobSchema.Requirement,
		EmployerID:  employerID,
	}

	// Assuming 'tx' is a valid gorm.DB transaction object available in the context

         jobresult := tx.Create(&job)
		 if jobresult.Error != nil{
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{
				"error":"Failed to create job",
			})
			return
		 }
		 tx.Commit()
	c.JSON(http.StatusOK, gin.H{
		"message": "Job created successfully",
	})
}