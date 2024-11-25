package student

import (
	"encoding/csv"
	"errors"
	"fme_backend/internal/artisans"
	"fme_backend/internal/config"
	myuser "fme_backend/internal/user"
	"fme_backend/internal/utilities"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)




func CreateFmeStudent(c *gin.Context) {
    if c.Bind(&CreateStudentSchema) != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "message": "Failed to read request body",
        })
        return
    }

    // Data Validation
    var stateOfOrigin string
    var result bool
    if (CreateStudentSchema.StateOfOrigin != "") {
    stateOfOrigin, result = utilities.ValidateState(CreateStudentSchema.StateOfOrigin)
    if !result {
        c.JSON(http.StatusBadRequest, gin.H{
            "message": "Incorrect state of origin",
        })
        return
    }
    }
    

    StateOfResidence, result := utilities.ValidateState(CreateStudentSchema.StateOfResidence)
    if !result {
        c.JSON(http.StatusBadRequest, gin.H{
            "message": "Incorrect state of residence",
        })
        return
    }

    if !utilities.IsNigerianPhoneNumber(CreateStudentSchema.PhoneNumber){
        c.JSON(http.StatusBadRequest, gin.H{
            "message": "Incorrect phone number",
        })
        return
    }


    dOB, err := utilities.ParseDoB(CreateStudentSchema.DOBstring)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "message": "Incorrect date format",
        })
        return
    }

    gender, result := utilities.ValidateGender(CreateStudentSchema.Gender)
    if !result {
        c.JSON(http.StatusBadRequest, gin.H{
            "message": "Incorrect Gender",
        })
        return
    }

    result = utilities.VeriryNINFormat(CreateStudentSchema.NationalIdentityNumber)
    if !result {
        c.JSON(http.StatusBadRequest, gin.H{
            "message": "Incorrect NIN format",
        })
        return
    }

    if CreateStudentSchema.IsDisabled {
        if !validateDisability(CreateStudentSchema.DisabilityName) {
            c.JSON(http.StatusBadRequest, gin.H{
                "message": "Wrong disability name",
            })
            return
        }
    }

    // Transaction Handling
    tx := config.DB.Begin() // Begin a transaction

    // Create user with transaction
    result, message, newUserID := myuser.CreateStudentUser(tx, CreateStudentSchema.Email, "dfcv")
    if !result {
        tx.Rollback() // Rollback if user creation fails
        c.JSON(http.StatusBadRequest, gin.H{
            "message": message,
        })
        return
    }

    


    student := Student{
        Firstname: CreateStudentSchema.Firstname,
        Lastname: CreateStudentSchema.Lastname,
        Gender: gender,
        StateOfOrigin: stateOfOrigin,
        StateOfResidence: StateOfResidence,
        DOB: dOB,
        UserID: newUserID,
        Fmestudent: true,
        SID: CreateStudentSchema.SID,
        Address: CreateStudentSchema.Address,
        PhoneNumber: CreateStudentSchema.PhoneNumber,
        NationalIdentityNumber: CreateStudentSchema.NationalIdentityNumber,
        LocalGovernment: CreateStudentSchema.LocalGovernment,
        IsDisabled: CreateStudentSchema.IsDisabled,
        DisabilityName: CreateStudentSchema.DisabilityName,
    }
    studentresult := tx.Create(&student) // Create student within transaction
    if studentresult.Error != nil {
        tx.Rollback() // Rollback if student creation fails
        c.JSON(http.StatusBadRequest, gin.H{
            "message": "Failed to create User",
        })
        return
    }

    // Add Course
    studentcourse := StudentCourse{
        StudentID: student.ID,
        CourseID: CreateStudentSchema.CourseID,
    }
    studentCourseResult := tx.Create(&studentcourse) // Create student within transaction
    if studentCourseResult.Error != nil {
        tx.Rollback() // Rollback if student creation fails
        c.JSON(http.StatusBadRequest, gin.H{
            "message": "Failed to create User",
        })
        return
    }

    tx.Commit() // Commit the transaction if both creations are successful
    // Success response
    c.JSON(http.StatusOK, gin.H{
        "message": "Student created successfully",
    })
}

func CreateMdaStudent(c *gin.Context) {
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

	
	// bind the post data
	if c.Bind(&CreateStudentSchema) != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "message": "Failed to read request body",
        })
        return
    }

    // Data Validation
    var stateOfOrigin string
    var result bool
    if (CreateStudentSchema.StateOfOrigin != "") {
    stateOfOrigin, result = utilities.ValidateState(CreateStudentSchema.StateOfOrigin)
    if !result {
        c.JSON(http.StatusBadRequest, gin.H{
            "message": "Incorrect state of origin",
        })
        return
    }
    }

    if !utilities.IsNigerianPhoneNumber(CreateStudentSchema.PhoneNumber){
        c.JSON(http.StatusBadRequest, gin.H{
            "message": "Incorrect phone number",
        })
        return
    }


    StateOfResidence, result := utilities.ValidateState(CreateStudentSchema.StateOfResidence)
    if !result {
        c.JSON(http.StatusBadRequest, gin.H{
            "message": "Incorrect state of residence",
        })
        return
    }

    dOB, err := utilities.ParseDoB(CreateStudentSchema.DOBstring)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "message": "Incorrect date format",
        })
        return
    }

    gender, result := utilities.ValidateGender(CreateStudentSchema.Gender)
    if !result {
        c.JSON(http.StatusBadRequest, gin.H{
            "message": "Incorrect Gender",
        })
        return
    }

    result = utilities.VeriryNINFormat(CreateStudentSchema.NationalIdentityNumber)
    if !result {
        c.JSON(http.StatusBadRequest, gin.H{
            "message": "Incorrect NIN format",
        })
        return
    }

    if CreateStudentSchema.IsDisabled {
        if !validateDisability(CreateStudentSchema.DisabilityName) {
            c.JSON(http.StatusBadRequest, gin.H{
                "message": "Wrong disability name",
            })
            return
        }
    }
    

    // Transaction Handling
    tx := config.DB.Begin() // Begin a transaction

    // Create user with transaction
    result, message, newUserID := myuser.CreateStudentUser(tx, CreateStudentSchema.Email, "dfcv")
    if !result {
        tx.Rollback() // Rollback if user creation fails
        c.JSON(http.StatusBadRequest, gin.H{
            "message": message,
        })
        return
    }

    student := Student{
        Firstname: CreateStudentSchema.Firstname,
        Lastname: CreateStudentSchema.Lastname,
        Gender: gender,
        StateOfOrigin: stateOfOrigin,
        StateOfResidence: StateOfResidence,
        DOB: dOB,
        UserID: newUserID,
        MdaID: mdaID,
        PhoneNumber: CreateStudentSchema.PhoneNumber,
        SID: CreateStudentSchema.SID,
        Address: CreateStudentSchema.Address,
        LocalGovernment: CreateStudentSchema.LocalGovernment,
        NationalIdentityNumber: CreateStudentSchema.NationalIdentityNumber,
        IsDisabled: CreateStudentSchema.IsDisabled,
        DisabilityName: CreateStudentSchema.DisabilityName,

    }
    studentresult := tx.Create(&student) // Create student within transaction
    if studentresult.Error != nil {
        tx.Rollback() // Rollback if student creation fails
        c.JSON(http.StatusBadRequest, gin.H{
            "message": "Failed to create User",
        })
        return
    }

    // Add Course
    studentcourse := StudentCourse{
        StudentID: student.ID,
        CourseID: CreateStudentSchema.CourseID,
    }
    studentCourseResult := tx.Create(&studentcourse) // Create student within transaction
    if studentCourseResult.Error != nil {
        tx.Rollback() // Rollback if student creation fails
        c.JSON(http.StatusBadRequest, gin.H{
            "message": "Failed to create User",
        })
        return
    }

    tx.Commit() // Commit the transaction if both creations are successful
    // Success response
    c.JSON(http.StatusOK, gin.H{
        "message": "Student created successfully",
    })
}

