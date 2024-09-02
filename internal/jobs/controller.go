package jobs

import (
	"fme_backend/internal/config"
	"fme_backend/internal/utilities"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateJob(c *gin.Context) {
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

  // create job object
  job := Job{
	JobTitle:    CreateJobSchema.JobTitle,
	Location:    CreateJobSchema.Location,	//revisit
	Budget:       CreateJobSchema.Budget,
	JobType:       typeOfJobs,
	Category:     CreateJobSchema.Category,
	Description:  CreateJobSchema.Description,
	Requirement: CreateJobSchema.Requirement,
	Responsibilities: CreateJobSchema.Responsibilities,
	EmployerID:  employerID,
	HiringStatus: true,
	Status: "open",
}

	// save the job object 
	err := config.DB.Save(&job).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error trying to create the job",
		})
		return
	}
	c.JSON(http.StatusOK,gin.H{
		"message":"job created cuccesfully",
	})
}

// add akk necessary filters to this and any auth user
func GetAllJobs(c *gin.Context){

	// add filters

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

	   var queryParams JobFilterSchema
	if err := c.ShouldBindQuery(&queryParams); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

	// query the db for the data using filters and employer id
	var jobs []GetJobSchema
	result := config.DB.Table("jobs").
					 Select("jobs.id AS id, jobs.job_title AS job_title, jobs.description AS description, jobs.budget AS amount, jobs.job_type AS job_type, jobs.status AS status")
	
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

	err = result.Limit(limit).Offset(offset).Order("jobs.created_at DESC").Find(&jobs).Error

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

	var job GetJobByIdsSchema
	
	result := config.DB.Table("jobs").
						Select("jobs.id AS id, jobs.job_title AS job_title, employers.id AS employer_id, employers.first_name AS employer_first_name, employers.last_name AS employer_last_name, jobs.skills AS skills, jobs.created_at AS created_at, jobs.location AS location, jobs.job_type AS job_type, jobs.requirement AS requirements, jobs.responsibilities AS responsibilities, jobs.description AS description, jobs.status AS status, jobs.hiring_status AS hiring_status").
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

// only artisans 
func SaveNewJob(c *gin.Context) {
	// get the job id
	idStr := c.Param("id")
	if idStr == ""{
		c.JSON(http.StatusBadRequest, gin.H{"error":"job Id is required"})
		return
	}

	jobId, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message":"Invalid job id provided"})
	     return
	}


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

	//create the save job model
	savedJob := SaveJob{
		ArtisanID: artisanID,
		JobID: uint(jobId),
	}
	var savedId uint

	// check if the job has been saved previously
	err = config.DB.Table("save_jobs").
				Select("id").
				Where("job_id = ? AND artisan_id = ?", savedJob.JobID, savedJob.ArtisanID).
				First(&savedId).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError,
				 gin.H{"message":"Error scanning database","savedid":savedId,})
	     return
		} 
		
		if savedId != 0 {
			c.JSON(http.StatusBadRequest, gin.H{"message":"Job already saved"})
	     return
		}

	// save the job
	err = config.DB.Save(&savedJob).Error
	if err !=nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error trying to save the job",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "job saved successfully",
	})
}

func ApplyForJob(c * gin.Context) {
	// get the job id
	idStr := c.Param("id")
	if idStr == ""{
		c.JSON(http.StatusBadRequest, gin.H{"error":"job Id is required"})
		return
	}

	jobId, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message":"Invalid job id provided"})
	     return
	}


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

	//create the save job model
	appliedJob := JobApplication{
		ArtisanID: artisanID,
		JobID: uint(jobId),
	}

	var appliedId uint

	// check if the job has been saved previously
	err = config.DB.Table("job_applications").
				Select("id").
				Where("job_id = ? AND artisan_id = ?", appliedJob.JobID, appliedJob.ArtisanID).
				Find(&appliedId).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError,
				 gin.H{"message":"Error scanning database","savedid":appliedId,})
	     return
		} 
		
		if appliedId != 0 {
			c.JSON(http.StatusBadRequest, gin.H{"message":"Job already saved"})
	     return
		}

	// save the job
	err = config.DB.Save(&appliedJob).Error
	if err !=nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error trying to save the job",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "job applied for successfully",
	})
}

func GetAppliedJobs(c *gin.Context){
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

	var appliedJobs []GetAppliedJobSchema
	// get all the applied jobs
	err := config.DB.Table("jobs").
    Select("jobs.id AS id, jobs.job_title AS name, jobs.description AS description, job_applications.application_status AS application_status").
    Joins("JOIN job_applications ON job_applications.job_id = jobs.id").
    Where("job_applications.artisan_id = ?", artisanID).
    Scan(&appliedJobs).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error trying retrieve jobs data",
		})
		return
	}
	
	c.JSON(http.StatusOK,gin.H{
		"jobs": appliedJobs,
	})
	

}

