package student

import (
	"encoding/csv"
	"errors"
	"fme_backend/internal/config"
	"fme_backend/internal/course"
	myuser "fme_backend/internal/user"
	"fme_backend/internal/utilities"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
        NsqLevel: CreateStudentSchema.NsqLevel,
        Address: CreateStudentSchema.Address,
        PhoneNumber: CreateStudentSchema.PhoneNumber,
        NationalIdentityNumber: CreateStudentSchema.NationalIdentityNumber,
        LocalGovernment: CreateStudentSchema.LocalGovernment,
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
    studentcourse := course.StudentCourse{
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
        NsqLevel: CreateStudentSchema.NsqLevel,
        Address: CreateStudentSchema.Address,
        LocalGovernment: CreateStudentSchema.LocalGovernment,
        NationalIdentityNumber: CreateStudentSchema.NationalIdentityNumber,
    }
    fmt.Println(student)
    studentresult := tx.Create(&student) // Create student within transaction
    if studentresult.Error != nil {
        tx.Rollback() // Rollback if student creation fails
        c.JSON(http.StatusBadRequest, gin.H{
            "message": "Failed to create User",
        })
        return
    }

    // Add Course
    studentcourse := course.StudentCourse{
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
        NsqLevel: CreateStudentSchema.NsqLevel,
        Address: CreateStudentSchema.Address,
        PhoneNumber: CreateStudentSchema.PhoneNumber,
        LocalGovernment: CreateStudentSchema.LocalGovernment,
        NationalIdentityNumber: CreateStudentSchema.NationalIdentityNumber,
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
    studentcourse := course.StudentCourse{
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
    if len(record) < 14 {
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

    courseID, err := strconv.ParseUint(record[11], 10, 32)
    if err != nil {
        return CreateStudentSchematype{}, fmt.Errorf("invalid course ID: %v", err)
    }

    if !utilities.VeriryNINFormat(record[12]) {
        return CreateStudentSchematype{}, fmt.Errorf("invalid NIN format: %v", err)
        
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
        NsqLevel:         record[9],
        Address:          record[10],
        CourseID:         uint(courseID),
        NationalIdentityNumber: record[12],
        LocalGovernment: record[13],
    }, nil
}

func createMdaStudentInstance(mdaID uint, schema CreateStudentSchematype) (Student, error) {

    dob, err := utilities.ParseDoB(schema.DOBstring)
    if err != nil {
        return Student{}, errors.New("wrong date format")
    }
    tx := config.DB.Begin()

    // Create user with transaction
    result, message, newUserID := myuser.CreateStudentUser(tx, schema.Email, "dfcv")
    if !result {
        tx.Rollback()
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
        NsqLevel:         schema.NsqLevel,
        Address:          schema.Address,
    }

    if err := tx.Create(&student).Error; err != nil {
        tx.Rollback()
        return Student{}, fmt.Errorf("failed to create student: %v", err)
    }

    studentCourse := course.StudentCourse{
        StudentID: student.ID,
        CourseID:  schema.CourseID,
    }

    if err := tx.Create(&studentCourse).Error; err != nil {
        tx.Rollback()
        return Student{}, fmt.Errorf("failed to create student course: %v", err)
    }

    tx.Commit()
    return student, nil
}

func createStcStudentInstance(stcID uint, schema CreateStudentSchematype) (Student, error) {

    dob, err := utilities.ParseDoB(schema.DOBstring)
    if err != nil {
        return Student{}, errors.New("wrong date format")
    }
    tx := config.DB.Begin()

    // Create user with transaction
    result, message, newUserID := myuser.CreateStudentUser(tx, schema.Email, "dfcv")
    if !result {
        tx.Rollback()
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
        NsqLevel:         schema.NsqLevel,
        Address:          schema.Address,
    }

    if err := tx.Create(&student).Error; err != nil {
        tx.Rollback()
        return Student{}, fmt.Errorf("failed to create student: %v", err)
    }

    studentCourse := course.StudentCourse{
        StudentID: student.ID,
        CourseID:  schema.CourseID,
    }

    if err := tx.Create(&studentCourse).Error; err != nil {
        tx.Rollback()
        return Student{}, fmt.Errorf("failed to create student course: %v", err)
    }

    tx.Commit()
    return student, nil
}


func createFmeStudentInstance(schema CreateStudentSchematype) (Student, error) {

    dob, err := utilities.ParseDoB(schema.DOBstring)
    if err != nil {
        return Student{}, errors.New("wrong date format")
    }
    tx := config.DB.Begin()

    // Create user with transaction
    result, message, newUserID := myuser.CreateStudentUser(tx, schema.Email, "dfcv")
    if !result {
        tx.Rollback()
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
        NsqLevel:         schema.NsqLevel,
        Address:          schema.Address,
    }

    if err := tx.Create(&student).Error; err != nil {
        tx.Rollback()
        return Student{}, fmt.Errorf("failed to create student: %v", err)
    }

    studentCourse := course.StudentCourse{
        StudentID: student.ID,
        CourseID:  schema.CourseID,
    }

    if err := tx.Create(&studentCourse).Error; err != nil {
        tx.Rollback()
        return Student{}, fmt.Errorf("failed to create student course: %v", err)
    }

    tx.Commit()
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

    for {
        record, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{
                "message": "Failed to read CSV file",
            })
            return
        }

        // Parse and validate each record
        studentSchema, err := parseAndValidateRecord(record)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{
                "message": fmt.Sprintf("Invalid data: %v", err),
            })
            return
        }

        // Create user and student instances
        student, err := createMdaStudentInstance(mdaID, studentSchema)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{
                "message": fmt.Sprintf("Failed to create student: %v", err),
            })
            return
        }

        students = append(students, student)
    }

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

    for {
        record, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{
                "message": "Failed to read CSV file",
            })
            return
        }

        // Parse and validate each record
        studentSchema, err := parseAndValidateRecord(record)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{
                "message": fmt.Sprintf("Invalid data: %v", err),
            })
            return
        }

        // Create user and student instances
        student, err := createFmeStudentInstance(studentSchema)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{
                "message": fmt.Sprintf("Failed to create student: %v", err),
            })
            return
        }

        students = append(students, student)
    }

    // If all students are successfully created
    c.JSON(http.StatusOK, gin.H{
        "message": "Students created successfully",
        "students": students,
    })
}



func CreateStcStudentFromCsv(c *gin.Context) {
    // Retrieve the MDA ID
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

    for {
        record, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{
                "message": "Failed to read CSV file",
            })
            return
        }

        // Parse and validate each record
        studentSchema, err := parseAndValidateRecord(record)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{
                "message": fmt.Sprintf("Invalid data: %v", err),
            })
            return
        }

        // Create user and student instances
        student, err := createStcStudentInstance(stcID, studentSchema)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{
                "message": fmt.Sprintf("Failed to create student: %v", err),
            })
            return
        }

        students = append(students, student)
    }

    // If all students are successfully created
    c.JSON(http.StatusOK, gin.H{
        "message": "Students created successfully",
        "students": students,
    })
}