func CreateStcStudent(c *gin.Context) {
	// retrieve the stc id
	stcIDStr,exists := c.Get("stcID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Problem with the authorization token",
		})
		return
	}
	stcID,ok := stcIDStr.(uint)
	if !ok{
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Problem with the authorization token",
		})
		return
	}

	

	if c.Bind(&CreateStudentSchema) != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "message": "Failed to read request body",
        })
        return
    }

    // Data Validation
    var stateOfOrigin string
    var result bool
    if (CreateStudentSchema.StateOfOrigin != "") {
    stateOfOrigin, result = utilities.ValidateState(CreateStudentSchema.StateOfOrigin)
    if !result {
        c.JSON(http.StatusBadRequest, gin.H{
            "message": "Incorrect state of origin",
        })
        return
    }
    }

    StateOfResidence, result := utilities.ValidateState(CreateStudentSchema.StateOfResidence)
    if !result {
        c.JSON(http.StatusBadRequest, gin.H{
            "message": "Incorrect state of residence",
        })
        return
    }

    if !utilities.IsNigerianPhoneNumber(CreateStudentSchema.PhoneNumber){
        c.JSON(http.StatusBadRequest, gin.H{
            "message": "Incorrect phone number",
        })
        return
    }

    dOB, err := utilities.ParseDoB(CreateStudentSchema.DOBstring)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "message": "Incorrect date format",
        })
        return
    }

    gender, result := utilities.ValidateGender(CreateStudentSchema.Gender)
    if !result {
        c.JSON(http.StatusBadRequest, gin.H{
            "message": "Incorrect Gender",
        })
        return
    }

    result = utilities.VeriryNINFormat(CreateStudentSchema.NationalIdentityNumber)
    if !result {
        c.JSON(http.StatusBadRequest, gin.H{
            "message": "Incorrect NIN format",
        })
        return
    }

    if CreateStudentSchema.IsDisabled {
        if !validateDisability(CreateStudentSchema.DisabilityName) {
            c.JSON(http.StatusBadRequest, gin.H{
                "message": "Wrong disability name",
            })
            return
        }
    }

    // Transaction Handling
    tx := config.DB.Begin() // Begin a transaction

    // Create user with transaction
    result, message, newUserID := myuser.CreateStudentUser(tx, CreateStudentSchema.Email, "dfcv")
    if !result {
        tx.Rollback() // Rollback if user creation fails
        c.JSON(http.StatusBadRequest, gin.H{
            "message": message,
        })
        return
    }

    student := Student{
        Firstname: CreateStudentSchema.Firstname,
        Lastname: CreateStudentSchema.Lastname,
        Gender: gender,
        StateOfOrigin: stateOfOrigin,
        StateOfResidence: StateOfResidence,
        DOB: dOB,
        UserID: newUserID,
        StcID: stcID,
        SID: CreateStudentSchema.SID,
        Address: CreateStudentSchema.Address,
        PhoneNumber: CreateStudentSchema.PhoneNumber,
        LocalGovernment: CreateStudentSchema.LocalGovernment,
        NationalIdentityNumber: CreateStudentSchema.NationalIdentityNumber,
        IsDisabled: CreateStudentSchema.IsDisabled,
        DisabilityName: CreateStudentSchema.DisabilityName,
    }
    studentresult := tx.Create(&student) // Create student within transaction
    if studentresult.Error != nil {
        tx.Rollback() // Rollback if student creation fails
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Failed to create User",
        })
        return
    }

    // Add Course
    studentcourse := StudentCourse{
        StudentID: student.ID,
        CourseID: CreateStudentSchema.CourseID,
    }
    studentCourseResult := tx.Create(&studentcourse) // Create student within transaction
    if studentCourseResult.Error != nil {
        tx.Rollback() // Rollback if student creation fails
        c.JSON(http.StatusBadRequest, gin.H{
            "message": "Failed to create User",
        })
        return
    }

    tx.Commit() // Commit the transaction if both creations are successful
    // Success response
    c.JSON(http.StatusOK, gin.H{
        "message": "Student created successfully",
    })
}

func GetAllStudents(c *gin.Context) {
    // get page value
    pageStr := c.Query("page")
    // mdaIDStr := c.Query("mda_id")
    // stcIDStr := c.Query("stc_id")

    limit:= 100

    page, err := strconv.Atoi(pageStr)
    if err != nil || page <= 0 {
        page = 1 // Default page
    }

    offset := (page - 1) * limit

    //get active filter
    activestr:= c.Query("active")

    // Get userID
    userIDstr,exists := c.Get("userID")
    if !exists{
        c.JSON(http.StatusUnauthorized,gin.H{"message":"unauthorized user"})
        return
    }

    userID,ok := userIDstr.(uint)
    if !ok {
        c.JSON(http.StatusUnauthorized,gin.H{"message":"unauthorized user failed to convert to uint"})
        return
    }

    // Get User Role
    userRoleStr,exists := c.Get("userRole")
    if !exists{
        c.JSON(http.StatusUnauthorized,gin.H{"message":"unauthorized user role"})
        return
    }

    userRole,ok := userRoleStr.(int)
    if !ok {
        c.JSON(http.StatusUnauthorized,gin.H{"message":"unauthorized user failed to convert to uint role"})
        return
    }

    // Make DB queries
    var students []GetAllStudentSchema
    switch (userRole) {
    case 1:
        query := config.DB.Table("students").
        Select("students.id AS id, students.gender AS gender, students.address AS address, MAX(students.created_at) AS created_at, students.firstname AS first_name, students.phone_number AS phone_number, students.state_of_residence AS state_of_residence, students.lastname AS last_name, users.is_active AS is_active, users.id AS user_id, users.email AS email, STRING_AGG(courses.name, ', ') AS courses_taken").
        Joins("JOIN users ON students.user_id = users.id").
        Joins("LEFT JOIN student_courses ON students.id = student_courses.student_id").
        Joins("LEFT JOIN courses ON student_courses.course_id = courses.id").
        Group("students.id, students.phone_number, students.firstname, students.lastname, students.state_of_residence, students.gender, students.address, users.is_active, users.id, users.email").
        Offset(offset).
        Limit(limit)

        // Add active filter:
        if (activestr != "") {
            var isActiveCondition string
            if (activestr == "true") {
                isActiveCondition = "users.is_active = true"
            } else if (activestr == "false"){
                isActiveCondition = "users.is_active = false"
            } else {
                c.JSON(http.StatusBadRequest, gin.H{"message":"incorrect active filter"})
                return
            }
            query = query.Where(isActiveCondition)

        }


        err := query.Scan(&students).Error

        if err!=nil {
            c.JSON(http.StatusBadRequest,gin.H{
                "message":"error retrieving students",
            })
            return
        }
        c.JSON(http.StatusOK,gin.H{"students":students})
        return

    case 2:
        //get mdaid
        var userMdaId uint 
		err := config.DB.Table("mdas").
		Where("user_id = ?", userID).
		Pluck("id",&userMdaId).Error
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "MdaAccount has issues",
			})
			return
		}

        // add filter for the mdaId and related stc
        // use left join to add 
        query := config.DB.Table("students").
            Select("students.id AS id, students.gender AS gender, students.address AS address, MAX(students.created_at) AS created_at, students.firstname AS first_name, students.phone_number AS phone_number, students.state_of_residence AS state_of_residence, students.lastname AS last_name, users.is_active AS is_active, users.id AS user_id, users.email AS email, STRING_AGG(courses.name, ', ') AS courses_taken").
            Joins("JOIN users ON students.user_id = users.id").
            Joins("LEFT JOIN student_courses ON students.id = student_courses.student_id").
            Joins("LEFT JOIN courses ON student_courses.course_id = courses.id").
            Joins("LEFT JOIN stcs ON students.stc_id = stcs.id"). // Joining stc table
            Group("students.id, students.phone_number, students.firstname, students.lastname, students.state_of_residence, students.gender, students.address, users.is_active, users.id, users.email").
            Offset(offset).
            Limit(limit)

            // Add active filter:
            if activestr != "" {
                var isActiveCondition string
                if activestr == "true" {
                    isActiveCondition = "users.is_active = true"
                } else if activestr == "false" {
                    isActiveCondition = "users.is_active = false"
                } else {
                    c.JSON(http.StatusBadRequest, gin.H{"message": "incorrect active filter"})
                    return
                }
                query = query.Where(isActiveCondition)
            }

            // Add condition for mdaid
            
            query = query.Where("students.mda_id = ? OR stcs.mda_id = ?", userMdaId, userMdaId)
           

            nerr := query.Scan(&students).Error
            if nerr!=nil {
                c.JSON(http.StatusBadRequest,gin.H{
                    "message":"error retrieving students",
                })
                return
            }
            c.JSON(http.StatusOK,gin.H{"students":students})
            return


    case 3:
        //get user stcid
			var userStcId uint
			err := config.DB.
					Table("stcs").
					Select("id").
					Where("user_id = ?", userID).
					Scan(&userStcId).Error
			if err != nil{
				c.JSON(http.StatusUnauthorized, gin.H{
					"message": "Error with authorization",
				})
				return
			}

            query := config.DB.Table("students").
                Select("students.id AS id, students.gender AS gender, students.address AS address, MAX(students.created_at) AS created_at, students.firstname AS first_name, students.phone_number AS phone_number, students.state_of_residence AS state_of_residence, students.lastname AS last_name, users.is_active AS is_active, users.id AS user_id, users.email AS email, STRING_AGG(courses.name, ', ') AS courses_taken").
                Joins("JOIN users ON students.user_id = users.id").
                Joins("LEFT JOIN student_courses ON students.id = student_courses.student_id").
                Joins("LEFT JOIN courses ON student_courses.course_id = courses.id").
                Group("students.id, students.phone_number, students.firstname, students.lastname, students.state_of_residence, students.gender, students.address, users.is_active, users.id, users.email").
                Offset(offset).
                Limit(limit)

