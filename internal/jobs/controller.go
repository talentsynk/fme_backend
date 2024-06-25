package jobs

import (
	"errors"
	"fme_backend/internal/config"
	"fmt"
	"net/http"
	"strconv"
	"time"
     "gorm.io/gorm"
	"github.com/gin-gonic/gin"
	
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
		Location:    CreateJobSchema.Location,
		Budget:       CreateJobSchema.Budget,
		JobType:       CreateJobSchema.JobType,
		Category:     CreateJobSchema.Category,
		Description:  CreateJobSchema.Description,
		Requirement: CreateJobSchema.Requirement,
		Responsibilities: CreateJobSchema.Responsibilities,
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


func GetAllJobs(c *gin.Context){

	   limitStr := c.Query("limit")
	   pageStr := c.Query("page")
	
	   

	   limit, err := strconv.Atoi(limitStr)
	   if err != nil || limit <= 0 {
		limit = 18
	   }

	   page, err := strconv.Atoi(pageStr)
	   if err != nil || page <= 0 {
		page = 1 
	   }

	   offset := (page - 1) * limit 

	  var jobs []GetAllJobsSchema

	 if   result := config.DB.Table("jobs").
	    Select("jobs.id AS id, jobs.created_at AS created_at, jobs.job_type AS job_type, jobs.location AS location, jobs.budget AS budget, jobs.description AS description, jobs.employer_id AS employer_id, employers.first_name AS first_name, employers.last_name As last_name").
		Joins("JOIN employers ON jobs.employer_id = employers.id").
		Limit(limit).
		Offset(offset).
		Find(&jobs); result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		      return
		}

		c.JSON(http.StatusOK, gin.H{"jobs": jobs})
}

func GetJobID(c *gin.Context){
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


	

	var job GetJobSchema

	
	result := config.DB.Table("jobs").
	    Select("jobs.id AS id, jobs.location AS location, jobs.description AS description, jobs.job_type AS job_type, jobs.job_title AS job_title, jobs.requirement AS requirement, jobs.responsibilities AS responsibilities, jobs.created_at AS created_at, employers.first_name AS first_name, employers.last_name AS last_name, employers.id AS employer_id").
		Joins("JOIN employers ON jobs.employer_id = employers.id").
		Where("jobs.id = ?", id).
		First(&job)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":"Job not found",
		})
		return
	}
	c.JSON(http.StatusOK, job)
}


func GetLatestJobs(c *gin.Context){
	limitStr := c.Query("limit")
	pageStr := c.Query("page")
 
	

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
	 limit = 18
	}


	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
	 page = 1 
	}

	offset := (page - 1) * limit 

	
	sevenDaysAgo := time.Now().AddDate(0,0,-7)
   var jobs []GetAllJobsSchema

  if   result := config.DB.Table("jobs").
	 Select("jobs.id AS id, jobs.created_at AS created_at, jobs.job_type AS job_type, jobs.location AS location, jobs.budget AS budget, jobs.description AS description, jobs.employer_id AS employer_id, employers.first_name AS first_name, employers.last_name As last_name").
	 Joins("JOIN employers ON jobs.employer_id = employers.id").
	 Where("jobs.created_at >= ?", sevenDaysAgo).
	 Limit(limit).
	 Offset(offset).
	 Find(&jobs); result.Error != nil {
		 c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		   return
	 }

	 c.JSON(http.StatusOK, gin.H{"jobs": jobs})

}

func GetAllAppliedJobs(c *gin.Context){

}

func GetAllSavedJobs(c *gin.Context){

}

func GetAllShorttermJobs(c *gin.Context){

}

func GetAllFulltimeJobs(c *gin.Context){

}