func GetArtisanJobProfile(c *gin.Context) {
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

	var artisanDetails struct {
		TotalApplications  int64
		AverageRating      float64
		TotalRatings       int64
		BusinessName       string
		Skills             string
		ArtisanID          uint
		BusinessDescription string
	}
	
	err = config.DB.Table("artisans").
		Joins("LEFT JOIN job_applications ON job_applications.artisan_id = artisans.id").
		Joins("LEFT JOIN job_application_ratings ON job_application_ratings.job_application_id = job_applications.id").
		Select(`artisans.id AS artisan_id,
				artisans.business_name,
				artisans.skill AS skills,
				artisans.business_description,
				COUNT(DISTINCT job_applications.id) AS total_applications,
				COALESCE(AVG(job_application_ratings.rating), 0) AS average_rating,
				COUNT(job_application_ratings.rating) AS total_ratings`).
		Where("artisans.id = ?", id).
		Group("artisans.id").
		Scan(&artisanDetails).Error

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error":"error retrieving db data"})
	     return
		}
		c.JSON(http.StatusOK, gin.H{"artisan":artisanDetails})
	
}


func GetSavedJobs(c *gin.Context){
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

	var savedJobs []GetSavedJobSchema
	// get all the applied jobs
	err := config.DB.Table("jobs").
    Select("jobs.id AS id, jobs.job_title AS name, jobs.description AS description, jobs.budget AS amount, jobs.job_type AS job, jobs.location AS location").
    Joins("JOIN save_jobs ON save_jobs.job_id = jobs.id").
    Where("save_jobs.artisan_id = ?", artisanID).
    Scan(&savedJobs).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error trying retrieve jobs data",
		})
		return
	}
	
	c.JSON(http.StatusOK,gin.H{
		"jobs": savedJobs,
	})
}

// accept a job 
// confirm if the job has already been hired out by checking job status
func HireArtisan(c *gin.Context) {
	// get the job id and the artisan id
	if err := c.BindJSON(&HireArtisanSchema); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read request body",
		})
		return
	}
	// get the employer id
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

	var jobApplication JobApplication
	var job Job

	err := config.DB.
				Where("id = ?", HireArtisanSchema.JobId).
				First(&job).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to retrieve job data",
		})
		return
	}

	// check if the employer can work on the job
	if employerID != job.EmployerID {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Wrong employer details",
		})
		return
	}


	// change the application status and job status
	err = config.DB.
				Where("artisan_id = ? AND job_id = ?", HireArtisanSchema.ArtisanId, HireArtisanSchema.JobId).
				First(&jobApplication).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to retrieve job data 2",
		})
		return
	}
	
	
	tx := config.DB.Begin()

	// change the job application
	jobApplication.ApplicationStatus = "hired"
	if err = tx.Save(&jobApplication).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to error updating job application status",
		})
		return
	}

	//change the job status
	job.Status = "ongoing"
	if err = tx.Save(&job).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to error updating job status",
		})
		return
	}

	//end transaction
	if err = tx.Commit().Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to commit the changes",
		})
		return
	}

	c.JSON(http.StatusOK,gin.H{
		"message": "artisan hired successfully",
	})
}

// view all applicants fix up rating
func ViewApplicants(c *gin.Context) {
	// get the job id 
	idStr := c.Param("id")
	if idStr == ""{
		c.JSON(http.StatusBadRequest, gin.H{"error":"job Id is required"})
		return
	}

	jobId, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message":"Invalid job id provided"})
	     return
	}


	// verify if the employer has aceess to view the applicants
	// get the employer id
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

	var job Job
	var applicants []GetApplicantSchema

	err = config.DB.Table("jobs").
				Where("id = ?", jobId).
				First(&job).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to retrieve job data 1",
		})
		return
	}

	// check if the employer can work on the job
	if employerID != job.EmployerID {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Wrong employer details",
		})
		return
	}


	// get all the artisan info from the applications.
	// write query to get all the artisans that match that jobid
	subQuery := config.DB.Table("job_application_ratings").
    Select("job_applications.artisan_id, COALESCE(AVG(job_application_ratings.rating), 0) as average_rating").
    Joins("JOIN job_applications ON job_applications.id = job_application_ratings.job_application_id").
    Group("job_applications.artisan_id")