// Add active filter:
            if activestr != "" {
                var isActiveCondition string
                if activestr == "true" {
                    isActiveCondition = "users.is_active = true"
                } else if activestr == "false" {
                    isActiveCondition = "users.is_active = false"
                } else {
                    c.JSON(http.StatusBadRequest, gin.H{"message": "incorrect active filter"})
                    return
                }
                query = query.Where(isActiveCondition)
            }

            // Add condition for stcid
            query = query.Where("students.stc_id = ?", userStcId)
            

            nerr := query.Scan(&students).Error
            if nerr != nil {
                c.JSON(http.StatusInternalServerError,gin.H{
                    "message":"error retrieving students",
                })
                return
            }
            c.JSON(http.StatusOK,gin.H{"students":students})
            return


    default:
        fmt.Println(userID)
        c.JSON(http.StatusUnauthorized,gin.H{"message":"default unauthorized user"})
        return
    }

}

func GetStudent(c *gin.Context) {
    // GET ID from parameter
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "path parameter not provided",
		})
		return
	}
	// Convert id to string
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "path parameter invalid",
		})
		return
	}

    // Get user ID
    userIDstr,exists := c.Get("userID")
    if !exists{
        c.JSON(http.StatusUnauthorized,gin.H{"message":"unauthorized user"})
        return
    }

    userID,ok := userIDstr.(uint)
    if !ok {
        c.JSON(http.StatusUnauthorized,gin.H{"message":"unauthorized user failed to convert to uint"})
        return
    }

    // Get user Role
    userRoleStr,exists := c.Get("userRole")
    if !exists{
        c.JSON(http.StatusUnauthorized,gin.H{"message":"unauthorized user"})
        return
    }

    userRole,ok := userRoleStr.(int)
    if !ok {
        c.JSON(http.StatusUnauthorized,gin.H{"message":"unauthorized user failed to convert to uint"})
        return
    }


    var student GetStudentSchema
    switch (userRole) {
    case 1:
        err := config.DB.Table("students").
        Select("students.id AS id, students.firstname AS first_name, students.lastname AS last_name, students.phone_number AS phone_number, users.is_active AS is_active,  users.email AS email, users.id AS user_id, STRING_AGG(courses.name, ', ') AS courses_taken, students.gender AS gender, students.state_of_residence AS state_of_residence, students.address AS address, MAX(students.created_at) AS created_at").
        Joins("JOIN users ON students.user_id = users.id").
        Joins("LEFT JOIN student_courses ON students.id = student_courses.student_id").
        Joins("LEFT JOIN courses ON student_courses.course_id = courses.id").
        Where("students.id = ?", id).
        Group("students.id, users.is_active, users.email, users.id, students.firstname, students.lastname, students.phone_number, students.gender, students.state_of_residence, students.address").
        Find(&student,id).Error


        if err!=nil {
            c.JSON(http.StatusBadRequest,gin.H{
                "message":"error retrieving students",
            })
            return
        }
        if (student.ID == 0) {
            c.JSON(http.StatusNotFound,gin.H{
                "message":"record does not exist",
            })
            return
        }
        c.JSON(http.StatusOK,gin.H{"students":student})

    case 2:
        // get mda id
        var userMdaId uint 
			err := config.DB.Table("mdas").
			Where("user_id = ?", userID).
			Pluck("id",&userMdaId).Error
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{
					"message": "MdaAccount has issues",
				})
				return
			}
        // get the instance mda and student id
        var instanceData struct{
            MdaId uint
            StcId uint
        }
        nerr := config.DB.Table("students").
        Select("mda_id, stc_id").
        Where("id = ?", id).
        Scan(&instanceData).Error
        if nerr != nil {
            c.JSON(http.StatusUnauthorized, gin.H{
                "message": "cannot get the instance data",
            })
            return
        }

        if instanceData.MdaId != 0{
            if userMdaId != instanceData.MdaId {
                c.JSON(http.StatusUnauthorized,gin.H{
                    "message":"authorization required to view this student",
                })
                return
            }
            err := config.DB.Table("students").
                    Select("students.id AS id, students.firstname AS first_name, students.lastname AS last_name, students.phone_number AS phone_number, users.is_active AS is_active,  users.email AS email, users.id AS user_id, STRING_AGG(courses.name, ', ') AS courses_taken, students.gender AS gender, students.state_of_residence AS state_of_residence, students.address AS address, MAX(students.created_at) AS created_at").
                    Joins("JOIN users ON students.user_id = users.id").
                    Joins("LEFT JOIN student_courses ON students.id = student_courses.student_id").
                    Joins("LEFT JOIN courses ON student_courses.course_id = courses.id").
                    Where("students.id = ?", id).
                    Group("students.id, users.is_active, users.email, users.id, students.firstname, students.lastname, students.phone_number, students.gender, students.state_of_residence, students.address").
                    Find(&student,id).Error


                    if err!=nil {
                        c.JSON(http.StatusBadRequest,gin.H{
                            "message":"error retrieving students",
                        })
                        return
                    }
                    if (student.ID == 0) {
                        c.JSON(http.StatusNotFound,gin.H{
                            "message":"record does not exist",
                        })
                        return
                    }
                    c.JSON(http.StatusOK,gin.H{"students":student})

        } else if (instanceData.StcId != 0) {
                    var insatnceStcMdaId uint	// get the related mda id 

					err := config.DB.Table("stcs").
					Where("id = ?", instanceData.StcId).
					Select("mda_id").
					Scan(&insatnceStcMdaId).Error
					if err != nil {
						c.JSON(http.StatusBadRequest, gin.H{
							"message": "Unable to update the user record",
						})
						return
					}
                    if userMdaId != insatnceStcMdaId {
                        c.JSON(http.StatusUnauthorized, gin.H{
							"message": "Need authorization to view this student",
						})
						return
                    }
                    nerr := config.DB.Table("students").
                    Select("students.id AS id, students.firstname AS first_name, students.lastname AS last_name, students.phone_number AS phone_number, users.is_active AS is_active,  users.email AS email, users.id AS user_id, STRING_AGG(courses.name, ', ') AS courses_taken, students.gender AS gender, students.state_of_residence AS state_of_residence, students.address AS address, MAX(students.created_at) AS created_at").
                    Joins("JOIN users ON students.user_id = users.id").
                    Joins("LEFT JOIN student_courses ON students.id = student_courses.student_id").
                    Joins("LEFT JOIN courses ON student_courses.course_id = courses.id").
                    Where("students.id = ?", id).
                    Group("students.id, users.is_active, users.email, users.id, students.firstname, students.lastname, students.phone_number, students.gender, students.state_of_residence, students.address").
                    Find(&student,id).Error


                    if nerr!=nil {
                        c.JSON(http.StatusBadRequest,gin.H{
                            "message":"error retrieving students",
                        })
                        return
                    }
                    if (student.ID == 0) {
                        c.JSON(http.StatusNotFound,gin.H{
                            "message":"record does not exist",
                        })
                        return
                    }
                    c.JSON(http.StatusOK,gin.H{"students":student})

        } else {
            c.JSON(http.StatusUnauthorized, gin.H{
                "message": "Authorization needed to view this user",
            })
            return
        }
      
    case 3:
       //get user stcid
		var userStcId uint
		err := config.DB.
				Table("stcs").
				Select("id").
				Where("user_id = ?", userID).
				Scan(&userStcId).Error
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Error trying to fetch data",
			})
			return
		}   
        
        // get student stc id
		var studentStcId uint
		nerr := config.DB.
				Table("students").
				Select("stc_id").
				Where("id = ?", id).
				Scan(&studentStcId).Error

		if nerr != nil{
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Error trying to fetch data",
			})
			return
		}

        if studentStcId != userStcId {
            c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Need permission to view this user",
			})
			return
        }

        verr := config.DB.Table("students").
                    Select("students.id AS id, students.firstname AS first_name, students.lastname AS last_name, students.phone_number AS phone_number, users.is_active AS is_active,  users.email AS email, users.id AS user_id, STRING_AGG(courses.name, ', ') AS courses_taken, students.gender AS gender, students.state_of_residence AS state_of_residence, students.address AS address, MAX(students.created_at) AS created_at").
                    Joins("JOIN users ON students.user_id = users.id").
                    Joins("LEFT JOIN student_courses ON students.id = student_courses.student_id").
                    Joins("LEFT JOIN courses ON student_courses.course_id = courses.id").
                    Where("students.id = ?", id).
                    Group("students.id, users.is_active, users.email, users.id, students.firstname, students.lastname, students.phone_number, students.gender, students.state_of_residence, students.address").
                    Find(&student,id).Error


        if verr!=nil {
            c.JSON(http.StatusBadRequest,gin.H{
                "message":"error retrieving students",
            })
            return
        }
        if (student.ID == 0) {
            c.JSON(http.StatusNotFound,gin.H{
                "message":"record does not exist",
            })
            return
        }
        c.JSON(http.StatusOK,gin.H{"students":student})

    
    default:
        fmt.Println(userID)
        c.JSON(http.StatusUnauthorized,gin.H{"message":"default unauthorized user"})
        return
    }
}

