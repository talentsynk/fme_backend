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

    fme := r.Group("/fme")
    fme.POST("/create", myuser.CreateFme)

    r.Run(":8080")
}
