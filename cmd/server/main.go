package main

import (
	"fme_backend/internal/config"
	"fme_backend/internal/course"
	"fme_backend/internal/dashboard"
	mda "fme_backend/internal/mdas"
	middleware "fme_backend/internal/middlewares"
	stc "fme_backend/internal/stc"
	"fme_backend/internal/student"
	myuser "fme_backend/internal/user"
	employer "fme_backend/internal/employers"

     job "fme_backend/internal/jobs"
	"time"
    "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {
	config.ConnectToDb()
}

func main() {
	r := gin.Default()

	config := cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	r.Use(cors.New(config))

	usergroup := r.Group("/user")
	usergroup.POST("/createfme", myuser.CreateFmeUser)
	usergroup.GET("/activate/:id", middleware.RequireAuth, myuser.ActivateUser)
	usergroup.GET("/suspend/:id", middleware.RequireAuth, myuser.SuspendUser)
	usergroup.POST("/login", myuser.Login)
	usergroup.POST("/otp/request", myuser.RequestOtp)
	usergroup.POST("/otp/verify", myuser.VerifyOtp)
	usergroup.POST("/changepassword", myuser.ChangePassword)

	mdagroup := r.Group("/mda")
	mdagroup.POST("/create-mda", middleware.RequireFme, mda.CreateMda)
	mdagroup.GET("/get-all-mdas", middleware.RequireFme, mda.GetAllMdas)
	mdagroup.GET("/search-mda", middleware.RequireAuth, mda.SearchMda)
	mdagroup.PATCH("/update-mda/:id", middleware.RequireAuth, mda.UpdateMda)
	mdagroup.GET("/get-mda/:id", middleware.RequireAuth, mda.GetMdaByID)
	mdagroup.GET("/total-mda", middleware.RequireAuth, mda.MdaTotal)
  mdagroup.GET("/get-ascending-mda", middleware.RequireAuth, mda.FilterMdaAscending)
  mdagroup.GET("/get-descending-mda", middleware.RequireAuth, mda.FilterMdaDescending)
	mdagroup.GET("/filter-by-state",middleware.RequireAuth, mda.FilterMdaByState )
  mdagroup.POST("/suspend-mda/:id", middleware.RequireAuth, mda.SuspendMda)
	mdagroup.POST("/activate-mda/:id",middleware.RequireAuth, mda.ActivateMda)
  mdagroup.GET("/profile",middleware.RequireMda,mda.GetMdaProfile)
mdagroup.GET("/download-csv", middleware.RequireFme, mda.DownloadMdaCsv)

    
   stcgroup := r.Group("/stc")
	 stcgroup.POST("/create-stc",middleware.RequireFme, stc.CreateFmeStc)
	 stcgroup.POST("/create-mda-stc", middleware.RequireMda, stc.CreateMdaStc)
	 stcgroup.GET("/get-all-stc", middleware.RequireAuth, stc.GetStc)
	 stcgroup.GET("/search-stc", middleware.RequireAuth, stc.SearchStc)
	 stcgroup.PATCH("/update-stc/:id", middleware.RequireAuth, stc.UpdateStc)
	 stcgroup.GET("/get-stc/:id", middleware.RequireAuth, stc.GetStcByID)
	 stcgroup.GET("/get-total-count", middleware.RequireAuth,stc.StcTotal)	
   stcgroup.GET("/get-ascending-stc", middleware.RequireAuth, stc.FilterStcAscending)
   stcgroup.GET("/get-descending-stc", middleware.RequireAuth, stc.FilterStcDescending)
	 stcgroup.POST("/suspend-stc/:id", middleware.RequireAuth, stc.SuspendStc)
	 stcgroup.POST("/activate-stc/:id",middleware.RequireAuth, stc.ActivateStc)
   stcgroup.GET("/filter-by-state", middleware.RequireAuth, stc.FilterStcByState)
	 stcgroup.GET("/get-all-mda-stc", middleware.RequireMda, stc.GetAllMdaStc)
	 stcgroup.GET("/get-mda-stc/:id", middleware.RequireMda, stc.GetMdaStcByID)
   stcgroup.GET("/get-mda-total", middleware.RequireMda, stc.StcMdaTotal)
	 stcgroup.GET("/profile",middleware.RequireStc,stc.GetStcProfile)


  studentgroup:= r.Group("/student")
	studentgroup.POST("/create-fme",middleware.RequireFme,student.CreateFmeStudent)
	studentgroup.POST("/create-mda",middleware.RequireMda,student.CreateMdaStudent)


	studentgroup.POST("/create-stc",middleware.RequireStc,student.CreateStcStudent)
	studentgroup.POST("/create-mda-csv",middleware.RequireMda,student.CreateMdaStudentFromCsv)
	studentgroup.POST("/create-stc-csv",middleware.RequireStc,student.CreateStcStudentFromCsv)
	studentgroup.POST("/create-fme-csv",middleware.RequireFme,student.CreateFmeStudentFromCsv)
	
	studentgroup.GET("/all",middleware.RequireAuth,student.GetAllStudents)
	studentgroup.GET("/:id",middleware.RequireAuth,student.GetStudent)
	studentgroup.GET("/total-info",middleware.RequireAuth,student.GetTotalStudentInfo)
	studentgroup.GET("/download-csv",middleware.RequireAuth,student.DownloadStudentsCsv)

	categorygroup := r.Group("category")
	categorygroup.POST("/create", middleware.RequireFme, course.CreateCategory)
	categorygroup.GET("/:id", middleware.RequireAuth, course.GetCategory)
	categorygroup.GET("/all", middleware.RequireAuth, course.GetAllCategories)

	coursegroup := r.Group("/course")
	coursegroup.POST("/create", middleware.RequireFme, course.CreateCourse)
	coursegroup.GET("/:id", middleware.RequireAuth, course.GetCourse)
	coursegroup.GET("/all", middleware.RequireAuth, course.GetAllCourses)

	dashgroup := r.Group("/dashboard")
	dashgroup.GET("/summary", middleware.RequireAuth, dashboard.GetDashSummary)
	dashgroup.GET("/course-percentage", middleware.RequireAuth, dashboard.GetStudentPercentPerCourse)
	dashgroup.GET("/top-5-mda", middleware.RequireAuth, dashboard.GetTopMda)
	dashgroup.GET("/top-5-stc", middleware.RequireAuth, dashboard.GetTopStc)


	employergroup  := r.Group("/employer")
	employergroup.POST("/create-employer", employer.CreateEmployer)
	employergroup.GET("/get-employer", middleware.RequireEmployer, employer.GetEmployer)
	employergroup.GET("/:id", middleware.RequireEmployer, employer.GetEmployerByID)
  

	jobgroup := r.Group("/job")
	jobgroup.POST("/create-job", middleware.RequireEmployer, job.CreateJob)
	jobgroup.GET("/get-jobs", middleware.RequireAuth, job.GetAllJobs)
	jobgroup.GET("/get-job/:id", middleware.RequireAuth,  job.GetJobID)
	jobgroup.GET("/get-latest-job", middleware.RequireAuth, job.GetLatestJobs)
	jobgroup.GET("/search-job", middleware.RequireAuth, job.SearchJob)
	jobgroup.POST("/applied", middleware.RequireStudent,  job.ApplyForJob)
	r.Run(":8000")
}
