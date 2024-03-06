package main

import (
	"os"
	"log"
    "fme_backend/config"
	"github.com/joho/godotenv"
	"github.com/gin-gonic/gin"
  route "fme_backend/user/routes"
	
)

func init() {
    config.ConnectToDb()
}

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	} 

	router := gin.New()

	route.AuthRoutes(router)
	route.UserRoutes(router)

	router.Run(":" + port)
	
}