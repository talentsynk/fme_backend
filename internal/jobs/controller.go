package jobs

import (
	"errors"
	"fme_backend/internal/config"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"fme_backend/internal/utilities"
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

   typeOfJobs, result := utilities.ValidateJobType(CreateJobSchema.JobType)
   if !result {
	c.JSON(http.StatusBadRequest, gin.H{
		"error":"Incorrect job type",
	})
  } 

	tx := config.DB.Begin() 

	job := Job{
		JobTitle:    CreateJobSchema.JobTitle,
		Location:    CreateJobSchema.Location,
		Budget:       CreateJobSchema.Budget,
		JobType:       typeOfJobs,
		Category:     CreateJobSchema.Category,
		Description:  CreateJobSchema.Description,
		Requirement: CreateJobSchema.Requirement,
		Responsibilities: CreateJobSchema.Responsibilities,
		EmployerID:  employerID,
		Status: true,
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
	  var total int64

	// Get the total number of jobs
	if result := config.DB.Table("jobs").
		Count(&total); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}

	 if   result := config.DB.Table("jobs").
	    Select("jobs.id AS id, jobs.created_at AS created_at, jobs.job_type AS job_type, jobs.location AS location, jobs.budget AS budget, jobs.description AS description, jobs.employer_id AS employer_id, employers.first_name AS first_name, employers.last_name As last_name").
		Joins("JOIN employers ON jobs.employer_id = employers.id").
		Limit(limit).
		Offset(offset).
		Find(&jobs); result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		      return
		}

		c.JSON(http.StatusOK, gin.H{"total":total, "jobs": jobs})
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

func GetJobType(c *gin.Context){
	jobType := c.Param("jobType")

    limitStr  :=  c.Query("limit")
	pageStr :=  c.Query("page")


	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 18
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil  || page <= 0{
		page = 1
	}

	offset := (page - 1) * limit
	var total int64

	if result := config.DB.Table("jobs").
	Where("jobs.job_type = ?", jobType).
	Count(&total); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}

	var jobs []GetAllJobsSchema
    if result := config.DB.Table("jobs").
	Select("jobs.id AS id, jobs.created_at AS created_at, jobs.job_type AS job_type,jobs.responsibilities, jobs.requirement, jobs.job_title, jobs.location AS location, jobs.budget AS budget, jobs.description AS description, jobs.employer_id AS employer_id, employers.first_name AS first_name, employers.last_name AS last_name").
	Joins("JOIN employers ON jobs.employer_id = employers.id").
	Where("jobs.job_type = ?", jobType).
	Limit(limit).
	Offset(offset).
	Find(&jobs); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total":total,
		"jobs":jobs,
	})
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


func GetAppliedJobs(c *gin.Context){
	studentIDStr, exists := c.Get("studentID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "problem with the authorization token"})
	     return
	}

	studentID, ok := studentIDStr.(uint)
	if !ok{
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Problem with the authorization token"})
        return
	}

	var appliedJobs []Job
	if err := config.DB.Joins("JOIN job_applications ON jobs.id = job_applications.job_id").
        Where("job_applications.student_id = ?", studentID).Find(&appliedJobs).Error; err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":"Database error", "details":err.Error()})
			return 
		}
		c.JSON(http.StatusOK, gin.H{"applied_jobs":appliedJobs})
  }




  func GetAllAppliedJobs(c *gin.Context){
	var appliedJobs []struct {
		JobID     uint   `json:"job_id"`
		FirstName        string 
		LastName         string 
		Location         string
		Description      string 
		JobType          string  
		JobTitle         string 
		Requirement      string
		Responsibilities string
        StudentID uint   `json:"student_id"`
        AppliedAt time.Time `json:"applied_at"` 
		EmployerID        uint  
	}

	if err := config.DB.Table("jobs").
	 Select("jobs.id as job_id, jobs.job_title, jobs.location, jobs.description, jobs.job_type, jobs.requirement, jobs.responsibilities,employers.first_name AS first_name, employers.last_name AS last_name, employers.id AS employer_id,job_applications.student_id, job_applications.created_at as applied_at").
     Joins("JOIN job_applications ON jobs.id = job_applications.job_id"). 
	 Joins("JOIN employers ON jobs.employer_id = employers.id"). 
     Find(&appliedJobs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error", "details": err.Error()})
        return
	}
	c.JSON(http.StatusOK, gin.H{"applied_jobs":appliedJobs})
}

