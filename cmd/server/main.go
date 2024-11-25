package main

import (
	"fme_backend/internal/artisans"
	"fme_backend/internal/config"
	"fme_backend/internal/course"
	"fme_backend/internal/dashboard"
	employer "fme_backend/internal/employers"
	mda "fme_backend/internal/mdas"
	middleware "fme_backend/internal/middlewares"
	"fme_backend/internal/scripts"
	stc "fme_backend/internal/stc"
	"fme_backend/internal/student"
	myuser "fme_backend/internal/user"

	job "fme_backend/internal/jobs"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {
	config.ConnectToDb()
}

func main() {
	scripts.CreateFmeAtStart()
	r := gin.Default()

	corsConfig := cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	r.Use(cors.New(corsConfig))
	r.Use(middleware.RequestLogger(config.DB))

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
	mdagroup.PATCH("/edit/:id", middleware.RequireFme, mda.EditMdaData)

    
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
	stcgroup.GET("/download-csv",middleware.RequireAuth, stc.DownloadStcCsv)
	stcgroup.PATCH("/edit/:id", middleware.RequireAuth, stc.EditStc)


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
	studentgroup.POST("/graduate-student/:id", middleware.RequireAuth, student.GraduateStudent)
	studentgroup.GET(("/disabilities"),student.GetDisabilityList)
	studentgroup.PATCH("/edit/:id", middleware.RequireAuth, student.EditStudent)


	categorygroup := r.Group("category")
	categorygroup.POST("/create", middleware.RequireFme, course.CreateCategory)
	categorygroup.GET("/:id", middleware.RequireAuth, course.GetCategory)
	categorygroup.GET("/all", middleware.RequireAuth, course.GetAllCategories)

	coursegroup := r.Group("/course")
	coursegroup.POST("/create", middleware.RequireFme, course.CreateCourse)
	coursegroup.GET("/:id", middleware.RequireAuth, course.GetCourse)
	coursegroup.GET("/all", middleware.RequireAuth, course.GetAllCourses)
	coursegroup.GET("/details/:id",middleware.RequireAuth,course.GetCourseDetails)
	coursegroup.GET("/:id/top-5-mda",middleware.RequireAuth,course.GetTopMDAsByStudentCountForCourse)
	coursegroup.GET("/:id/top-5-stc",middleware.RequireAuth,course.GetTopSTCsByStudentCountForCourse)


	dashgroup := r.Group("/dashboard")
	dashgroup.GET("/summary", middleware.RequireAuth, dashboard.GetDashSummary)
	dashgroup.GET("/course-percentage", middleware.RequireAuth, dashboard.GetStudentPercentPerCourse)
	dashgroup.GET("/top-5-mda", middleware.RequireAuth, dashboard.GetTopMda)
	dashgroup.GET("/top-5-stc", middleware.RequireAuth, dashboard.GetTopStc)


	employergroup  := r.Group("/employer")
	employergroup.POST("/create-employer", employer.CreateEmployer)
	employergroup.GET("/get-employer", middleware.RequireEmployer, employer.GetEmployer)
	employergroup.GET("/:id", middleware.RequireAuth, employer.GetEmployerByID)
	employergroup.GET("/get-all-employer", middleware.RequireEmployer, employer.GetAllEmployer)
	employergroup.GET("/dash-stats", middleware.RequireEmployer, employer.EmployerDashboard)
	employergroup.GET("/profile-stats", middleware.RequireEmployer, employer.GetEmployerProfileStats)
	employergroup.GET("/profile-stats/:id", middleware.RequireArtisan, employer.GetEmployerProfileStatsByArtisan)
	employergroup.GET("/ratings/:id",middleware.RequireAuth,employer.GetEmployerRating)
	employergroup.GET("/jobs/:id",middleware.RequireArtisan,employer.GetAllEmployerJobs)
	employergroup.GET("/similar/:id",middleware.RequireArtisan,employer.GetSimilarEmployerDetails)



  

	jobgroup := r.Group("/job")
	jobgroup.POST("/create-job", middleware.RequireEmployer, job.CreateJob)
	jobgroup.GET("/all", middleware.RequireArtisan, job.GetAllJobs)
	jobgroup.GET("/:id", middleware.RequireArtisan,  job.GetJobID)
	jobgroup.POST("/save/:id", middleware.RequireArtisan, job.SaveNewJob)
	jobgroup.POST("/apply/:id", middleware.RequireArtisan, job.ApplyForJob)
	jobgroup.GET("/applied-jobs", middleware.RequireArtisan, job.GetAppliedJobs)
	jobgroup.GET("/saved-jobs", middleware.RequireArtisan, job.GetSavedJobs)
	jobgroup.POST("/hire", middleware.RequireEmployer, job.HireArtisan)
	jobgroup.GET("/applicants/:id", middleware.RequireEmployer, job.ViewApplicants)
	jobgroup.POST("/rate/:id", middleware.RequireEmployer, job.CompleteJob)
	jobgroup.GET("/my-jobs", middleware.RequireEmployer, job.GetMyJobs)
	jobgroup.GET("/close/:id",middleware.RequireEmployer,job.CloseJob)	
	jobgroup.GET("/open/:id",middleware.RequireEmployer,job.OpenJob)
	jobgroup.POST("/rate/employer/:id",middleware.RequireArtisan,job.RateEmployer)
	jobgroup.GET("/profile/artisan/:id",middleware.RequireEmployer, job.GetArtisanJobProfile)
	jobgroup.GET("/applicants/decline/:id",middleware.RequireEmployer,job.DeclineApplicant)
	jobgroup.GET("/applicants/short-list/:id",middleware.RequireEmployer,job.ShortlistApplicant)
	jobgroup.GET("/similar/:id",middleware.RequireArtisan,job.GetSimilarJobs)
	jobgroup.POST("/search",middleware.RequireArtisan,job.SearchForJob)
	jobgroup.GET("/general/:id", middleware.RequireAuth,  job.GetJobGeneral)



	artisanGroup := r.Group("/artisan")
	artisanGroup.GET("/all", middleware.RequireAuth, artisans.GetAllArtisans)
	artisanGroup.GET("/:id", middleware.RequireAuth, artisans.GetArtisanProfile)
	artisanGroup.GET("/job-stats",middleware.RequireArtisan,artisans.GetArtisanJobStat)
	artisanGroup.GET("/profile-stats",middleware.RequireArtisan,artisans.GetArtisanProfileStat)
	artisanGroup.GET("/profile-stats/:id",middleware.RequireEmployer,artisans.GetArtisanProfileStatByEmployer)
	artisanGroup.GET("/ratings/:id",middleware.RequireAuth,artisans.GetArtisanRating)
	artisanGroup.GET("/contact/:id", middleware.RequireAuth, artisans.GetContactDetails)
	artisanGroup.GET("/me", middleware.RequireArtisan, artisans.GetMyDetails)
	artisanGroup.GET("/download-data",middleware.RequireFme,artisans.DownloadArtisanData)
	




	r.Run(":8080")
}
