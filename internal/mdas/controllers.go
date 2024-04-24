package mda

import (
	"errors"
	"fme_backend/internal/config"
	myuser "fme_backend/internal/user"
	"fme_backend/internal/utilities"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"time"

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
   result, message, newUserID := myuser.CreateMdaUser(tx,MdaCreateSchema.Email, "dfcv")
   if !result {
       tx.Rollback() // Rollback if user creation fails
       c.JSON(http.StatusBadRequest, gin.H{
           "error": message,
       })
       return
   }

    mda := Mda{
        RegisterName: MdaCreateSchema.Name,
        Address:MdaCreateSchema.Address,
        StateOfOperation:stateOfOperation,
        // PhoneNumber: MdaCreateSchema.PhoneNumber,
        UserID: newUserID,
    }

	fmt.Println(mda)
    mdaresult := tx.Create(&mda) 
    if mdaresult.Error != nil {
        tx.Rollback() 
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Failed to create Mda",
        })
        return
    }
    tx.Commit() 
    c.JSON(http.StatusOK, gin.H{
        "message": "Mda created successfully",
    })
}
func GetAllMdas(c *gin.Context) {
    fmt.Println("controller started")
    limitStr := c.Query("limit")
    pageStr := c.Query("page")

    limit, err := strconv.Atoi(limitStr)
    if err != nil || limit <= 0 {
        limit = 100
    }

    page, err := strconv.Atoi(pageStr)
    if err != nil || page <= 0 {
        page = 1
    }

    offset := (page - 1) * limit

    var mdas []struct {
        Id          uint
        StateOfOperation    string
        Name         string
        Address      string
        IsActive     bool   `json:"is_active"`
        STCCount     int    `json:"stc_count"`
        StudentCount int    `json:"student_count"`
        UserId  uint
    }

    if result := config.DB.Table("mdas").
    Select("mdas.id AS id, mdas.register_name AS name, mdas.address AS address,mdas.state_of_operation AS state_of_operation, users.is_active AS is_active, users.id AS user_id, COUNT(DISTINCT stcs.id) AS stc_count, COUNT(DISTINCT students.id) AS student_count").
    Joins("JOIN users ON mdas.user_id = users.id").
    Joins("LEFT JOIN stcs ON mdas.id = stcs.mda_id").
    Joins("LEFT JOIN students ON mdas.id = students.mda_id").
    Group("mdas.id, mdas.register_name, mdas.address, users.is_active, users.id, mdas.state_of_operation"). // Include all non-aggregated columns in GROUP BY
    Limit(limit).
    Offset(offset).
    Find(&mdas); result.Error != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"mdas": mdas})
}



func GetMdaByID(c *gin.Context) {
    idStr := c.Param("id")
    if idStr == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Mda ID is required"})
        return
    }

    id, err := strconv.Atoi(idStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "message": "Invalid Mda ID provided",
        })
        return
    }

    var mda struct {
        Id      uint
        Name       string
        CreatedAt   time.Time
        CourseCount     uint
        Address         string
        Email     string `json:"email"`
        IsActive  bool   `json:"is_active"`
        STCCount  int    `json:"stc_count"`
        StudentCount int  `json:"student_count"`
        UserId      uint
    }

    result := config.DB.Table("mdas").
        Select("mdas.id AS id, mdas.register_name AS name, mdas.address AS address, MAX(mdas.created_at) AS created_at, users.is_active  AS is_active, users.email AS email, users.id AS user_id, COUNT(DISTINCT stcs.id) AS stc_count, COUNT(DISTINCT students.id) AS student_count, COUNT(DISTINCT mda_courses.course_id) AS course_count").
        Joins("JOIN users ON mdas.user_id = users.id").
        Joins("LEFT JOIN stcs ON mdas.id = stcs.mda_id").
        Joins("LEFT JOIN students ON mdas.id = students.mda_id").
        Joins("LEFT JOIN mda_courses ON mdas.id = mda_courses.mda_id").
        Group("mdas.id, mdas.register_name, mdas.address, users.is_active, users.id,users.email").
        First(&mda, id)

    if result.Error != nil {
        c.JSON(http.StatusNotFound, gin.H{
            "error": "MDA not found",
        })
        return
    }

    c.JSON(http.StatusOK, mda)
}

func UpdateMda(c *gin.Context) {
    id := c.Param("id")
    if id == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Mda ID is required"})
        return
    }

    var mda Mda

    // Fetch the MDA including the associated user information
    if err := config.DB.Table("mdas").
        Select("mdas.*, users.email, users.is_active").
        Joins("JOIN users ON mdas.user_id = users.id").
        Joins("LEFT JOIN stcs ON mdas.id = stcs.mda_id").
        Joins("LEFT JOIN students ON mdas.id = students.mda_id").
        First(&mda, id).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            c.JSON(http.StatusNotFound, gin.H{"error": "Mda not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch Mda"})
        }
        return
    }

    // Update the MDA with the provided data
    if err := c.BindJSON(&mda); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
        return
    }

    // Save the updated MDA to the database
    if err := config.DB.Save(&mda).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update Mda", "details": err.Error()})
        return
    }

    c.JSON(http.StatusOK, mda)
}



