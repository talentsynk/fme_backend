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
	usergroup.PATCH("/deactivate",myuser.RequireAuth, myuser.DeactivateUser)
	usergroup.PATCH("/activate",myuser.RequireAuth, myuser.ActivateUser)
	usergroup.PATCH("/suspend",myuser.RequireAuth, myuser.SuspendUser)
	usergroup.POST("/login", myuser.Login)
	usergroup.GET("/requestotp",myuser.RequireAuth,myuser.RequestOtp)


	

    r.Run(":8080")
}
