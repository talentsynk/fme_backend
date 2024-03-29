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
	usergroup.PATCH("/activate",middleware.RequireAuth, myuser.ActivateUser)
	usergroup.PATCH("/suspend",middleware.RequireAuth, myuser.SuspendUser)
	usergroup.POST("/login", myuser.Login)
	usergroup.GET("/requestotp",middleware.RequireAuth,myuser.RequestOtp)

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