func SearchJob(c *gin.Context){
	query := c.Query("query")
	if query == ""{
		c.JSON(http.StatusBadGateway, gin.H{"error":"Search query is required"})
		return
	}


	var jobsearch []GetAllJobsSchema
	if err := config.DB.Table("jobs").
	Select("jobs.id AS id, jobs.created_at AS created_at, jobs.job_type AS job_type, jobs.location AS location, jobs.budget AS budget, jobs.description AS description, jobs.job_title AS job_title, jobs.requirement AS requirement, jobs.responsibilities AS responsibilities, jobs.category AS category, jobs.employer_id AS employer_id, employers.first_name AS first_name, employers.last_name AS last_name").
	Joins("JOIN employers ON jobs.employer_id = employers.id").
	Where("jobs.job_title LIKE ? OR jobs.location LIKE ? OR jobs.budget LIKE ? OR jobs.category LIKE ? OR jobs.job_type LIKE ? OR jobs.description LIKE ? OR employers.first_name LIKE ? OR employers.last_name LIKE ?", "%"+query+"%", "%"+query+"%", "%"+query+"%", "%"+query+"%", "%"+query+"%", "%"+query+"%", "%"+query+"%", "%"+query+"%").
	Find(&jobsearch).Error; err != nil {
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search for Jobs", "details": err.Error()})
	return
}


  if len(jobsearch) == 0 {
	c.JSON(http.StatusOK, gin.H{"message": "No matching mdas found"})
  }

	c.JSON(http.StatusOK, jobsearch)
}



// func SearchJob(c *gin.Context) {
// 	job_title := c.Query("job_title")
	

// 	var jobsearch []GetAllJobsSchema
// 	if err := config.DB.Table("jobs").
// 		Select("jobs.id AS id, jobs.created_at AS created_at, jobs.job_type AS job_type, jobs.location AS location, jobs.budget AS budget, jobs.description AS description, jobs.job_title AS job_title, jobs.requirement AS requirement, jobs.responsibilities AS responsibilities, jobs.category AS category, jobs.employer_id AS employer_id, employers.first_name AS first_name, employers.last_name AS last_name").
// 		Joins("JOIN employers ON jobs.employer_id = employers.id").
// 		Where("jobs.job_title LIKE ?", "%"+job_title+"%").
// 		Find(&jobsearch).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search for Jobs", "details": err.Error()})
// 		return
// 	}

// 	if len(jobsearch) == 0 {
// 		c.JSON(http.StatusOK, gin.H{"message": "No matching jobs found"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"jobs": jobsearch})
// }



// func ApplyForJob(c *gin.Context) {
//     studentIDStr, exist := c.Get("studentID")
//     if !exist {
//         c.JSON(http.StatusUnauthorized, gin.H{"message": "Problem with the authorization token"})
//         return
//     }

//     studentID, ok := studentIDStr.(uint)
//     if !ok {
//         c.JSON(http.StatusUnauthorized, gin.H{"message": "Problem with the authorization token"})
//         return
//     }

//     // Get the job ID from the request
//     var applyJobSchema struct {
//         JobID uint `json:"job_id" binding:"required"`
//     }

//     if err := c.ShouldBindJSON(&applyJobSchema); err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
//         return
//     }

//     // Begin a new transaction
//     tx := config.DB.Begin()

//     // Check if the application already exists
// 	var existingApplication JobApplication
// 	if result := tx.Where("job_id = ? AND student_id = ?", applyJobSchema.JobID, studentID).First(&existingApplication); result.Error == nil {
// 		tx.Rollback()
// 		c.JSON(http.StatusConflict, gin.H{"error": "You have already applied for this job"})
// 		return
// 	} else if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
// 		tx.Rollback()
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error occurred while checking for existing application"})
// 		return
// 	}
	

//     // Create the new job application
//     jobApplication := JobApplication{
//         JobID:     applyJobSchema.JobID,
//         StudentID: studentID,
//     }

//     if result := tx.Create(&jobApplication); result.Error != nil {
//         tx.Rollback()
//         c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to apply for the job", "details": result.Error.Error()})
//         return
//     }

//     // Commit the transaction
//     if err := tx.Commit().Error; err != nil {
//         tx.Rollback()
//         c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction", "details": err.Error()})
//         return
//     }