err = config.DB.Table("artisans").
    Select("job_applications.id AS application_id, job_applications.application_status AS application_status,artisans.id AS artisan_id, artisans.business_name AS business_name, artisans.first_name, artisans.last_name, artisans.business_description AS business_description, job_applications.created_at AS job_application_date, COALESCE(subquery.average_rating, 0) AS average_rating").
    Joins("JOIN job_applications ON job_applications.artisan_id = artisans.id").
    Joins("LEFT JOIN (?) AS subquery ON subquery.artisan_id = artisans.id", subQuery). // Use LEFT JOIN to ensure all artisans are included
    Where("job_applications.job_id = ?", jobId).
    Scan(&applicants).Error
	
		if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to retrieve job data",
		})
		return
	}

	c.JSON(http.StatusOK,gin.H{
		"applicants": applicants,
	})
}

// complete a job with ratings  fix up the 
func CompleteJob(c *gin.Context) {
	// get the request data 
	if err := c.BindJSON(&CompleteJobSchema); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read request body",
		})
		return
	}

	// get the job id 
	idStr := c.Param("id")
	if idStr == ""{
		c.JSON(http.StatusBadRequest, gin.H{"error":"job Id is required"})
		return
	}

	jobId, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message":"Invalid job id provided"})
	     return
	}


	// verify if the employer has aceess to view the applicants
	// get the employer id
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

	var job Job

	err = config.DB.
				Where("id = ?", jobId).
				First(&job).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to retrieve job data",
		})
		return
	}

	// check if the employer can work on the job
	if employerID != job.EmployerID {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Wrong employer details",
		})
		return
	}

	if job.Status == "completed" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "job has already been rated",
		})
		return
	}

	// update the job status
	tx := config.DB.Begin()

	job.Status  = "completed"

	err = tx.Save(&job).Error
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to change job status",
		})
		return

	}
	

	// rate the artisan
	// get the application id
	var jobApplication JobApplication
	err = config.DB.Where("job_id = ? AND application_status = ?", jobId, "hired").
    Find(&jobApplication).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to get the jobApplication",
		})
		return
	}


	rating := JobApplicationRating{
		JobApplicationID: jobApplication.ID,
		Rating: CompleteJobSchema.Rating,
		Description: CompleteJobSchema.Description,
		
	}

	err = tx.Save(&rating).Error
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to change job status",
		})
		return
	}

	tx.Commit()
	//return successful
	c.JSON(http.StatusOK, gin.H{
		"message": "rating successful",
	})
}

//get jobs stats

// get my personal jobs with filters
func GetMyJobs(c *gin.Context) {
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


	// get the filters 
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

	err := result.Order("jobs.created_at DESC").Find(&jobs).Error

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

// close a job application
func CloseJob(c *gin.Context) {
	// get the job id 
	// get the job id 
	idStr := c.Param("id")
	if idStr == ""{
		c.JSON(http.StatusBadRequest, gin.H{"error":"job Id is required"})
		return
	}

	jobId, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message":"Invalid job id provided"})
	     return
	}


	// verify if the employer has aceess to view the applicants
	// get the employer id
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

	var job Job

	err = config.DB.Table("jobs").
				Where("id = ?", jobId).
				First(&job).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to retrieve job data 1",
		})
		return
	}

	// check if the employer can work on the job
	if employerID != job.EmployerID {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Wrong employer details",
		})
		return
	}

	if job.Status != "open" {
		c.JSON(http.StatusBadRequest,gin.H{
			"error": "job is currently not open",
		})
		return
	}
	if !job.HiringStatus {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "job already closed",
		})
		return
	}

	job.HiringStatus = false
	err = config.DB.Save(&job).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to change job status",
		})
		return
	}
	c.JSON(http.StatusOK,gin.H{
		"message":"job closed successfully",
	})

}

func OpenJob(c *gin.Context) {
	// get the job id 
	// get the job id 
	idStr := c.Param("id")
	if idStr == ""{
		c.JSON(http.StatusBadRequest, gin.H{"error":"job Id is required"})
		return
	}

	jobId, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message":"Invalid job id provided"})
	     return
	}


	// verify if the employer has aceess to view the applicants
	// get the employer id
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

	var job Job

	err = config.DB.Table("jobs").
				Where("id = ?", jobId).
				First(&job).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to retrieve job data 1",
		})
		return
	}

	// check if the employer can work on the job
	if employerID != job.EmployerID {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Wrong employer details",
		})
		return
	}

	if job.Status != "open" {
		c.JSON(http.StatusBadRequest,gin.H{
			"error": "job is currently no open",
		})
		return
	}

	if job.HiringStatus {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "job already open",
		})
		return
	}

	job.HiringStatus = true
	err = config.DB.Save(&job).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to change job status",
		})
		return
	}
	c.JSON(http.StatusOK,gin.H{
		"message":"job opened successfully",
	})

}

