package stc

import (
	"errors"
	"fme_backend/internal/config"
    myuser "fme_backend/internal/user"
    "fme_backend/internal/utilities"
	"fmt"
	"net/http"
	"sort"
    "strconv"
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
  result, message, newUserID := myuser.CreateStcUser(tx, StcCreateSchema.Email, "dfcv")
  if !result {
      tx.Rollback() // Rollback if user creation fails
      c.JSON(http.StatusBadRequest, gin.H{
          "error": message,
      })
      return
  }
  	stc := Stc{
		Name:              StcCreateSchema.Name,
		Address:           StcCreateSchema.Address,
		State: 			   state,
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
        "message": "Stc created successfully",
    })
}


func CreateMdaStc(c *gin.Context){
    
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
  result, message, newUserID := myuser.CreateStcUser(tx,  StcCreateSchema.Email, "dfcv")
  if !result {
      tx.Rollback() // Rollback if user creation fails
      c.JSON(http.StatusBadRequest, gin.H{
          "error": message,
      })
      return
  }

  stc := Stc{
    Name:              StcCreateSchema.Name,
    Address:           StcCreateSchema.Address,
    State: 			   state,
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
    "message": "Stc created successfully",
})
}




func GetStc(c *gin.Context){
     fmt.Println("controller started")
    //  userRole := c.GetString("user_role")
    //  if userRole != "fme" && userRole != "mda" {
    //      c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
    //      return
    //  }
     limitStr := c.Query("limit")
     pageStr := c.Query("page")



     limit, err := strconv.Atoi(limitStr)
    if err != nil || limit <= 0{
        limit = 100
    }
     page, err := strconv.Atoi(pageStr)
     if err != nil || page <= 0 {
         page = 1
     }
 
     offset := (page - 1) * limit

    var stcs []struct{
        Stc
        Email    string  `json:"email"`
        IsActive bool    `json:"is_active"`
    }
    if result := config.DB.Table("stcs").
        Select("stcs.*, users.email, users.is_active").
        Joins("JOIN users ON stcs.user_id = users.id").
        Limit(limit).
        Offset(offset).
        Find(&stcs); result.Error != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"stc":stcs})
}


func GetStcByID(c *gin.Context){
    id := c.Param("id")
    if id == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error":"Stc ID is required"})
        return
    }

    var stc struct {
        Stc
        Email    string  `json:"email"`
        IsActive bool    `json:"is_active"`
    }
    result := config.DB.Table("stcs").
    Select("mdas.*, users.email, users.is_active").
    Joins("JOIN users ON stcs.user_id = users.id").
    First(&stc, id)
if result.Error != nil {
    c.JSON(http.StatusNotFound, gin.H{
        "error": "STC not found",
    })
    return
}
  

    c.JSON(http.StatusOK, stc)
}


// func SearchStc(c *gin.Context) {
//     query := c.Query("query")
//     if query == "" {
//         c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
//         return
//     }

//     var stcsearch []Stc
//     if err := config.DB.Where("ownership LIKE ? OR centre_code LIKE ? OR name LIKE ? OR local_government LIKE ? OR state LIKE ? OR  certificate_of_operational_url LIKE ?", "%"+query+"%", "%"+query+"%","%"+query+"%","%"+query+"%","%"+query+"%","%"+query+"%").Find(&stcsearch).Error; err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
//         return
//     }

//     if len(stcsearch) == 0 {
//         c.JSON(http.StatusOK, gin.H{"message": "No matching stc found"})
//         return
//     }

//     c.JSON(http.StatusOK, stcsearch)
// }


func SearchStc(c *gin.Context) {
    query := c.Query("query")
    if query == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
        return
    }

    var stcsearch []struct {
        Stc
        Email    string `json:"email"`
        IsActive bool   `json:"is_active"`
    }
    if err := config.DB.Table("mdas").
        Select("stcs.*, users.email, users.is_active").
        Joins("JOIN users ON stcs.user_id = users.id").
        Where("name LIKE ? OR address LIKE ? OR state LIKE ?", "%"+query+"%", "%"+query+"%", "%"+query+"%").
        Find(&stcsearch).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search for Stcs", "details": err.Error()})
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
    if err := config.DB.Table("stcs").
    Select("stcs.*, users.email, users.is_active").
    Joins("JOIN users ON stcs.user_id = users.id").
    First(&stc, id).Error; err != nil {
    if errors.Is(err, gorm.ErrRecordNotFound) {
        c.JSON(http.StatusNotFound, gin.H{"error": "Stc not found"})
    } else {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch Stc"})
    }
    return
}

    if err := c.BindJSON(&stc); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
        return
    }

    if err := config.DB.Save(&stc).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update Stc"})
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

    
    sort.Slice(stcs, func(i, j int) bool {
        return stcs[i].Name < stcs[j].Name
    })

    
    c.JSON(http.StatusOK, gin.H{"stcs": stcs})
}


func FilterStcAscending(c *gin.Context) {
    var stcs []Stc
    if result := config.DB.Find(&stcs); result.Error != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
        return
    }

   
    sort.Slice(stcs, func(i, j int) bool {
        return stcs[i].Name > stcs[j].Name
    })

    
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

    // Retrieve the associated user using the UserID field in the Stc model
    var user myuser.User
    if err := config.DB.First(&user, stc.UserID).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch associated user"})
        return
    }

    // Update the IsActive field of the associated user to false
    user.IsActive = false
    if err := config.DB.Save(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to suspend Stc"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"suspend": "Stc Suspended"})
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
            c.JSON(http.StatusNotFound, gin.H{"error": "Stc not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
        }
        return
    }

    // Retrieve the associated user using the UserID field in the Stc model
    var user myuser.User
    if err := config.DB.First(&user, stc.UserID).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch associated user"})
        return
    }

    // Update the IsActive field of the associated user to true
    user.IsActive = true
    if err := config.DB.Save(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to activate Stc"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"activate": "Stc Activated"})
}


func FilterStcByState(c *gin.Context){
    state := c.Query("state")
    if state == ""{
        c.JSON(http.StatusBadRequest, gin.H{
            "error":"State parameter is required",
        })
        return
    }

    var stcs []Stc
    tx := config.DB.Where("state = ?", state).Find(&stcs)
    if tx.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error":"Failed to retrieve STCs",
        })
         return
    }
    c.JSON(http.StatusOK, gin.H{"stcs":stcs})
}



func StcTotal(c *gin.Context){
    fmt.Println("Get Total number of stc, active stc, inactive stc")
  
    var totalCount, activeCount, inactiveCount int64
    if result := config.DB.Model(&Stc{}).Count(&totalCount); result.Error != nil{
        c.JSON(http.StatusInternalServerError, gin.H{"error":result.Error.Error()})
        return
    }

    if result := config.DB.Model(&Stc{}).Where("is_operational  = ?", true).Count(&activeCount); result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        return
    }

    if result := config.DB.Model(&Stc{}).Where("is_operational = ?", false).Count(&inactiveCount); result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"total_count":totalCount, "total_is_active":activeCount, "total_in_active":inactiveCount})
}
