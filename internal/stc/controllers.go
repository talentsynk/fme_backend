package stc

import (
	"errors"
	"fme_backend/internal/config"
    myuser "fme_backend/internal/user"
    "fme_backend/internal/utilities"
	"fmt"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)



func CreateFmeStc(c *gin.Context){
	fmt.Println("controller started")
	if c.BindJSON(&StcCreateSchema) != nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"error":"Failed to read request body",
		})
		return
	}
    state, result := utilities.ValidateState(StcCreateSchema.State)
    if !result {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Incorrect state of residence",
        })
        return
    }
  // Transaction Handling
  tx := config.DB.Begin() // Begin a transaction

  // Create user with transaction
  result, message, newUserID := myuser.CreateStcUser(tx, StcCreateSchema.Phone, StcCreateSchema.Email, "dfcv")
  if !result {
      tx.Rollback() // Rollback if user creation fails
      c.JSON(http.StatusBadRequest, gin.H{
          "error": message,
      })
      return
  }
  	stc := Stc{
		Ownership:         StcCreateSchema.Ownership, 
		CentreCode:        StcCreateSchema.CentreCode,
		Name:              StcCreateSchema.Name,
		LocalGovernment:   StcCreateSchema.LocalGovernment,
		State: 			   state,
		isOperational:     true,
		CertificateOfOperationURL: StcCreateSchema.CertificateOfOperationURL,
        UserID: newUserID,
        Fmestc: true,
    }

	fmt.Println(stc)
    stcresult := tx.Create(&stc) // Create student within transaction
    if stcresult.Error != nil {
        tx.Rollback() // Rollback if student creation fails
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Failed to create User",
        })
        return
    }

    tx.Commit() // Commit the transaction if both creations are successful
    // Success response
    c.JSON(http.StatusOK, gin.H{
        "message": "Student created successfully",
    })
}


func CreateMdaStc(c *gin.Context){
    // retrieve the mda id
    mdaIDStr,exists := c.Get("mdaID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Problem with the authorization token",
		})
		return
	}
    mdaID,ok := mdaIDStr.(uint)
	if !ok{
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Problem with the authorization token",
		})
		return
	}
    if c.Bind(&StcCreateSchema) != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Failed to read request body",
        })
        return
    }

    state, result := utilities.ValidateState(StcCreateSchema.State)
    if !result {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Incorrect state of residence",
        })
        return
    }
     // Transaction Handling
  tx := config.DB.Begin() // Begin a transaction

  // Create user with transaction
  result, message, newUserID := myuser.CreateStcUser(tx, StcCreateSchema.Phone, StcCreateSchema.Email, "dfcv")
  if !result {
      tx.Rollback() // Rollback if user creation fails
      c.JSON(http.StatusBadRequest, gin.H{
          "error": message,
      })
      return
  }

  stc := Stc{
    Ownership:         StcCreateSchema.Ownership, 
    CentreCode:        StcCreateSchema.CentreCode,
    Name:              StcCreateSchema.Name,
    LocalGovernment:   StcCreateSchema.LocalGovernment,
    State: 			   state,
    isOperational:     true,
    CertificateOfOperationURL: StcCreateSchema.CertificateOfOperationURL,
    UserID: newUserID,
    MdaID: mdaID,
}
fmt.Println(stc)
stcresult := tx.Create(&stc) // Create student within transaction
if stcresult.Error != nil {
    tx.Rollback() // Rollback if student creation fails
    c.JSON(http.StatusBadRequest, gin.H{
        "error": "Failed to create User",
    })
    return
}

tx.Commit() 
c.JSON(http.StatusOK, gin.H{
    "message": "Student created successfully",
})
}




func GetStc(c *gin.Context){
    fmt.Println("controller started")

    var stc []Stc
    if result := config.DB.Find(&stc); result.Error != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"stc":stc})
}


func GetStcByID(c *gin.Context){
    id := c.Param("id")
    if id == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error":"Stc ID is required"})
        return
    }

    var stc Stc
    if err := config.DB.First(&stc, id).Error; err != nil{
        if errors.Is(err, gorm.ErrRecordNotFound) {
            c.JSON(http.StatusNotFound, gin.H{"error":"Mda not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error":"Internal Server Error"})
        }
        return
    }

    c.JSON(http.StatusOK, stc)
}


