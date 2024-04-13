package mda

import (
	"errors"
	"fme_backend/internal/config"
    myuser "fme_backend/internal/user"
	"fmt"
	"net/http"
	"sort"
    "strconv"
    "fme_backend/internal/utilities"
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

    stateOfOperation, result := utilities.ValidateState( MdaCreateSchema.StateOfOperation)
    if !result {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Incorrect state of origin",
        })
        return
    }
   // Transaction Handling
   tx := config.DB.Begin() 
   result, message, newUserID := myuser.CreateMdaUser(tx, MdaCreateSchema.PhoneNumber, MdaCreateSchema.Email, "dfcv")
   if !result {
       tx.Rollback() // Rollback if user creation fails
       c.JSON(http.StatusBadRequest, gin.H{
           "error": message,
       })
       return
   }

    mda := Mda{
        RegisterName:  MdaCreateSchema.RegisterName,
        Address: MdaCreateSchema.Address,
        StateOfOperation: stateOfOperation,
        UserID: newUserID,
    }

	fmt.Println(mda)
    mdaresult := tx.Create(&mda) 
    if mdaresult.Error != nil {
        tx.Rollback() 
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Failed to create User",
        })
        return
    }

    tx.Commit() 
    c.JSON(http.StatusOK, gin.H{
        "message": "Mda created successfully",
    })
}



func GetAllMdas(c *gin.Context){
    fmt.Println("controller started")
    limitStr := c.Query("limit")
    pageStr  := c.Query("page")

    limit, err := strconv.Atoi(limitStr)
    if err != nil || limit <= 0{
        limit = 10 
    }

    page, err := strconv.Atoi(pageStr)
    if err != nil || page <= 0 {
        page = 1
    }

    offset := (page - 1) * limit
 
    var mdas []Mda
    if result := config.DB.Limit(limit).Offset(offset).Find(&mdas); result.Error != nil {
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
    if err := config.DB.Where("register_name LIKE ? OR email LIKE ? OR address LIKE ? OR state_of_operation LIKE ?", "%"+query+"%", "%"+query+"%" , "%"+query+"%","%"+query+"%").Find(&mdasearch).Error; err != nil {
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






func FilterMdaDescending(c *gin.Context) {
    var mdas []Mda
    if result := config.DB.Find(&mdas); result.Error != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
        return
    }

    // Sort MDAs by name in ascending order
    sort.Slice(mdas, func(i, j int) bool {
        return mdas[i].RegisterName < mdas[j].RegisterName
    })

    // Send the sorted MDAs as the response
    c.JSON(http.StatusOK, gin.H{"mdas": mdas})
}



func FilterMdaAscending(c *gin.Context) {
    var mdas []Mda
    if result := config.DB.Find(&mdas); result.Error != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
        return
    }

    // Sort MDAs by name in ascending order
    sort.Slice(mdas, func(i, j int) bool {
        return mdas[i].RegisterName > mdas[j].RegisterName
    })

    // Send the sorted MDAs as the response
    c.JSON(http.StatusOK, gin.H{"mdas": mdas})
}


func FilterMdaByState(c *gin.Context){
    StateOfOperation := c.Query("StateOfOperation")
    if StateOfOperation == ""{
        c.JSON(http.StatusBadRequest, gin.H{
            "error":"State parameter is required",
        })
        return
    }

    var mdas []Mda
    tx := config.DB.Where("StateOfOperation = ?",StateOfOperation).Find(&mdas)
    if tx.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error":"Failed to retrieve STCs",
        })
         return
    }
    c.JSON(http.StatusOK, gin.H{"mdas":mdas})
}

// func SuspendMda(c *gin.Context) {
//     id := c.Param("id")
//     if id == "" {
//         c.JSON(http.StatusBadRequest, gin.H{"error": "Mda ID is required"})
//         return
//     }

//     var mda Mda
//     if err := config.DB.First(&mda, id).Error; err != nil {
//         if errors.Is(err, gorm.ErrRecordNotFound) {
//             c.JSON(http.StatusNotFound, gin.H{"error": "Mda not found"})
//         } else {
//             c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
//         }
//         return
//     }

//     mda.IsActive = false // Assuming "suspend" means setting IsActive to false

//     if err := config.DB.Save(&mda).Error; err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to suspend Mda"})
//         return
//     }

//     c.JSON(http.StatusOK, gin.H{"suspend":"Mda Suspended"})
// }



// func ActivateMda(c *gin.Context) {
//     id := c.Param("id")
//     if id == "" {
//         c.JSON(http.StatusBadRequest, gin.H{"error": "Mda ID is required"})
//         return
//     }

//     var mda Mda
//     if err := config.DB.First(&mda, id).Error; err != nil {
//         if errors.Is(err, gorm.ErrRecordNotFound) {
//             c.JSON(http.StatusNotFound, gin.H{"error": "Mda not found"})
//         } else {
//             c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
//         }
//         return
//     }

//     mda.IsActive = true // Assuming "activate" means setting IsActive to true

//     if err := config.DB.Save(&mda).Error; err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to activate Mda"})
//         return
//     }

//     c.JSON(http.StatusOK, gin.H{"Activate":"Mda Activated"})
// }