func GetTotalStudentInfo(c *gin.Context) {
    // Get user ID
    userIDstr,exists := c.Get("userID")
    if !exists{
        c.JSON(http.StatusUnauthorized,gin.H{"message":"unauthorized user"})
        return
    }

    userID,ok := userIDstr.(uint)
    if !ok {
        c.JSON(http.StatusUnauthorized,gin.H{"message":"unauthorized user failed to convert to uint"})
        return
    }

    // Get user Role
    userRoleStr,exists := c.Get("userRole")
    if !exists{
        c.JSON(http.StatusUnauthorized,gin.H{"message":"unauthorized user"})
        return
    }

    userRole,ok := userRoleStr.(int)
    if !ok {
        c.JSON(http.StatusUnauthorized,gin.H{"message":"unauthorized user failed to convert to uint"})
        return
    }

    var studentinfo TotalStudentInfo
    switch userRole {
    case 1:
        err:= config.DB.Table("students").
        Joins("JOIN users ON students.user_id = users.id").
        Select("COUNT(DISTINCT students.id) as total_students, SUM(CASE WHEN users.is_active = true THEN 1 ELSE 0 END) as total_active_students, SUM(CASE WHEN users.is_active = false THEN 1 ELSE 0 END) as total_inactive_students").
        Scan(&studentinfo).Error

        if err!=nil {
            c.JSON(http.StatusBadRequest,gin.H{
                "message":"error retrieving students",
            })
            return
        }

        c.JSON(http.StatusOK,gin.H{"studentInfo":studentinfo})

    case 2:
        // Get user MDA ID
			var userMdaId uint 
			err := config.DB.Table("mdas").
			Where("user_id = ?", userID).
			Pluck("id",&userMdaId).Error
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{
					"message": "MdaAccount has issues",
				})
				return
			}

            err = config.DB.Table("students").
                    Joins("JOIN users ON students.user_id = users.id").
                    Joins("LEFT JOIN stcs on students.stc_id = stcs.id").
                    Select("COUNT(DISTINCT students.id) as total_students, SUM(CASE WHEN users.is_active = true THEN 1 ELSE 0 END) as total_active_students, SUM(CASE WHEN users.is_active = false THEN 1 ELSE 0 END) as total_inactive_students").
                    Where("students.mda_id = ? OR stcs.mda_id = ?", userMdaId, userMdaId).
                    Scan(&studentinfo).Error
            if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Error retrieving students",
				})
				return
            }
            c.JSON(http.StatusOK,gin.H{"studentInfo":studentinfo})

        case 3:
            //get user stcid
			var userStcId uint
			err := config.DB.
					Table("stcs").
					Select("id").
					Where("user_id = ?", userID).
					Scan(&userStcId).Error
			if err != nil{
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Error trying to suspend this user",
				})
				return
			}

            err = config.DB.Table("students").
                    Joins("JOIN users ON students.user_id = users.id").
                    Joins("LEFT JOIN stcs on students.stc_id = stcs.id").
                    Select("COUNT(DISTINCT students.id) as total_students, SUM(CASE WHEN users.is_active = true THEN 1 ELSE 0 END) as total_active_students, SUM(CASE WHEN users.is_active = false THEN 1 ELSE 0 END) as total_inactive_students").
                    Where("students.stc_id = ?", userStcId).
                    Scan(&studentinfo).Error
            if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Error retrieving students",
				})
				return
            }
            c.JSON(http.StatusOK,gin.H{"studentInfo":studentinfo})



        default:
        fmt.Println(userID)
        c.JSON(http.StatusUnauthorized,gin.H{"message":"default unauthorized user"})
        return
    }
}


func parseAndValidateRecord(record []string) (CreateStudentSchematype, error) {
    if len(record) < 15 {
        return CreateStudentSchematype{}, fmt.Errorf("incorrect number of fields in the CSV record")
    }

    _, err := utilities.ParseDoB(record[5])
    if err != nil {
        return CreateStudentSchematype{}, fmt.Errorf("invalid date of birth format: %v", err)
    }

    gender, valid := utilities.ValidateGender(record[2])
    if !valid {
        return CreateStudentSchematype{}, fmt.Errorf("invalid gender")
    }

    stateOfOrigin, valid := utilities.ValidateState(record[3])
    if !valid {
        return CreateStudentSchematype{}, fmt.Errorf("invalid state of origin")
    }

    stateOfResidence, valid := utilities.ValidateState(record[4])
    if !valid {
        return CreateStudentSchematype{}, fmt.Errorf("invalid state of residence")
    }

    if !utilities.IsNigerianPhoneNumber("0"+record[7]) {
        return CreateStudentSchematype{}, fmt.Errorf("invalid phone number")
    }

    courseID, err := strconv.ParseUint(record[10], 10, 32)
    if err != nil {
        return CreateStudentSchematype{}, fmt.Errorf("invalid course id: %v", err)
    }

    if !utilities.VeriryNINFormat(record[12]) {
        return CreateStudentSchematype{}, fmt.Errorf("invalid nin format: %v", err)
        
    }
    isDisabled,err := strconv.ParseBool(record[13])
    if err != nil {
        return CreateStudentSchematype{}, fmt.Errorf("invalid is disabled format: %v", err)

    }
    if !validateDisability(record[14]) {
        return CreateStudentSchematype{}, fmt.Errorf("invalid disability name: %v", err)

    }

    return CreateStudentSchematype{
        Firstname:        record[0],
        Lastname:         record[1],
        Gender:           gender,
        StateOfOrigin:    stateOfOrigin,
        StateOfResidence: stateOfResidence,
        DOBstring:        record[5],
        Email:            record[6],
        PhoneNumber:      "0"+record[7],
        SID:              record[8],
        Address:          record[9],
        CourseID:         uint(courseID),
        NationalIdentityNumber: record[11],
        LocalGovernment: record[12],
        IsDisabled: isDisabled,
        DisabilityName: record[14],
    }, nil
}

func createMdaStudentInstance(mdaID uint, schema CreateStudentSchematype,tx *gorm.DB) (Student, error) {

    dob, err := utilities.ParseDoB(schema.DOBstring)
    if err != nil {
        return Student{}, errors.New("wrong date format")
    }

    // Create user with transaction
    result, message, newUserID := myuser.CreateStudentUser(tx, schema.Email, "dfcv")
    if !result {
        return Student{}, fmt.Errorf(message)
    }



    student := Student{
        Firstname:        schema.Firstname,
        Lastname:         schema.Lastname,
        Gender:           schema.Gender,
        StateOfOrigin:    schema.StateOfOrigin,
        StateOfResidence: schema.StateOfResidence,
        DOB:              dob,
        UserID:           newUserID,
        MdaID:            mdaID,
        PhoneNumber:      schema.PhoneNumber,
        SID:              schema.SID,
        Address:          schema.Address,
        IsDisabled: schema.IsDisabled,
        DisabilityName: CreateStudentSchema.DisabilityName,
    }

    if err := tx.Create(&student).Error; err != nil {
        return Student{}, fmt.Errorf("failed to create student: %v", err)
    }

    studentCourse := StudentCourse{
        StudentID: student.ID,
        CourseID:  schema.CourseID,
    }

    if err := tx.Create(&studentCourse).Error; err != nil {
        return Student{}, fmt.Errorf("failed to create student course: %v", err)
    }

    return student, nil
}

