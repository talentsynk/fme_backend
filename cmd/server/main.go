package main

import (
	"fme_backend/internal/config"
	mda "fme_backend/internal/mdas"
	 "fme_backend/internal/course"
	
	middleware "fme_backend/internal/middlewares"
	"fme_backend/internal/student"
	stc "fme_backend/internal/stc"
	myuser "fme_backend/internal/user"

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
		AllowAllOrigins: true,
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:   []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge: 12 * time.Hour,
	  }
	r.Use(cors.New(config))

    usergroup := r.Group("/user")
    usergroup.POST("/createfme", myuser.CreateFmeUser)
	usergroup.GET("/activate/:id",middleware.RequireAuth, myuser.ActivateUser)
	usergroup.GET("/suspend/:id",middleware.RequireAuth, myuser.SuspendUser)
	usergroup.POST("/login", myuser.Login)
	usergroup.POST("/otp/request",myuser.RequestOtp)
	usergroup.POST("/otp/verify" ,myuser.VerifyOtp)
	usergroup.POST("/changepassword" ,myuser.ChangePassword)

    mdagroup := r.Group("/mda")
	mdagroup.POST("/create-mda",middleware.RequireFme, mda.CreateMda)
	mdagroup.GET("/get-all-mdas", middleware.RequireFme, mda.GetAllMdas)
	mdagroup.GET("/search-mda",middleware.RequireAuth, mda.SearchMda)
	mdagroup.PATCH("/update-mda/:id",middleware.RequireAuth, mda.UpdateMda)
	mdagroup.GET("/get-mda/:id",middleware.RequireAuth, mda.GetMdaByID)
	mdagroup.GET("/total-mda", middleware.RequireAuth, mda.MdaTotal)
    mdagroup.GET("/get-ascending-mda", middleware.RequireAuth, mda.FilterMdaAscending)
    mdagroup.GET("/get-descending-mda", middleware.RequireAuth, mda.FilterMdaDescending)
	mdagroup.GET("/filter-by-state",middleware.RequireAuth, mda.FilterMdaByState )
    // mdagroup.PATCH("/suspend-mda/:id", middleware.RequireAuth, mda.SuspendMda)
	// mdagroup.PATCH("/activate-mda/:id",middleware.RequireAuth, mda.ActivateMda)
    
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
	 stcgroup.PATCH("/suspend-stc/:id", middleware.RequireAuth, stc.SuspendStc)
	 stcgroup.PATCH("/activate-stc/:id",middleware.RequireAuth, stc.ActivateStc)
     stcgroup.GET("/filter-by-state", middleware.RequireAuth, stc.FilterStcByState)
	 stcgroup.GET("/get-mda-stc",middleware.RequireMda, stc.GetTotalStcsByMdaID)
	 


  	studentgroup:= r.Group("/student")
	studentgroup.POST("/create-fme",student.CreateFmeStudent)
	studentgroup.POST("/create-mda",middleware.RequireMda,student.CreateMdaStudent)
	studentgroup.POST("/create-stc",middleware.RequireStc,student.CreateStcStudent)
	studentgroup.GET("/all-fme",middleware.RequireFme,student.GetAllStudents)
	studentgroup.GET("/fme/:id",middleware.RequireFme,student.GetStudent)

	categorygroup := r.Group("category")
	categorygroup.POST("/create",middleware.RequireFme,course.CreateCategory)
	categorygroup.GET("/:id",middleware.RequireAuth,course.GetCategory)
	categorygroup.GET("/all",middleware.RequireAuth,course.GetAllCategories)

	coursegroup := r.Group("/course")
	coursegroup.POST("/create",middleware.RequireFme,course.CreateCourse)
	coursegroup.GET("/:id", middleware.RequireAuth, course.GetCourse)
	coursegroup.GET("/all",middleware.RequireAuth,course.GetAllCourses)

	


    r.Run(":8000")
}
