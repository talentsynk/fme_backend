package student

import (
	"fme_backend/internal/config"
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
            "error": "Failed to read request body",
        })
        return
    }

    // Data Validation
    stateOfOrigin, result := utilities.ValidateState(CreateStudentSchema.StateOfOrigin)
    if !result {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Incorrect state of origin",
        })
        return
    }

    StateOfResidence, result := utilities.ValidateState(CreateStudentSchema.StateOfResidence)
    if !result {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Incorrect state of residence",
        })
        return
    }

    dOB, err := utilities.ParseDoB(CreateStudentSchema.DOBstring)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Incorrect date format",
        })
        return
    }

    gender, result := utilities.ValidateGender(CreateStudentSchema.Gender)
    if !result {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Incorrect Gender",
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
            "error": message,
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
    }

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

func CreateMdaStudent(c *gin.Context) {
	// retrieve the mda id
	mdaIDStr,exists := c.Get("mdaID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Problem with the aithorization token",
		})
		return
	}
	mdaID,ok := mdaIDStr.(uint)
	if !ok{
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Problem with the aithorization token",
		})
		return
	}

	
	// bind the post data
	if c.Bind(&CreateStudentSchema) != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Failed to read request body",
        })
        return
    }

    // Data Validation
    stateOfOrigin, result := utilities.ValidateState(CreateStudentSchema.StateOfOrigin)
    if !result {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Incorrect state of origin",
        })
        return
    }

    StateOfResidence, result := utilities.ValidateState(CreateStudentSchema.StateOfResidence)
    if !result {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Incorrect state of residence",
        })
        return
    }

    dOB, err := utilities.ParseDoB(CreateStudentSchema.DOBstring)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Incorrect date format",
        })
        return
    }

    gender, result := utilities.ValidateGender(CreateStudentSchema.Gender)
    if !result {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Incorrect Gender",
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
            "error": message,
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

func CreateStcStudent(c *gin.Context) {
	// retrieve the stc id
	stcIDStr,exists := c.Get("stcID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Problem with the aithorization token",
		})
		return
	}
	stcID,ok := stcIDStr.(uint)
	if !ok{
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Problem with the aithorization token",
		})
		return
	}

	

	if c.Bind(&CreateStudentSchema) != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Failed to read request body",
        })
        return
    }

    // Data Validation
    stateOfOrigin, result := utilities.ValidateState(CreateStudentSchema.StateOfOrigin)
    if !result {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Incorrect state of origin",
        })
        return
    }

    StateOfResidence, result := utilities.ValidateState(CreateStudentSchema.StateOfResidence)
    if !result {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Incorrect state of residence",
        })
        return
    }

    dOB, err := utilities.ParseDoB(CreateStudentSchema.DOBstring)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Incorrect date format",
        })
        return
    }

    gender, result := utilities.ValidateGender(CreateStudentSchema.Gender)
    if !result {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Incorrect Gender",
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
            "error": message,
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
	limitStr := c.Query("limit")
    pageStr := c.Query("page")
    mdaIDStr := c.Query("mda_id")
    stcIDStr := c.Query("stc_id")

    limit, err := strconv.Atoi(limitStr)
    if err != nil || limit <= 0 {
        limit = 10 // Default limit
    }

    page, err := strconv.Atoi(pageStr)
    if err != nil || page <= 0 {
        page = 1 // Default page
    }

    offset := (page - 1) * limit

	

    // Build the query with optional filtering
    query := config.DB.Select("id","firstname","lastname","dob","state_of_origin","state_of_residence","gender","graduation_status").Offset(offset).Limit(limit)
    if mdaIDStr != "" {
        mdaID, err := strconv.Atoi(mdaIDStr)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "invalid mda_id"})
            return
        }
        query = query.Where("mda_id = ?", mdaID)
    } else if stcIDStr != "" {
        stcID, err := strconv.Atoi(stcIDStr)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "invalid stc_id"})
            return
        }
        query = query.Where("stc_id = ?", stcID)
    }
	var students []Student
	result := query.Find(&students)
    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "students": students,
    })
}

func GetStudent(c *gin.Context) {
	// GET ID
	idStr:= c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "path parameter not provided",
		})
		return
	}
	// Convert id to string
	id,err :=strconv.Atoi(idStr)
	if err!= nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "path parameter invalid",
		})
		return
	}

	var student Student
	instance_result := config.DB.Select("id","firstname","lastname","dob","state_of_origin","state_of_residence","gender","graduation_status").First(&student, id)
	fmt.Println(id)

	if instance_result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "instance does not exist",
		})
		return
	}

	c.JSON(http.StatusOK, student)
}