//     c.JSON(http.StatusOK, gin.H{"message": "Applied for job successfully"})
// }




// func SaveJobs(c *gin.Context) {
//     studentIDStr, exist := c.Get("studentID")
//     if !exist {
//         c.JSON(http.StatusUnauthorized, gin.H{"message": "Problem with the authorization token"})
//         return
//     }

//     studentID, ok := studentIDStr.(uint)
//     if !ok {
//         c.JSON(http.StatusUnauthorized, gin.H{"message": "Problem with the authorization token"})
//         return
//     }

//     // Get the job ID from the request
//     var savedJobSchema struct {
//         JobID uint `json:"job_id" binding:"required"`
//     }

//     if err := c.ShouldBindJSON(&savedJobSchema); err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
//         return
//     }

//     // Begin a new transaction
//     tx := config.DB.Begin()



//     // Create the new job application
//     savedJob := SaveJob{
//         JobID:     savedJobSchema.JobID,
//         StudentID: studentID,
//     }

//     if result := tx.Create(&savedJob); result.Error != nil {
//         tx.Rollback()
//         c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to  save  job", "details": result.Error.Error()})
//         return
//     }

//     // Commit the transaction
//     tx.Commit()
        

//     c.JSON(http.StatusOK, gin.H{"message":"The job has been saved successfully"})
// }


func ApplyorSaveJob(c *gin.Context){
	studentIDStr, exists := c.Get("studentID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message":"problem with the authorization token"})
	     return
	}

	studentID, ok := studentIDStr.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized,gin.H{"message":"Problem with the authorization token"})
		return
	}


	var jobActionSchema struct{
		JobID  uint `json:"job_id" binding:"required"`
        Action string  `json:"action" binding:"required"`
	}

	if err := c.ShouldBindJSON(&jobActionSchema); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error":"Failed to read request body"})
	    return
	}


	tx := config.DB.Begin()

	switch jobActionSchema.Action {
	case  "apply":
		var existingApplication JobApplication
		if result := tx.Where("job_id = ? AND student_id = ?", jobActionSchema.JobID, studentID).First(&existingApplication); result.Error == nil{
			tx.Rollback()
			c.JSON(http.StatusConflict, gin.H{"error":"You have already applied for this job"})
		     return  
		}else if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound){
           tx.Rollback()
		   c.JSON(http.StatusInternalServerError, gin.H{"error":"Database error occure while checking for existing application"})
		   return
		}
     
		jobApplication := JobApplication{
			JobID: jobActionSchema.JobID,
			StudentID: studentID,
		}

		if result := tx.Create(&jobApplication); result.Error != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error":"Failed to commit tranction", "details": result.Error.Error()})
			return
		}
		
		if err := tx.Commit().Error; err != nil {
            tx.Rollback()
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction", "details": err.Error()})
            return
        }

        c.JSON(http.StatusOK, gin.H{"message": "Applied for job successfully"})

	case "save":
		var existingSave SaveJob
       if result := tx.Where("job_id = ? AND student_id = ?", jobActionSchema.JobID, studentID).First(&existingSave); result.Error == nil {
		  tx.Rollback()
		  c.JSON(http.StatusConflict, gin.H{"error":"You have already save this job"})
		  return
	   }else if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound){
		 tx.Rollback()
		 c.JSON(http.StatusInternalServerError, gin.H{"error":"Database error occured while checking for existing saved job", "details": result.Error.Error()})
		 return
	   }

        savedJob := SaveJob{
			JobID: jobActionSchema.JobID,
			StudentID: studentID,
		}
		
		if result := tx.Create(&savedJob); result.Error  != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error":"Failed  to save job", "detailes":result.Error.Error()})
		     return
		}
		
		if err := tx.Commit().Error; err != nil{
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to commit transaction", "details":err.Error()})
		    return
		}
		c.JSON(http.StatusOK, gin.H{"message":"The job has been saved successfully"})
	default:
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error":"Invalid action"})
	}

}