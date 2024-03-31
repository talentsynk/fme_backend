package main

import (
	"fme_backend/internal/config"
	myuser "fme_backend/internal/user"
	 mda "fme_backend/internal/mdas"
	 middleware "fme_backend/internal/middlewares"


	"github.com/gin-gonic/gin"
)

func init() {
    config.ConnectToDb()
}

func main() {
	r := gin.Default()

    usergroup := r.Group("/user")
      usergroup.POST("/createfme", myuser.CreateFmeUser)
	    usergroup.PATCH("/deactivate",middleware.RequireAuth, myuser.DeactivateUser)
	    usergroup.PATCH("/activate/:id",middleware.RequireAuth, myuser.ActivateUser)
	    usergroup.PATCH("/suspend/:id",middleware.RequireAuth, myuser.SuspendUser)
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
	mdagroup.PATCH("/suspend/:id", middleware.RequireAuth, mda.SuspendMda)
	mdagroup.PATCH("/activate/:id", middleware.RequireAuth, mda.ActivateMda)

    r.Run(":8080")
}
