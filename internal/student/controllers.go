package student

import (
	"fme_backend/internal/config"
	"fme_backend/internal/course"
	myuser "fme_backend/internal/user"
	"fme_backend/internal/utilities"
	"fmt"
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
    stateOfOrigin, result := utilities.ValidateState(CreateStudentSchema.StateOfOrigin)
    if !result {
        c.JSON(http.StatusBadRequest, gin.H{
            "message": "Incorrect state of origin",
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

    // Transaction Handling
    tx := config.DB.Begin() // Begin a transaction

    // Create user with transaction
    result, message, newUserID := myuser.CreateStudentUser(tx, CreateStudentSchema.PhoneNumber, CreateStudentSchema.Email, "dfcv")
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
    stateOfOrigin, result := utilities.ValidateState(CreateStudentSchema.StateOfOrigin)
    if !result {
        c.JSON(http.StatusBadRequest, gin.H{
            "message": "Incorrect state of origin",
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

    // Transaction Handling
    tx := config.DB.Begin() // Begin a transaction

    // Create user with transaction
    result, message, newUserID := myuser.CreateStudentUser(tx, CreateStudentSchema.PhoneNumber, CreateStudentSchema.Email, "dfcv")
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

    tx.Commit() 
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
    stateOfOrigin, result := utilities.ValidateState(CreateStudentSchema.StateOfOrigin)
    if !result {
        c.JSON(http.StatusBadRequest, gin.H{
            "message": "Incorrect state of origin",
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

    // Transaction Handling
    tx := config.DB.Begin() // Begin a transaction

    // Create user with transaction
    result, message, newUserID := myuser.CreateStudentUser(tx, CreateStudentSchema.PhoneNumber, CreateStudentSchema.Email, "dfcv")
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
    }
    fmt.Println(student)
    studentresult := tx.Create(&student) // Create student within transaction
    if studentresult.Error != nil {
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

func GetAllStudents(c *gin.Context) {
    pageStr := c.Query("page")
    // mdaIDStr := c.Query("mda_id")
    // stcIDStr := c.Query("stc_id")

    limit:= 100

    page, err := strconv.Atoi(pageStr)
    if err != nil || page <= 0 {
        page = 1 // Default page
    }

    offset := (page - 1) * limit
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
    var students []GetAllStudentSchema
    switch (userID) {
    case 1:
        err := config.DB.Table("students").
        Select("students.id AS student_id, students.firstname AS first_name, students.state_of_residence AS state_of_residence, students.lastname AS last_name, users.is_active AS is_active, users.email AS email, GROUP_CONCAT(courses.name SEPARATOR ', ') AS courses_taken").
        Joins("JOIN users ON students.user_id = users.id").
        Joins("LEFT JOIN student_courses ON students.id = student_courses.student_id").
        Joins("LEFT JOIN courses ON student_courses.course_id = courses.id").
        Group("students.id").
        Offset(offset).
        Limit(limit).
        Scan(&students).Error
        if err!=nil {
            c.JSON(http.StatusBadRequest,gin.H{
                "message":"error retrieving students",
            })
            return
        }
        c.JSON(http.StatusOK,gin.H{"students":students})

    default:
        c.JSON(http.StatusUnauthorized,gin.H{"message":"default unauthorized user"})
        return
    }

}

func GetStudent(c *gin.Context) {
    // GET ID
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
    var student GetStudentSchema
    switch (userID) {
    case 1:
        err := config.DB.Table("students").
        Select("students.id AS student_id, students.firstname AS first_name, students.lastname AS last_name, users.is_active AS is_active, users.email AS email, GROUP_CONCAT(courses.name SEPARATOR ', ') AS courses_taken, students.gender AS gender, students.state_of_residence AS state_of_residence, students.address AS address, students.created_at AS created_at").
        Joins("JOIN users ON students.user_id = users.id").
        Joins("LEFT JOIN student_courses ON students.id = student_courses.student_id").
        Joins("LEFT JOIN courses ON student_courses.course_id = courses.id").
        Where("students.id = ?", id).
        Group("students.id").
        Scan(&student).Error
        if err!=nil {
            c.JSON(http.StatusBadRequest,gin.H{
                "message":"error retrieving students",
            })
            return
        }
        c.JSON(http.StatusOK,gin.H{"students":student})

    default:
        c.JSON(http.StatusUnauthorized,gin.H{"message":"default unauthorized user"})
        return
    }
}



