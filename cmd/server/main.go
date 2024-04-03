package main

import (
	"fme_backend/internal/config"
	mda "fme_backend/internal/mdas"
	middleware "fme_backend/internal/middlewares"
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
	mdagroup.POST("/create-mda", middleware.RequireAuth,  mda.CreateMda)
	mdagroup.GET("/get-all-mdas", middleware.RequireAuth, mda.GetMdas)
	mdagroup.GET("/search-mda", middleware.RequireAuth, mda.SearchMda)
	mdagroup.PATCH("/update-mda/:id", middleware.RequireAuth, mda.UpdateMda)
	mdagroup.GET("/get-mda/:id", middleware.RequireAuth, mda.GetMdaByID)
	mdagroup.GET("/get-total-number-mda", middleware.RequireAuth, mda.TotalNumberOfMda)
	mdagroup.GET("/get-total-isactive-mda", middleware.RequireAuth, mda.TotalNumberOfActiveMda)
	mdagroup.GET("/get-total-inactive-mda", middleware.RequireAuth, mda.TotalNumberOfActiveMda)

    r.Run(":8000")
}