func createStcStudentInstance(stcID uint, schema CreateStudentSchematype,tx *gorm.DB) (Student, error) {

    dob, err := utilities.ParseDoB(schema.DOBstring)
    if err != nil {
        return Student{}, errors.New("wrong date format")
    }

    // Create user with transaction
    result, message, newUserID := myuser.CreateStudentUser(tx, schema.Email, "dfcv")
    if !result {
        return Student{}, fmt.Errorf(message)
    }



    student := Student{
        Firstname:        schema.Firstname,
        Lastname:         schema.Lastname,
        Gender:           schema.Gender,
        StateOfOrigin:    schema.StateOfOrigin,
        StateOfResidence: schema.StateOfResidence,
        DOB:              dob,
        UserID:           newUserID,
        StcID:            stcID,
        PhoneNumber:      schema.PhoneNumber,
        SID:              schema.SID,
        Address:          schema.Address,
        DisabilityName: schema.DisabilityName,
        IsDisabled: schema.IsDisabled,
    }

    if err := tx.Create(&student).Error; err != nil {
        return Student{}, fmt.Errorf("failed to create student: %v", err)
    }

    studentCourse := StudentCourse{
        StudentID: student.ID,
        CourseID:  schema.CourseID,
    }

    if err := tx.Create(&studentCourse).Error; err != nil {
        return Student{}, fmt.Errorf("failed to create student course: %v", err)
    }

   
    return student, nil
}


func createFmeStudentInstance(schema CreateStudentSchematype, tx *gorm.DB) (Student, error) {

    dob, err := utilities.ParseDoB(schema.DOBstring)
    if err != nil {
        return Student{}, errors.New("wrong date format")
    }

    // Create user with transaction
    result, message, newUserID := myuser.CreateStudentUser(tx, schema.Email, "dfcv")
    if !result {
        return Student{}, fmt.Errorf(message)
    }



    student := Student{
        Firstname:        schema.Firstname,
        Lastname:         schema.Lastname,
        Gender:           schema.Gender,
        StateOfOrigin:    schema.StateOfOrigin,
        StateOfResidence: schema.StateOfResidence,
        DOB:              dob,
        UserID:           newUserID,
        Fmestudent:       true,
        PhoneNumber:      schema.PhoneNumber,
        SID:              schema.SID,
        Address:          schema.Address,
        DisabilityName: schema.DisabilityName,
        IsDisabled: CreateStudentSchema.IsDisabled,
    }

    if err := tx.Create(&student).Error; err != nil {
        
        return Student{}, fmt.Errorf("failed to create student: %v", err)
    }

    studentCourse := StudentCourse{
        StudentID: student.ID,
        CourseID:  schema.CourseID,
    }

    if err := tx.Create(&studentCourse).Error; err != nil {
        return Student{}, fmt.Errorf("failed to create student course: %v", err)
    }

    return student, nil
}


func CreateMdaStudentFromCsv(c *gin.Context) {
    // Retrieve the MDA ID
    mdaIDStr, exists := c.Get("mdaID")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{
            "message": "Problem with the authorization token",
        })
        return
    }
    mdaID, ok := mdaIDStr.(uint)
    if !ok {
        c.JSON(http.StatusUnauthorized, gin.H{
            "message": "Problem with the authorization token",
        })
        return
    }

    // Get the CSV file from the request
    file, _, err := c.Request.FormFile("file")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "message": "Failed to read request body",
        })
        return
    }
    defer file.Close()

    // Read the CSV file
    reader := csv.NewReader(file)
    var students []Student

    // Skip the header
    if _, err := reader.Read(); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "message": "Failed to read CSV header",
        })
        return
    }
    tx := config.DB.Begin()
    for {
        record, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            tx.Rollback()
            c.JSON(http.StatusBadRequest, gin.H{
                "message": "Failed to read CSV file",
            })
            return
        }

        // Parse and validate each record
        studentSchema, err := parseAndValidateRecord(record)
        if err != nil {
            tx.Rollback()
            c.JSON(http.StatusBadRequest, gin.H{
                "message": fmt.Sprintf("Invalid data: %v", err),
            })
            return
        }

        // Create user and student instances
        student, err := createMdaStudentInstance(mdaID, studentSchema,tx)
        if err != nil {
            tx.Rollback()
            c.JSON(http.StatusBadRequest, gin.H{
                "message": fmt.Sprintf("Failed to create student: %v", err),
            })
            return
        }

        students = append(students, student)
    }
    tx.Commit()

    // If all students are successfully created
    c.JSON(http.StatusOK, gin.H{
        "message": "Students created successfully",
        "students": students,
    })
}



func CreateFmeStudentFromCsv(c *gin.Context) {
    // Get the CSV file from the request
    file, _, err := c.Request.FormFile("file")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "message": "Failed to read request body",
        })
        return
    }
    defer file.Close()

    // Read the CSV file
    reader := csv.NewReader(file)
    var students []Student

    // Skip the header
    if _, err := reader.Read(); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "message": "Failed to read CSV header",
        })
        return
    }

    tx := config.DB.Begin()

    for {
        record, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            tx.Rollback()
            c.JSON(http.StatusBadRequest, gin.H{
                "message": "Failed to read CSV file",
            })
            return
        }

        // Parse and validate each record
        studentSchema, err := parseAndValidateRecord(record)
        if err != nil {
            tx.Rollback()
            c.JSON(http.StatusBadRequest, gin.H{
                "message": fmt.Sprintf("Invalid data: %v", err),
            })
            return
        }

        // Create user and student instances
        student, err := createFmeStudentInstance(studentSchema,tx)
        if err != nil {
            tx.Rollback()
            c.JSON(http.StatusBadRequest, gin.H{
                "message": fmt.Sprintf("Failed to create student: %v", err),
            })
            return
        }

        students = append(students, student)
    }
    tx.Commit()

    // If all students are successfully created
    c.JSON(http.StatusOK, gin.H{
        "message": "Students created successfully",
        "students": students,
    })
}



func CreateStcStudentFromCsv(c *gin.Context) {
    // Retrieve the stc ID
    stcIDStr, exists := c.Get("stcID")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{
            "message": "Problem with the authorization token",
        })
        return
    }
    stcID, ok := stcIDStr.(uint)
    if !ok {
        c.JSON(http.StatusUnauthorized, gin.H{
            "message": "Problem with the authorization token",
        })
        return
    }

    // Get the CSV file from the request
    file, _, err := c.Request.FormFile("file")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "message": "Failed to read request body",
        })
        return
    }
    defer file.Close()

    // Read the CSV file
    reader := csv.NewReader(file)
    var students []Student

    // Skip the header  
    if _, err := reader.Read(); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "message": "Failed to read CSV header",
        })
        return
    }

    // validate the header

    tx := config.DB.Begin()
    // read other records
    for {
        record, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            tx.Rollback()
            c.JSON(http.StatusBadRequest, gin.H{
                "message": "Failed to read CSV file",
            })
            return
        }

        // Parse and validate each record
        studentSchema, err := parseAndValidateRecord(record)
        if err != nil {
            tx.Rollback()
            c.JSON(http.StatusBadRequest, gin.H{
                "message": fmt.Sprintf("Invalid data: %v", err),
            })
            return
        }

        // Create user and student instances
        student, err := createStcStudentInstance(stcID, studentSchema,tx)
        if err != nil {
            tx.Rollback()
            c.JSON(http.StatusBadRequest, gin.H{
                "message": fmt.Sprintf("Failed to create student: %v", err),
            })
            return
        }

        students = append(students, student)
    }
    tx.Commit()

    // If all students are successfully created
    c.JSON(http.StatusOK, gin.H{
        "message": "Students created successfully",
        "students": students,
    })
}



