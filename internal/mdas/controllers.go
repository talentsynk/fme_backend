package mda

import (
	"errors"
	"fme_backend/internal/config"
	"fmt"
	"net/http"
    "github.com/gin-gonic/gin"
	"gorm.io/gorm"
)


func CreateMda(c *gin.Context) {
	fmt.Println("controller started")
    if  c.BindJSON(&MdaCreateSchema) != nil{
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Failed to read request body",
        })
        return
    }

    userID, exists := c.Get("userID") // Make sure this key matches what you set in the middleware
    if !exists {
        c.JSON(http.StatusBadRequest, gin.H{"error": "UserID not found in context"})
        return
    }

    mda := Mda{
        Name:       MdaCreateSchema.Name,
        AgencyCode: MdaCreateSchema.AgencyCode,
        IsActive:   true,
        UserID:     userID.(uint),
    }

	fmt.Println(mda)
    if result := config.DB.Create(&mda); result.Error != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Mda created successfully"})
}



func GetMdas(c *gin.Context){
    fmt.Println("controller started")

    var mdas []Mda
    if result := config.DB.Find(&mdas); result.Error != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"mdas":mdas})
}


func GetMdaByID(c *gin.Context){
    id := c.Param("id")
    if id == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error":"Mda ID is required"})
        return
    }

    var mda Mda
    if err := config.DB.First(&mda, id).Error; err != nil{
        if errors.Is(err, gorm.ErrRecordNotFound) {
            c.JSON(http.StatusNotFound, gin.H{"error":"Mda not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error":"Internal Server Error"})
        }
        return
    }

    c.JSON(http.StatusOK, mda)
}


func TotalNumberOfMda(c *gin.Context){
     fmt.Println("Get Total Number Of Mda")
     var count int64

   if result := config.DB.Model(&Mda{}).Count(&count); result.Error != nil{
     c.JSON(http.StatusInternalServerError, gin.H{"error":result.Error.Error()})
     return
   }

   c.JSON(http.StatusOK, gin.H{"total_count":count})
}


func TotalNumberOfActiveMda(c *gin.Context){
    fmt.Println("Get Total Number Of Active Mda")
    var count int64

    if result := config.DB.Model(&Mda{}).Where("is_active = ?", true).Count(&count); result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"total_is_active_count":count})
}


func TotalNumberOfInActiveMda(c *gin.Context){
    fmt.Println("Get Total Number Of InActive Mda")
    var count int64

    if result := config.DB.Model(&Mda{}).Where("is_active = ?", false).Count(&count); result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"total_in_active_count":count})
}





func SearchMda(c *gin.Context) {
    query := c.Query("query")
    if query == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
        return
    }

    var mdasearch []Mda
    if err := config.DB.Where("name LIKE ? OR agency_code LIKE ?", "%"+query+"%", "%"+query+"%").Find(&mdasearch).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
        return
    }

    if len(mdasearch) == 0 {
        c.JSON(http.StatusOK, gin.H{"message": "No matching mdas found"})
        return
    }

    c.JSON(http.StatusOK, mdasearch)
}


func UpdateMda(c *gin.Context) {
    id := c.Param("id")
    if id == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Mda ID is required"})
        return
    }

    var mda Mda
    if err := config.DB.First(&mda, id).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            c.JSON(http.StatusNotFound, gin.H{"error": "Mda not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
        }
        return
    }

    if err := c.BindJSON(&mda); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
        return
    }

    if err := config.DB.Save(&mda).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update Mda"})
        return
    }

    c.JSON(http.StatusOK, mda)
}




func SuspendMda( c *gin.Context){
    id := c.Param("id")
    if id == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error":"Mda ID is required"})
        return
    }

    var mda Mda
 if err := config.DB.First(&mda, id).Error; err != nil {
    if errors.Is(err, gorm.ErrRecordNotFound){
        c.JSON(http.StatusNotFound, gin.H{"error":"Mda not  found"})
    }else {
        c.JSON(http.StatusInternalServerError, gin.H{"error":"Internal Server Error"})
    }
    return
 }

    mda.IsActive = false

    if err := config.DB.Save(&mda).Error; err != nil{
        c.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to suspend Mda"})
    }
    c.JSON(200, gin.H{"message":"Mda Suspend successfully"})

}




func ActivateMda( c *gin.Context){
    id := c.Param("id")
    if id == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error":"Mda ID is required"})
        return
    }

    var mda Mda
 if err := config.DB.First(&mda, id).Error; err != nil {
    if errors.Is(err, gorm.ErrRecordNotFound){
        c.JSON(http.StatusNotFound, gin.H{"error":"Mda not  found"})
    }else {
        c.JSON(http.StatusInternalServerError, gin.H{"error":"Internal Server Error"})
    }
    return
 }

    mda.IsActive = true

    if err := config.DB.Save(&mda).Error; err != nil{
        c.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to activate Mda"})
    }
    c.JSON(200, gin.H{"message":"Mda Activate successfully"})

}