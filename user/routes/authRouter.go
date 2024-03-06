package route

import (

  controller   "fme_backend/user/controllers"
	"github.com/gin-gonic/gin"
)



func AuthRoutes(incomingRoutes *gin.Engine){
incomingRoutes.POST("/user/login", controller.Login)
incomingRoutes.POST("/user/create-fme", controller.CreateFmeUser)

}