func DownloadStudentsCsv(c *gin.Context) {
    // get active filter
    activestr := c.Query("active")

    // Get userID
    userIDstr, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized user"})
        return
    }

    userID, ok := userIDstr.(uint)
    if !ok {
        c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized user failed to convert to uint"})
        return
    }

    // Get User Role
    userRoleStr, exists := c.Get("userRole")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized user role"})
        return
    }

    userRole, ok := userRoleStr.(int)
    if !ok {
        c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized user failed to convert to uint role"})
        return
    }

    // Make DB queries
    var students []GetAllStudentSchema
    switch userRole {
    case 1:
        query := config.DB.Table("students").
            Select("students.id AS id, students.gender AS gender, students.address AS address, MAX(students.created_at) AS created_at, students.firstname AS first_name, students.phone_number AS phone_number, students.state_of_residence AS state_of_residence, students.lastname AS last_name, users.is_active AS is_active, users.id AS user_id, users.email AS email, STRING_AGG(courses.name, ', ') AS courses_taken").
            Joins("JOIN users ON students.user_id = users.id").
            Joins("LEFT JOIN student_courses ON students.id = student_courses.student_id").
            Joins("LEFT JOIN courses ON student_courses.course_id = courses.id").
            Group("students.id, students.phone_number, students.firstname, students.lastname, students.state_of_residence, students.gender, students.address, users.is_active, users.id, users.email")
            

        // Add active filter:
        if activestr != "" {
            var isActiveCondition string
            if activestr == "true" {
                isActiveCondition = "users.is_active = true"
            } else if activestr == "false" {
                isActiveCondition = "users.is_active = false"
            } else {
                c.JSON(http.StatusBadRequest, gin.H{"message": "incorrect active filter"})
                return
            }
            query = query.Where(isActiveCondition)
        }

        err := query.Scan(&students).Error
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{
                "message": "error retrieving students",
            })
            return
        }

    case 2:
        // get mdaid
        var userMdaId uint
        err := config.DB.Table("mdas").
            Where("user_id = ?", userID).
            Pluck("id", &userMdaId).Error
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{
                "message": "MdaAccount has issues",
            })
            return
        }

        // add filter for the mdaId and related stc
        query := config.DB.Table("students").
            Select("students.id AS id, students.gender AS gender, students.address AS address, MAX(students.created_at) AS created_at, students.firstname AS first_name, students.phone_number AS phone_number, students.state_of_residence AS state_of_residence, students.lastname AS last_name, users.is_active AS is_active, users.id AS user_id, users.email AS email, STRING_AGG(courses.name, ', ') AS courses_taken").
            Joins("JOIN users ON students.user_id = users.id").
            Joins("LEFT JOIN student_courses ON students.id = student_courses.student_id").
            Joins("LEFT JOIN courses ON student_courses.course_id = courses.id").
            Joins("LEFT JOIN stcs ON students.stc_id = stcs.id"). // Joining stc table
            Group("students.id, students.phone_number, students.firstname, students.lastname, students.state_of_residence, students.gender, students.address, users.is_active, users.id, users.email")
            

        // Add active filter:
        if activestr != "" {
            var isActiveCondition string
            if activestr == "true" {
                isActiveCondition = "users.is_active = true"
            } else if activestr == "false" {
                isActiveCondition = "users.is_active = false"
            } else {
                c.JSON(http.StatusBadRequest, gin.H{"message": "incorrect active filter"})
                return
            }
            query = query.Where(isActiveCondition)
        }

        // Add condition for mdaid
        query = query.Where("students.mda_id = ? OR stcs.mda_id = ?", userMdaId, userMdaId)

        err = query.Scan(&students).Error
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{
                "message": "error retrieving students",
            })
            return
        }

    case 3:
        // get user stcid
        var userStcId uint
        err := config.DB.
            Table("stcs").
            Select("id").
            Where("user_id = ?", userID).
            Scan(&userStcId).Error
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{
                "message": "Error with authorization",
            })
            return
        }

        query := config.DB.Table("students").
            Select("students.id AS id, students.gender AS gender, students.address AS address, MAX(students.created_at) AS created_at, students.firstname AS first_name, students.phone_number AS phone_number, students.state_of_residence AS state_of_residence, students.lastname AS last_name, users.is_active AS is_active, users.id AS user_id, users.email AS email, STRING_AGG(courses.name, ', ') AS courses_taken").
            Joins("JOIN users ON students.user_id = users.id").
            Joins("LEFT JOIN student_courses ON students.id = student_courses.student_id").
            Joins("LEFT JOIN courses ON student_courses.course_id = courses.id").
            Group("students.id, students.phone_number, students.firstname, students.lastname, students.state_of_residence, students.gender, students.address, users.is_active, users.id, users.email")
            

        // Add active filter:
        if activestr != "" {
            var isActiveCondition string
            if activestr == "true" {
                isActiveCondition = "users.is_active = true"
            } else if activestr == "false" {
                isActiveCondition = "users.is_active = false"
            } else {
                c.JSON(http.StatusBadRequest, gin.H{"message": "incorrect active filter"})
                return
            }
            query = query.Where(isActiveCondition)
        }

        // Add condition for stcid
        query = query.Where("students.stc_id = ?", userStcId)

        nerr := query.Scan(&students).Error
        if nerr != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "message": "error retrieving students",
            })
            return
        }

    default:
        fmt.Println(userID)
        c.JSON(http.StatusUnauthorized, gin.H{"message": "default unauthorized user"})
        return
    }

    // Generate CSV
    csvData := [][]string{
        {"ID", "Gender", "Address", "Created At", "First Name", "Phone Number", "State of Residence", "Last Name", "Is Active", "User ID", "Email", "Courses Taken"},
    }

    for _, student := range students {
        csvData = append(csvData, []string{
            strconv.Itoa(int(student.ID)),
            student.Gender,
            student.Address,
            student.CreatedAt.String(),
            student.FirstName,
            student.PhoneNumber,
            student.StateOfResidence,
            student.LastName,
            strconv.FormatBool(student.IsActive),
            strconv.Itoa(int(student.UserID)),
            student.Email,
            student.CoursesTaken,
        })
    }

    c.Header("Content-Disposition", "attachment; filename=students.csv")
    c.Header("Content-Type", "text/csv")

    w := csv.NewWriter(c.Writer)
    defer w.Flush()

    for _, record := range csvData {
        if err := w.Write(record); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"message": "error writing csv"})
            return
        }
    }
}