func RateEmployer(c *gin.Context) {

	if err := c.BindJSON(&CreateJobSchema); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read request body",
		})
		return
	}

	idStr := c.Param("id")
	if idStr == ""{
		c.JSON(http.StatusBadRequest, gin.H{"error":"job Id is required"})
		return
	}

	jobId, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message":"Invalid job id provided"})
	     return
	}

	// check if the job is completed
	var job Job
	err = config.DB.First(&job,jobId).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message":"Error getting Job"})
	     return
	}
	if job.Status != "completed" {
		c.JSON(http.StatusBadRequest, gin.H{"message":"Cannot rate uncompleted job"})
	     return
	}

	// create the rating
	rating := EmployerJobRating{
		JobID: job.ID,
		Ratings: CompleteJobSchema.Rating,
		Description: CompleteJobSchema.Description,
	}

	// save the rating
	err = config.DB.Save(&rating).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message":"Unable to save rating"})
	     return
	}
	c.JSON(http.StatusOK, gin.H{"message":"employer rating succesful"})

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
	
	err = config.DB.Table("job_application_ratings").
		Joins("JOIN job_applications ON job_applications.id = job_application_ratings.job_application_id").
		Joins("JOIN jobs ON jobs.id = job_applications.job_id").
		Joins("JOIN employers ON employers.id = jobs.employer_id").
		Select(`job_application_ratings.created_at AS created_at,
				job_application_ratings.rating AS rating,
				job_application_ratings.description AS description,
				jobs.employer_id AS employer_id,
				employers.first_name AS first_name,
				employers.last_name AS last_name`).
		Where("job_applications.artisan_id = ?", id).
		Scan(&ratings).Error
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message":"error querying the database"})
	     return
	}

	c.JSON(http.StatusOK,gin.H{
		"ratings":ratings,
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

	var ratings []struct {
		CreatedAt     time.Time
		Rating        uint
		Description   string
		ArtisanID     uint
		BusinessName  string
	}
	
	err = config.DB.Table("employer_job_ratings").
		Joins("JOIN job_applications ON job_applications.id = employer_job_ratings.job_id").
		Joins("JOIN artisans ON artisans.id = job_applications.artisan_id").
		Joins("JOIN jobs ON jobs.id = job_applications.job_id").
		Select(`employer_job_ratings.created_at AS created_at,
				employer_job_ratings.ratings AS rating,
				employer_job_ratings.description AS description,
				artisans.id AS artisan_id,
				artisans.business_name AS business_name`).
		Where("jobs.employer_id = ?", id).
		Scan(&ratings).Error
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message":"error querying the database"})
	     return
	}

	c.JSON(http.StatusOK,gin.H{
		"ratings":ratings,
	})
	
}

func DeclineApplicant(c *gin.Context) {
	// get the job id
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
	// get the employer id
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

	var jobApplication JobApplication
	// change the application status and job status
	err = config.DB.
	Where("id = ?", id).
	First(&jobApplication).Error

	if err != nil {
	c.JSON(http.StatusBadRequest, gin.H{
	"error": "Failed to retrieve job data 2",
	})
	return
	}

	var job Job

	err = config.DB.
				Where("id = ?", jobApplication.JobID).
				First(&job).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to retrieve job data",
		})
		return
	}

	// check if the employer can work on the job
	if employerID != job.EmployerID {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Wrong employer details",
		})
		return
	}
	if jobApplication.ApplicationStatus == "hired" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Cannot deline an hired artisan",
		})
	}

	if jobApplication.ApplicationStatus == "declined" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Artisan already declined",
		})
	}

	jobApplication.ApplicationStatus = "declined"
	err = config.DB.Save(&jobApplication).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Unable to save decline artisan",
		})
	}

	c.JSON(http.StatusOK,gin.H{
		"message":"application declined successfully",
	})
		
}

func ShortlistApplicant(c *gin.Context) {
	// get the job id
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
	// get the employer id
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

	var jobApplication JobApplication
	// change the application status and job status
	err = config.DB.
	Where("id = ?", id).
	First(&jobApplication).Error

	if err != nil {
	c.JSON(http.StatusBadRequest, gin.H{
	"error": "Failed to retrieve job data 2",
	})
	return
	}

	var job Job

	err = config.DB.
				Where("id = ?", jobApplication.JobID).
				First(&job).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to retrieve job data",
		})
		return
	}

	// check if the employer can work on the job
	if employerID != job.EmployerID {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Wrong employer details",
		})
		return
	}
	if jobApplication.ApplicationStatus == "hired" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Cannot deline an hired artisan",
		})
	}

	if jobApplication.ApplicationStatus == "listed" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Artisan already declined",
		})
	}

	jobApplication.ApplicationStatus = "listed"
	err = config.DB.Save(&jobApplication).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Unable to save decline artisan",
		})
	}
	c.JSON(http.StatusOK,gin.H{
		"message":"application listed successfully",
	})
		
}

func GetEmployerJobStats(c *gin.Context) {
	// get the employer id 


	// variables needed
}