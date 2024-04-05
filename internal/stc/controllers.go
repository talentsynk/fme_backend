package stc

import (
	"fme_backend/internal/config"
	"fmt"
	"net/http"
    "github.com/gin-gonic/gin"
)



func CreateStc(c *gin.Context){
	fmt.Println("controller started")
	if c.BindJSON(&StcCreateSchema) != nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"error":"Failed to read request body",
		})
		return
	}
	 userID, exists := c.Get("userID")
	 if !exists{
		c.JSON(http.StatusBadRequest, gin.H{"error":"UserID not found in Context"})
	 }

	 mdaID, exists := c.Get("mdaID")
	 if !exists{
		c.JSON(http.StatusBadRequest, gin.H{"error":"mdaID not found in Context"})
	 }

	stc := Stc{
		Ownership:         StcCreateSchema.Ownership, 
		CentreCode:        StcCreateSchema.CentreCode,
		Name:              StcCreateSchema.Name,
		LocalGovernment:   StcCreateSchema.LocalGovernment,
		State: 			   StcCreateSchema.State,
		isOperational:     true,
		CertificateOfOperationURL: StcCreateSchema.CertificateOfOperationURL,
		MdaID: mdaID.(uint),
		UserID: userID.(uint),
	}

	fmt.Println(stc)
	if result := config.DB.Create(&stc); result.Error != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error":result.Error.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"message":"Stc created successfully"})
}