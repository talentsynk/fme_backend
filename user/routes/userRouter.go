package route


import (
 controller	"fme_backend/user/controllers"
  	"github.com/gin-gonic/gin"
 middleware	 "fme_backend/user/middlewares"
)



func UserRoutes(incomingRoutes *gin.Engine){
incomingRoutes.Use(middleware.RequireAuth)
incomingRoutes.PATCH("/user/deactivate",controller.DeactivateUser)
incomingRoutes.PATCH("/user/activate", controller.ActivateUser)
incomingRoutes.PATCH("/user/suspend", controller.SuspendUser)
incomingRoutes.GET("/user/request-otp",controller.RequestOtp)
}