func TotalNumberOfStc(c *gin.Context){
     fmt.Println("Get Total Number Of Mda")
     var count int64

   if result := config.DB.Model(&Stc{}).Count(&count); result.Error != nil{
     c.JSON(http.StatusInternalServerError, gin.H{"error":result.Error.Error()})
     return
   }

   c.JSON(http.StatusOK, gin.H{"total_count":count})
}


func TotalNumberOfOperationalStc(c *gin.Context){
    fmt.Println("Get Total Number Of Active Mda")
    var count int64

    if result := config.DB.Model(&Stc{}).Where("is_operational  = ?", true).Count(&count); result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"total_is_active_count":count})
}


func TotalNumberOfInOperationalStc(c *gin.Context){
    fmt.Println("Get Total Number Of InActive Mda")
    var count int64

    if result := config.DB.Model(&Stc{}).Where("is_operational = ?", false).Count(&count); result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"total_in_active_count":count})
}





func SearchStc(c *gin.Context) {
    query := c.Query("query")
    if query == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
        return
    }

    var stcsearch []Stc
    if err := config.DB.Where("ownership LIKE ? OR center_code LIKE ? OR name LIKE ? OR local_government ? OR state LIKE ? OR LIKE certificate_of_operational_url", "%"+query+"%", "%"+query+"%","%"+query+"%","%"+query+"%","%"+query+"%","%"+query+"%").Find(&stcsearch).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
        return
    }

    if len(stcsearch) == 0 {
        c.JSON(http.StatusOK, gin.H{"message": "No matching stc found"})
        return
    }

    c.JSON(http.StatusOK, stcsearch)
}


func UpdateStc(c *gin.Context) {
    id := c.Param("id")
    if id == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Stc ID is required"})
        return
    }

    var stc Stc
    if err := config.DB.First(&stc, id).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            c.JSON(http.StatusNotFound, gin.H{"error": "Mda not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
        }
        return
    }

    if err := c.BindJSON(&stc); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
        return
    }

    if err := config.DB.Save(&stc).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update Mda"})
        return
    }

    c.JSON(http.StatusOK, stc)
}



func FilterStcDescending(c *gin.Context) {
    var stcs []Stc
    if result := config.DB.Find(&stcs); result.Error != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
        return
    }

    // Sort MDAs by name in ascending order
    sort.Slice(stcs, func(i, j int) bool {
        return stcs[i].Ownership < stcs[j].Ownership
    })

    // Send the sorted MDAs as the response
    c.JSON(http.StatusOK, gin.H{"stcs": stcs})
}


func FilterStcAscending(c *gin.Context) {
    var stcs []Stc
    if result := config.DB.Find(&stcs); result.Error != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
        return
    }

    // Sort MDAs by name in ascending order
    sort.Slice(stcs, func(i, j int) bool {
        return stcs[i].Ownership > stcs[j].Ownership
    })

    // Send the sorted MDAs as the response
    c.JSON(http.StatusOK, gin.H{"stcs": stcs})
}




func SuspendStc(c *gin.Context) {
    id := c.Param("id")
    if id == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Stc ID is required"})
        return
    }

    var stc Stc
    if err := config.DB.First(&stc, id).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            c.JSON(http.StatusNotFound, gin.H{"error": "Stc not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
        }
        return
    }

    stc.isOperational = false 

    if err := config.DB.Save(&stc).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to suspend Stc"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"suspend":"Stc  Suspended"})
}



func ActivateStc(c *gin.Context) {
    id := c.Param("id")
    if id == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Stc ID is required"})
        return
    }

    var stc Stc
    if err := config.DB.First(&stc, id).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            c.JSON(http.StatusNotFound, gin.H{"error": "STC not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
        }
        return
    }

    stc.isOperational = true 

    if err := config.DB.Save(&stc).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to activate STC"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"Activate":"STC Activated"})
}