func GraduateStudent(c *gin.Context) {
    idStr:= c.Param("id")
		if idStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "path parameter not provided",
			})
			return
		}

    id, err := strconv.Atoi(idStr)

    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "message": "path parameter invalid",
        })
        return
    }

    if c.Bind(&GraduateStudentSchema) != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "message": "Failed to read request body",
        })
        return
    }

    // check the date
    gradDate, err := utilities.ParseDoB(GraduateStudentSchema.DateOfGrad)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "message": "Incorrect date format",
        })
        return
    }

    // get the user id and then figure out the permission issues
    // Get user ID
    userIDstr,exists := c.Get("userID")
    if !exists{
        c.JSON(http.StatusUnauthorized,gin.H{"message":"unauthorized user"})
        return
    }

    userID,ok := userIDstr.(uint)
    if !ok {
        c.JSON(http.StatusUnauthorized,gin.H{"message":"unauthorized user failed to convert to uint"})
        return
    }

    // Get user Role
    userRoleStr,exists := c.Get("userRole")
    if !exists{
        c.JSON(http.StatusUnauthorized,gin.H{"message":"unauthorized user"})
        return
    }

    userRole,ok := userRoleStr.(int)
    if !ok {
        c.JSON(http.StatusUnauthorized,gin.H{"message":"unauthorized user failed to convert to uint"})
        return
    }

    var instance Student

    switch userRole {
    case 1:
        // get the student instance
        err := config.DB.First(&instance, id).Error
        if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Student Id invalid",
			})
			return
		}
        // verify if the student is an fme student 
        if !instance.Fmestudent {
            c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Cannot perfom this operation on this student",
			})
			return
        }

        // set the graduation status to true
        instance.GraduationStatus = true
        instance.GraduationDate = gradDate
        instance.NsqLevel = GraduateStudentSchema.NsqLevel

        err = config.DB.Save(&instance).Error
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "message": "Unable to update the user record",
            })
            return
        }

        // set the user role to artisan
        var userinstance myuser.User
        err = config.DB.First(&userinstance, instance.UserID).Error
        if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Student Id invalid",
			})
			return
		}
        userinstance.Role = 6

        err = config.DB.Save(&userinstance).Error
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "message": "Unable to update the user record",
            })
            return
        }

        // create and save the artisan
        artisans := artisans.Artisans{
            UserID: instance.UserID,
            LastName: instance.Lastname,
            FirstName: instance.Firstname,
            LGA: instance.LocalGovernment,
            StateOfResidence: instance.StateOfResidence,

        }

        err = config.DB.Save(&artisans).Error
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "message": "Unable to update the user record",
            })
            return
        } 


        // send successfull response
        c.JSON(http.StatusOK, gin.H{
            "message": "The student has been successfully graduated",
        })   
        
    case 2:  //for the mda user
        var userMdaId uint 
        err := config.DB.Table("mdas").
        Where("user_id = ?", userID).
        Pluck("id",&userMdaId).Error
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{
                "message": "MdaAccount has issues",
            })
            return
        }

        err = config.DB.First(&instance, id).Error
        if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Student Id invalid",
			})
			return
		}

        if userMdaId != instance.MdaID {
            c.JSON(http.StatusInternalServerError, gin.H{
                "message": "Unable to update the user record",
            })
            return
        }

        // update the students graduation status
        instance.GraduationStatus = true
        instance.GraduationDate = gradDate
        instance.NsqLevel = GraduateStudentSchema.NsqLevel

        err = config.DB.Save(&instance).Error
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "message": "Unable to update the user record",
            })
            return
        }
         // set the user role to artisan
         var userinstance myuser.User
         err = config.DB.First(&userinstance, instance.UserID).Error
         if err != nil {
             c.JSON(http.StatusBadRequest, gin.H{
                 "message": "Student Id invalid",
             })
             return
         }
         userinstance.Role = 6
 
         err = config.DB.Save(&userinstance).Error
         if err != nil {
             c.JSON(http.StatusInternalServerError, gin.H{
                 "message": "Unable to update the user record",
             })
             return
         }

        // create and save the artisan
        artisans := artisans.Artisans{
            UserID: instance.UserID,
            LastName: instance.Lastname,
            FirstName: instance.Firstname,
            LGA: instance.LocalGovernment,
            StateOfResidence: instance.StateOfResidence,

        }

        err = config.DB.Save(&artisans).Error
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "message": "Unable to update the user record",
            })
            return
        } 


        // send successfull response
        c.JSON(http.StatusOK, gin.H{
            "message": "The student has been successfully graduated",
        })   

    case 3:
        //get user stcid
		var userStcId uint
		err := config.DB.
				Table("stcs").
				Select("id").
				Where("user_id = ?", userID).
				Scan(&userStcId).Error
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Error trying to fetch data",
			})
			return
		}   

        err = config.DB.First(&instance, id).Error
        if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Student Id invalid",
			})
			return
		}

        if userStcId != instance.StcID {
            c.JSON(http.StatusInternalServerError, gin.H{
                "message": "Unable to update the user record",
            })
            return
        }

        // update the graduation status
        instance.GraduationStatus = true
        instance.GraduationDate = gradDate
        instance.NsqLevel = GraduateStudentSchema.NsqLevel

        err = config.DB.Save(&instance).Error
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "message": "Unable to update the user record",
            })
            return
        }

         // set the user role to artisan
         var userinstance myuser.User
         err = config.DB.First(&userinstance, instance.UserID).Error
         if err != nil {
             c.JSON(http.StatusBadRequest, gin.H{
                 "message": "Student Id invalid",
             })
             return
         }
         userinstance.Role = 6
 
         err = config.DB.Save(&userinstance).Error
         if err != nil {
             c.JSON(http.StatusInternalServerError, gin.H{
                 "message": "Unable to update the user record",
             })
             return
         }

        // create and save the artisan
        artisans := artisans.Artisans{
            UserID: instance.UserID,
            LastName: instance.Lastname,
            FirstName: instance.Firstname,
            LGA: instance.LocalGovernment,
            StateOfResidence: instance.StateOfResidence,

        }

        err = config.DB.Save(&artisans).Error
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "message": "Unable to update the user record",
            })
            return
        } 

        // send successfull response
        c.JSON(http.StatusOK, gin.H{
            "message": "The student has been successfully graduated",
        })  

    default:
        c.JSON(http.StatusUnauthorized, gin.H{
            "message": "Cannot perform this operation on this student",
        })  

    }
 
}
var disabilities []Disability = []Disability{
    {ID: 1, Name: "Blindness"},
    {ID: 2, Name: "Low vision"},
    {ID: 3, Name: "Color blindness"},
    {ID: 4, Name: "Deafness"},
    {ID: 5, Name: "Hard of hearing"},
    {ID: 6, Name: "Amputation"},
    {ID: 7, Name: "Cerebral palsy"},
    {ID: 8, Name: "Muscular dystrophy"},
    {ID: 9, Name: "Spinal cord injury"},
    {ID: 10, Name: "Spina bifida"},
    {ID: 11, Name: "Intellectual disability"},
    {ID: 12, Name: "Speech impairment"},
    {ID: 13, Name: "Language impairment"},
    {ID: 14, Name: "Depression"},
    {ID: 15, Name: "Anxiety"},
    {ID: 16, Name: "Bipolar disorder"},
    {ID: 17, Name: "Schizophrenia"},
}

func GetDisabilityList(c *gin.Context) {
    c.JSON(http.StatusOK, disabilities)
}

func validateDisability(name string) bool {
    for _, disability := range disabilities {
        if strings.EqualFold(disability.Name, name) {
            return true
        }
    }
    return false
}