func SearchMda(c *gin.Context) {
    query := c.Query("query")
    if query == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
        return
    }

    var mdasearch []struct {
        Mda
        Email    string `json:"email"`
        IsActive bool   `json:"is_active"`
        STCCount  int    `json:"stc_count"`
        StudentCount int  `json:"student_count"`
    }

    // Update the SQL query to include the email and isActive fields from the users table
    if err := config.DB.Table("mdas").
        // Select("mdas.*, users.email, users.is_active").
	 Select("mdas.*, MAX(users.email) AS email, MAX(users.is_active) AS is_active, COUNT(stcs.id) AS stc_count, COUNT(students.id) AS student_count").
        Joins("JOIN users ON mdas.user_id = users.id").
        Joins("LEFT JOIN stcs ON mdas.id = stcs.mda_id").
        Joins("LEFT JOIN students ON mdas.id = students.mda_id").
        Where("register_name LIKE ? OR address LIKE ? OR state_of_operation LIKE ?", "%"+query+"%", "%"+query+"%", "%"+query+"%").
        Find(&mdasearch).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search for MDAs", "details": err.Error()})
        return
    }

    if len(mdasearch) == 0 {
        c.JSON(http.StatusOK, gin.H{"message": "No matching mdas found"})
        return
    }

    c.JSON(http.StatusOK, mdasearch)
}




func FilterMdaAscending(c *gin.Context) {
    var mdas []struct {
        Mda
        Email    string `json:"email"`
        IsActive bool   `json:"is_active"`
    }

    // Query to join Mda and User tables
    result := config.DB.Table("mdas").
        Select("mdas.*, users.email, users.is_active").
        Joins("JOIN users ON mdas.user_id = users.id").
        Find(&mdas)

    if result.Error != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
        return
    }

    // Sort MDAs by name in ascending order
    sort.Slice(mdas, func(i, j int) bool {
        return mdas[i].RegisterName < mdas[j].RegisterName
    })

    // Send the sorted MDAs with user details as the response
    c.JSON(http.StatusOK, gin.H{"mdas": mdas})
}

func FilterMdaDescending(c *gin.Context) {
    var mdas []struct {
        Mda
        Email    string `json:"email"`
        IsActive bool   `json:"is_active"`
    }

    // Query to join Mda and User tables
    result := config.DB.Table("mdas").
        Select("mdas.*, users.email, users.is_active").
        Joins("JOIN users ON mdas.user_id = users.id").
        Find(&mdas)

    if result.Error != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
        return
    }

    // Sort MDAs by name in descending order
    sort.Slice(mdas, func(i, j int) bool {
        return mdas[i].RegisterName > mdas[j].RegisterName
    })

    // Send the sorted MDAs with user details as the response
    c.JSON(http.StatusOK, gin.H{"mdas": mdas})
}




func FilterMdaByState(c *gin.Context) {
    StateOfOperation := c.Query("state_of_operation")
    if StateOfOperation == "" {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "State parameter is required",
        })
        return
    }

    // Retrieve MDAs matching the provided StateOfOperation
    var mdas []struct {
        Mda
        Email    string `json:"email"`
        IsActive bool   `json:"is_active"`
    }
    tx := config.DB.Table("mdas").
        Select("mdas.*, users.email, users.is_active").
        Joins("JOIN users ON mdas.user_id = users.id").
        Where("state_of_operation = ?", StateOfOperation).
        Find(&mdas)
    if tx.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error":   "Failed to retrieve MDAs",
            "details": tx.Error.Error(), // Include details of the error
        })
        return
    }

    // Return the retrieved MDAs as the response
    c.JSON(http.StatusOK, gin.H{"mdas": mdas})
}




func SuspendMda(c *gin.Context) {
    id := c.Param("id")
    if id == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Mda ID is required"})
        return
    }

    var mda Mda
    if err := config.DB.First(&mda, id).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            c.JSON(http.StatusNotFound, gin.H{"error": "mda not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
        }
        return
    }

    // Retrieve the associated user using the UserID field in the Stc model
    var user myuser.User
    if err := config.DB.First(&user, mda.UserID).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch associated user"})
        return
    }

    // Update the IsActive field of the associated user to false
    user.IsActive = false
    if err := config.DB.Save(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to suspend Stc"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"suspend": "Mda Suspended"})
}

func ActivateMda(c *gin.Context) {
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

    // Retrieve the associated user using the UserID field in the Stc model
    var user myuser.User
    if err := config.DB.First(&user, mda.UserID).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch associated user"})
        return
    }

    // Update the IsActive field of the associated user to true
    user.IsActive = true
    if err := config.DB.Save(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to activate Mda"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"activate": "Mda Activated"})
}


func MdaTotal(c *gin.Context) {
    fmt.Println("Get Total number of MDA, active MDA, inactive MDA")
    var totalCount, activeCount, inactiveCount int64

    // Total number of MDAs
    if result := config.DB.Model(&Mda{}).Count(&totalCount); result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        return
    }

    // Total number of active MDAs
    if result := config.DB.Model(&Mda{}).
        Joins("JOIN users ON mdas.user_id = users.id").
        Where("users.is_active = ?", true).
        Count(&activeCount); result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        return
    }

    // Total number of inactive MDAs
    inactiveCount = totalCount - activeCount

    c.JSON(http.StatusOK, gin.H{
        "total_mda":            totalCount,
        "total_active_mda":     activeCount,
        "total_inactive_mda":   inactiveCount,
    })
}
