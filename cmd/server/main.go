package main

import (
	"fme_backend/internal/config"
	myuser "fme_backend/internal/user"

	"github.com/gin-gonic/gin"
)

func init() {
    config.ConnectToDb()
}

func main() {
	r := gin.Default()

    usergroup := r.Group("/user")
    usergroup.POST("/createfme", myuser.CreateFmeUser)
	usergroup.POST("/login", myuser.Login)
	usergroup.GET("/activate/:id",myuser.RequireAuth, myuser.ActivateUser)
	usergroup.GET("/suspend/:id",myuser.RequireAuth, myuser.SuspendUser)
	usergroup.POST("/otp/request" ,myuser.RequestOtp)
	usergroup.POST("/otp/verify" ,myuser.VerifyOtp)
	usergroup.POST("/changepassword" ,myuser.ChangePassword)


    r.Run(":8080")
}