func EditStudent(c *gin.Context) {
    // get the student id
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

    // bind the incoming schema
    if c.ShouldBind(&CreateStudentSchema) != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Failed to read request body",
        })
        return
    }

    // validate the incoming data - come back to it 
    // Data Validation

    var stateOfOrigin string = ""
    var result bool
    var dOB time.Time
    var stateOfResidence string = ""
    var gender string = ""
    if (CreateStudentSchema.StateOfOrigin != "") {

        stateOfOrigin, result = utilities.ValidateState(CreateStudentSchema.StateOfOrigin)
        if !result {
            c.JSON(http.StatusBadRequest, gin.H{
                "message": "Incorrect state of origin",
            })
            return
        }
    }
    
    if (CreateStudentSchema.StateOfResidence != "") {
        stateOfResidence, result = utilities.ValidateState(CreateStudentSchema.StateOfResidence)
        if !result {
            c.JSON(http.StatusBadRequest, gin.H{
                "message": "Incorrect state of residence",
            })
            return
        }
    }
    
    if (CreateStudentSchema.PhoneNumber != "") {
        if !utilities.IsNigerianPhoneNumber(CreateStudentSchema.PhoneNumber){
            c.JSON(http.StatusBadRequest, gin.H{
                "message": "Incorrect phone number",
            })
            return
        }
    }
    


    if (CreateStudentSchema.DOBstring != "") {
        dOB, err = utilities.ParseDoB(CreateStudentSchema.DOBstring)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{
                "message": "Incorrect date format",
            })
            return
        }
    }
    
    if (CreateStudentSchema.Gender != "") {
        gender, result = utilities.ValidateGender(CreateStudentSchema.Gender)
        if !result {
            c.JSON(http.StatusBadRequest, gin.H{
                "message": "Incorrect Gender",
            })
            return
        } 
    }
    
    if (CreateStudentSchema.NationalIdentityNumber != "") {
        result = utilities.VeriryNINFormat(CreateStudentSchema.NationalIdentityNumber)
        if !result {
            c.JSON(http.StatusBadRequest, gin.H{
                "message": "Incorrect NIN format",
            })
            return
        }
    }
    

    if CreateStudentSchema.IsDisabled {
        if !validateDisability(CreateStudentSchema.DisabilityName) {
            c.JSON(http.StatusBadRequest, gin.H{
                "message": "Wrong disability name",
            })
            return
        }
    }


    // get the user details 
    userIDstr,exists := c.Get("userID")
    if !exists{
        c.JSON(http.StatusUnauthorized,gin.H{"message":"unauthorized user"})
        return
    }

    userID,ok := userIDstr.(uint)
    if !ok {
        c.JSON(http.StatusUnauthorized,gin.H{"message":"unauthorized user failed to convert to uint"})
        return
    }

    // Get user Role
    userRoleStr,exists := c.Get("userRole")
    if !exists{
        c.JSON(http.StatusUnauthorized,gin.H{"message":"unauthorized user"})
        return
    }

    userRole,ok := userRoleStr.(int)
    if !ok {
        c.JSON(http.StatusUnauthorized,gin.H{"message":"unauthorized user failed to convert to uint"})
        return
    }

    var studentInstance Student
    var studentCourseInstance StudentCourse
    switch userRole {
    case 1:
        //fme case
        // get the student instance
        err := config.DB.First(&studentInstance, id).Error
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{
                "message": "Unable to fetch student instance",
            })
            return
        }

        // get the student course instance
        err = config.DB.First(&studentCourseInstance, id).Error
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{
                "message": "Unable to fetch student instance",
            })
            return
        }

        tx := config.DB.Begin()

        // set the new student fields
        if CreateStudentSchema.Firstname != "" {
            studentInstance.Firstname = CreateStudentSchema.Firstname
        }
        if CreateStudentSchema.Lastname != "" {
            studentInstance.Lastname = CreateStudentSchema.Lastname
        }
        if CreateStudentSchema.Gender != "" {
            studentInstance.Gender = gender
        }

        if CreateStudentSchema.PhoneNumber != "" {
            studentInstance.PhoneNumber = CreateStudentSchema.PhoneNumber
        }

        if CreateStudentSchema.StateOfOrigin != "" {
            studentInstance.StateOfOrigin = stateOfOrigin
        }

        if CreateStudentSchema.StateOfResidence != "" {
            studentInstance.StateOfResidence = stateOfResidence
        }

        if CreateStudentSchema.DOBstring != "" {
            studentInstance.DOB = dOB
        }

        if CreateStudentSchema.SID != "" {
            studentInstance.SID = CreateStudentSchema.SID
        }
        if CreateStudentSchema.CourseID != 0 {
            // change the student course details
            studentCourseInstance.CourseID = CreateStudentSchema.CourseID
        }

        if CreateStudentSchema.Address != "" {
            studentInstance.Address = CreateStudentSchema.Address
        }

        if CreateStudentSchema.NationalIdentityNumber != "" {
            studentInstance.NationalIdentityNumber = CreateStudentSchema.NationalIdentityNumber
        }

        if CreateStudentSchema.LocalGovernment != "" {
            studentInstance.LocalGovernment = CreateStudentSchema.LocalGovernment
        }

        if CreateStudentSchema.IsDisabled {
            studentInstance.IsDisabled = CreateStudentSchema.IsDisabled
            studentInstance.DisabilityName = CreateStudentSchema.DisabilityName
        } else {
            studentInstance.IsDisabled = CreateStudentSchema.IsDisabled
        }
        


        // save the new fields
        err = tx.Save(&studentInstance).Error
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "message": "Unable to update the student record",
            })
            return
        } 

        err = tx.Save(&studentCourseInstance).Error
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "message": "Unable to update the student course record",
            })
            return
        } 

        tx.Commit()

        // decide the new 
        c.JSON(http.StatusOK, gin.H{
            "message": "The student has been updated",
        })  

    case 2: // mda case
        // get the mda id 
        var userMdaId uint 
        err := config.DB.Table("mdas").
        Where("user_id = ?", userID).
        Pluck("id",&userMdaId).Error
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{
                "message": "MdaAccount has issues",
            })
            return
        }

        // get the student and the student course
        err = config.DB.First(&studentInstance, id).Error
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{
                "message": "Unable to fetch student instance",
            })
            return
        }

        // get the student course instance
        err = config.DB.First(&studentCourseInstance, id).Error
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{
                "message": "Unable to fetch student instance",
            })
            return
        }


        // verify if the student has the same mda id and make the changes
        if userMdaId != studentInstance.MdaID {
            c.JSON(http.StatusBadRequest, gin.H{
                "message": "time limit for changes exceeded",
            })
            return
        }

        // verify if the student can still be operated on 
        if time.Since(studentInstance.CreatedAt) > time.Hour {
            c.JSON(http.StatusBadRequest, gin.H{
				"message": "time to change details elapsed",
			})
			return
        }

        // make the changes 
        tx := config.DB.Begin()

        // set the new student fields
        if CreateStudentSchema.Firstname != "" {
            studentInstance.Firstname = CreateStudentSchema.Firstname
        }
        if CreateStudentSchema.Lastname != "" {
            studentInstance.Lastname = CreateStudentSchema.Lastname
        }
        if CreateStudentSchema.Gender != "" {
            studentInstance.Gender = gender
        }

        if CreateStudentSchema.PhoneNumber != "" {
            studentInstance.PhoneNumber = CreateStudentSchema.PhoneNumber
        }

        if CreateStudentSchema.StateOfOrigin != "" {
            studentInstance.StateOfOrigin = stateOfOrigin
        }

        if CreateStudentSchema.StateOfResidence != "" {
            studentInstance.StateOfResidence = stateOfResidence
        }

        if CreateStudentSchema.DOBstring != "" {
            studentInstance.DOB = dOB
        }

        if CreateStudentSchema.SID != "" {
            studentInstance.SID = CreateStudentSchema.SID
        }
        if CreateStudentSchema.CourseID != 0 {
            // change the student course details
            studentCourseInstance.CourseID = CreateStudentSchema.CourseID
        }

        if CreateStudentSchema.Address != "" {
            studentInstance.Address = CreateStudentSchema.Address
        }

        if CreateStudentSchema.NationalIdentityNumber != "" {
            studentInstance.NationalIdentityNumber = CreateStudentSchema.NationalIdentityNumber
        }

        if CreateStudentSchema.LocalGovernment != "" {
            studentInstance.LocalGovernment = CreateStudentSchema.LocalGovernment
        }

        if CreateStudentSchema.IsDisabled {
            studentInstance.IsDisabled = CreateStudentSchema.IsDisabled
            studentInstance.DisabilityName = CreateStudentSchema.DisabilityName
        } else {
            studentInstance.IsDisabled = CreateStudentSchema.IsDisabled
        }
        


        // save the new fields
        err = tx.Save(&studentInstance).Error
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "message": "Unable to update the student record",
            })
            return
        } 

        err = tx.Save(&studentCourseInstance).Error
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "message": "Unable to update the student course record",
            })
            return
        } 

        tx.Commit()

        // decide the new 
        c.JSON(http.StatusOK, gin.H{
            "message": "The student has been updated",
        }) 

    case 3:
        // get the mda id 
        var userStcId uint 
        err := config.DB.Table("stcs").
        Where("user_id = ?", userID).
        Pluck("id",&userStcId).Error
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{
                "message": "MdaAccount has issues",
            })
            return
        }

        // get the student and the student course
        err = config.DB.First(&studentInstance, id).Error
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{
                "message": "Unable to fetch student instance",
            })
            return
        }

        // get the student course instance
        err = config.DB.First(&studentCourseInstance, id).Error
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{
                "message": "Unable to fetch student instance",
            })
            return
        }


        // verify if the student has the same mda id and make the changes
        if userStcId != studentInstance.StcID {
            c.JSON(http.StatusBadRequest, gin.H{
                "message": "time limit for changes exceeded",
            })
            return
        }

        // verify if the student can still be operated on 
        if time.Since(studentInstance.CreatedAt) > time.Hour {
            c.JSON(http.StatusBadRequest, gin.H{
				"message": "time to change details elapsed",
			})
			return
        }

        // make the changes 
        tx := config.DB.Begin()

        // set the new student fields
        if CreateStudentSchema.Firstname != "" {
            studentInstance.Firstname = CreateStudentSchema.Firstname
        }
        if CreateStudentSchema.Lastname != "" {
            studentInstance.Lastname = CreateStudentSchema.Lastname
        }
        if CreateStudentSchema.Gender != "" {
            studentInstance.Gender = gender
        }

        if CreateStudentSchema.PhoneNumber != "" {
            studentInstance.PhoneNumber = CreateStudentSchema.PhoneNumber
        }

        if CreateStudentSchema.StateOfOrigin != "" {
            studentInstance.StateOfOrigin = stateOfOrigin
        }

        if CreateStudentSchema.StateOfResidence != "" {
            studentInstance.StateOfResidence = stateOfResidence
        }

        if CreateStudentSchema.DOBstring != "" {
            studentInstance.DOB = dOB
        }

        if CreateStudentSchema.SID != "" {
            studentInstance.SID = CreateStudentSchema.SID
        }
        if CreateStudentSchema.CourseID != 0 {
            // change the student course details
            studentCourseInstance.CourseID = CreateStudentSchema.CourseID
        }

        if CreateStudentSchema.Address != "" {
            studentInstance.Address = CreateStudentSchema.Address
        }

        if CreateStudentSchema.NationalIdentityNumber != "" {
            studentInstance.NationalIdentityNumber = CreateStudentSchema.NationalIdentityNumber
        }

        if CreateStudentSchema.LocalGovernment != "" {
            studentInstance.LocalGovernment = CreateStudentSchema.LocalGovernment
        }

        if CreateStudentSchema.IsDisabled {
            studentInstance.IsDisabled = CreateStudentSchema.IsDisabled
            studentInstance.DisabilityName = CreateStudentSchema.DisabilityName
        } else {
            studentInstance.IsDisabled = CreateStudentSchema.IsDisabled
        }
        


        // save the new fields
        err = tx.Save(&studentInstance).Error
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "message": "Unable to update the student record",
            })
            return
        } 

        err = tx.Save(&studentCourseInstance).Error
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "message": "Unable to update the student course record",
            })
            return
        } 

        tx.Commit()

        // decide the new 
        c.JSON(http.StatusOK, gin.H{
            "message": "The student has been updated",
        })
    default:
        c.JSON(http.StatusUnauthorized, gin.H{
            "message": "You are unable to perform this action",
        })

    }
}