func GetJobStat(c *gin.Context){
	var totalAppliedJobs int64
	var totalJobRecommendation int64
	var totalCompletedJobs int64

   // Count the total number of jobs applied for
	if result := config.DB.Table("job_applications").Count(&totalAppliedJobs); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error", "details": result.Error.Error()})
	    return
	}


	 // Count the total number of job recommendations
	if result := config.DB.Table("job_recommendations").Count(&totalJobRecommendation); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error", "details": result.Error.Error()})
	    return
	}

	if result := config.DB.Table("completed_jobs").Count(&totalCompletedJobs); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error", "details": result.Error.Error()})
		return
	}


	c.JSON(http.StatusOK, gin.H{
		"total_applied_jobs":totalAppliedJobs,
		"total_job_recommendations":totalJobRecommendation,
		"total_completed_jobs":totalCompletedJobs,
	})


}
func GetStudentJobStat(c *gin.Context) {
	studentIDStr, exists := c.Get("studentID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Problem with the authorization token"})
		return
	}

	studentID, ok := studentIDStr.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Problem with the authorization token"})
		return
	}

	var totalAppliedJobs int64
	var totalJobRecommendation int64
	var totalCompletedJobs int64

	// Count the total number of jobs applied for by the student
	if result := config.DB.Table("job_applications").Where("student_id = ?", studentID).Count(&totalAppliedJobs); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error", "details": result.Error.Error()})
		return
	}

	// Count the total number of job recommendations for the student
	if result := config.DB.Table("job_recommendations").Where(" job_recommendations.student_id = ?", studentID).Count(&totalJobRecommendation); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error", "details": result.Error.Error()})
		return
	}

	// Count the total number of completed jobs for the student
	if result := config.DB.Table("completed_jobs").Where("completed_jobs.student_id = ?", studentID).Count(&totalCompletedJobs); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error", "details": result.Error.Error()})
		return
	}


	c.JSON(http.StatusOK, gin.H{
		"student_id":                studentID,
		"total_applied_jobs":        totalAppliedJobs,
		"total_job_recommendations": totalJobRecommendation,
		"total_completed_jobs":      totalCompletedJobs,
	})
}


func GetSavedJobs(c *gin.Context){
	studentIDStr, exists := c.Get("studentID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "problem with the authorization token"})
	     return
	}

	studentID, ok := studentIDStr.(uint)
	if !ok{
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Problem with the authorization token"})
        return
	}

	var savedJobs []Job
	if err := config.DB.Joins("JOIN save_jobs ON jobs.id = save_jobs.job_id").
        Where("save_jobs.student_id = ?", studentID).Find(&savedJobs).Error; err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":"Database error", "details":err.Error()})
			return 
		}
		c.JSON(http.StatusOK, gin.H{"saved_jobs":savedJobs})
  }




  func GetAllSavedJobs(c *gin.Context){
	var savedJobs []struct {
		JobID             uint   `json:"job_id"`
		FirstName         string 
		LastName          string 
		Location          string
		Description       string 
		JobType           string  
		JobTitle          string 
		Requirement       string
		Responsibilities  string
        StudentID         uint   `json:"student_id"`
        AppliedAt         time.Time `json:"applied_at"` 
		EmployerID        uint  
	}

	if err := config.DB.Table("jobs").
	 Select("jobs.id as job_id, jobs.job_title, jobs.location, jobs.description, jobs.job_type, jobs.requirement, jobs.responsibilities,employers.first_name AS first_name, employers.last_name AS last_name, employers.id AS employer_id,save_jobs.student_id, save_jobs.created_at as saved_at").
     Joins("JOIN save_jobs ON jobs.id = save_jobs.job_id"). 
	 Joins("JOIN employers ON jobs.employer_id = employers.id"). 
     Find(&savedJobs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error", "details": err.Error()})
        return
	}
	c.JSON(http.StatusOK, gin.H{"save_jobs":savedJobs})
}