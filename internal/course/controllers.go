package course

import (
	"fme_backend/internal/config"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateCourse(c *gin.Context) {
	if c.Bind(&CreateCourseSchema) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Bad request body",
		})
	}

	var courseCheck Course
	config.DB.Where("name= ?", CreateCourseSchema.Name).First(&courseCheck)
	if courseCheck.ID != 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Course already exists",
		})
		return
	}

	course := Course{
		Name:        CreateCourseSchema.Name,
		Description: CreateCourseSchema.Description,
		CategoryID: CreateCourseSchema.CategoryID,
	}

	result := config.DB.Create(&course)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to create course",
		})
		return
	}

	c.JSON(200, gin.H{"message": "Course created successfully"})
}

func CreateCategory(c *gin.Context) {
	if c.Bind(&CreateCategorySchema) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Bad request body",
		})
	}

	var categoryCheck Category
	config.DB.Where("name= ?", CreateCategorySchema.Name).First(&categoryCheck)
	if categoryCheck.ID != 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Category already exists",
		})
		return
	}

	category := Category{
		Name:        CreateCategorySchema.Name,
		Description: CreateCategorySchema.Description,
	}

	result := config.DB.Create(&category)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to create Category",
		})
		return
	}

	c.JSON(200, gin.H{"message": "Category created successfully"})
}

func GetCourse(c *gin.Context) {
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

	 var course GetCourseSchema
	 switch (userRole) {
	 case 1:
		err := config.DB.Table("courses").
        Select("COUNT(DISTINCT students.id) AS total_students, COUNT(DISTINCT mda_courses.mda_id) AS total_mda, COUNT(DISTINCT stc_courses.stc_id) AS total_stc, courses.description AS description, courses.name AS name, courses.id AS id").
        Joins("LEFT JOIN student_courses ON courses.id = student_courses.course_id").
        Joins("LEFT JOIN students ON student_courses.student_id = students.id").
        Joins("LEFT JOIN mda_courses ON courses.id = mda_courses.course_id").
        Joins("LEFT JOIN stc_courses ON courses.id = stc_courses.course_id").
        Where("courses.id = ?", id).
        Group("courses.id").
        Scan(&course).Error
		 if err!=nil {
			 c.JSON(http.StatusBadRequest,gin.H{
				 "message":"error retrieving students",
			 })
			 return
		 }
		 c.JSON(http.StatusOK,gin.H{"course":course})
 
	 case 2: 
	 err := config.DB.Table("courses").
	 Select("COUNT(DISTINCT students.id) AS total_students, COUNT(DISTINCT mda_courses.mda_id) AS total_mda, COUNT(DISTINCT stc_courses.stc_id) AS total_stc, courses.description AS description, courses.name AS name, courses.id AS id").
	 Joins("LEFT JOIN student_courses ON courses.id = student_courses.course_id").
	 Joins("LEFT JOIN students ON student_courses.student_id = students.id").
	 Joins("LEFT JOIN mda_courses ON courses.id = mda_courses.course_id").
	 Joins("LEFT JOIN stc_courses ON courses.id = stc_courses.course_id").
	 Where("courses.id = ?", id).
	 Group("courses.id").
	 Scan(&course).Error
	  if err!=nil {
		  c.JSON(http.StatusBadRequest,gin.H{
			  "message":"error retrieving students",
		  })
		  return
	  }
	  c.JSON(http.StatusOK,gin.H{"course":course})



	case 3:
		err := config.DB.Table("courses").
        Select("COUNT(DISTINCT students.id) AS total_students, COUNT(DISTINCT mda_courses.mda_id) AS total_mda, COUNT(DISTINCT stc_courses.stc_id) AS total_stc, courses.description AS description, courses.name AS name, courses.id AS id").
        Joins("LEFT JOIN student_courses ON courses.id = student_courses.course_id").
        Joins("LEFT JOIN students ON student_courses.student_id = students.id").
        Joins("LEFT JOIN mda_courses ON courses.id = mda_courses.course_id").
        Joins("LEFT JOIN stc_courses ON courses.id = stc_courses.course_id").
        Where("courses.id = ?", id).
        Group("courses.id").
        Scan(&course).Error
		 if err!=nil {
			 c.JSON(http.StatusBadRequest,gin.H{
				 "message":"error retrieving students",
			 })
			 return
		 }
		 c.JSON(http.StatusOK,gin.H{"course":course})
		 
	 default:
		fmt.Println(userID)
		 c.JSON(http.StatusUnauthorized,gin.H{"message":"default unauthorized user"})
		 return
	 }
}

func GetCoursesByName(c *gin.Context) {
	name := c.Query("name") // Get the query string parameter named "name"

	var courses []Course

	// Build the query with filtering based on the name parameter
	result := config.DB.Select("id", "name", "description").Where("name LIKE ?", "%"+name+"%").Find(&courses)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": result.Error.Error()})
		return
	}

	// Respond with the list of matching courses
	c.JSON(200, courses)
}

func GetAllCourses(c *gin.Context) {
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

	var courses []GetAllCoursesSchema
	switch (userRole) {
	case 1:
	   err := config.DB.Table("courses").
	   Select("id,name,description").
	   Find(&courses).Error
		if err!=nil {
			c.JSON(http.StatusBadRequest,gin.H{
				"message":"error retrieving students",
			})
			return
		}
		c.JSON(http.StatusOK,gin.H{"course":courses})

	case 2:
		err := config.DB.Table("courses").
	   Select("id,name,description").
	   Find(&courses).Error
		if err!=nil {
			c.JSON(http.StatusBadRequest,gin.H{
				"message":"error retrieving students",
			})
			return
		}
		c.JSON(http.StatusOK,gin.H{"course":courses})

	case 3:
		err := config.DB.Table("courses").
	   Select("id,name,description").
	   Find(&courses).Error
		if err!=nil {
			c.JSON(http.StatusBadRequest,gin.H{
				"message":"error retrieving students",
			})
			return
		}
		c.JSON(http.StatusOK,gin.H{"course":courses})

	default:
		fmt.Println(userID)
		c.JSON(http.StatusUnauthorized,gin.H{"message":"default unauthorized user"})
		return
	}
}

func GetAllCategories(c *gin.Context) {
	var categories []GetAllCategoriesSchema
	var count int64

	// Find all courses from the database
	result := config.DB.Table("categories").
	Select("categories.id AS id, categories.name AS name, categories.description AS description").
	Find(&categories)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "no course found"})
		return
	}
	totalCourses := config.DB.Model(&Category{}).Count(&count)
	if totalCourses.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "unable to get total Categories"})
		return
	}

	// Respond with the list of courses
	c.JSON(200, gin.H{
		"Categories":       categories,
		"total_Categories": count,
	})
}

func GetCategory(c *gin.Context) {

}

func GetCourseDetails(c *gin.Context) {
	// get the course id
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

	var courseStats struct {
		CourseName      string `json:"course_name"`
		CourseID        uint   `json:"course_id"`
		TotalStudents   int64  `json:"total_students"`
		UniqueMDACount  int64  `json:"unique_mda_count"`
		UniqueSTCCount  int64  `json:"unique_stc_count"`
		GraduatedCount  int64  `json:"graduated_count"`
	}
	

	result := config.DB.Table("courses").
	Select(`
		courses.name AS course_name,
		student_courses.course_id AS course_id,
		COUNT(student_courses.student_id) AS total_students,
		COUNT(DISTINCT CASE WHEN COALESCE(students.mda_id, stcs.mda_id) <> 0 THEN COALESCE(students.mda_id, stcs.mda_id) END) AS unique_mda_count,
		COUNT(DISTINCT CASE WHEN students.stc_id <> 0 THEN students.stc_id END) AS unique_stc_count,
		COUNT(CASE WHEN students.graduation_status = true THEN student_courses.student_id END) AS graduated_count
	`).
	Joins("LEFT JOIN student_courses ON student_courses.course_id = courses.id").
	Joins("LEFT JOIN students ON student_courses.student_id = students.id").
	Joins("LEFT JOIN stcs ON students.stc_id = stcs.id").
	Where("courses.id = ?",id).
	Group("student_courses.course_id,courses.name").
	Find(&courseStats)

    if result.Error != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
        return
    }

    c.JSON(http.StatusOK, courseStats)
}

func GetTopMDAsByStudentCountForCourse(c *gin.Context) {
	// Get course ID from path
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

	// Struct to hold results
	var topMDAs []struct {
		MDAName       string `json:"mda_name"`
		MDAID         uint   `json:"mda_id"`
		StudentCount  int64  `json:"student_count"`
	}

	// Query
	result := config.DB.Table("mdas").
		Select(`
			mdas.register_name AS mda_name,
			mdas.id AS mda_id,
			COUNT(DISTINCT student_courses.student_id) AS student_count
		`).
		Joins("LEFT JOIN students ON students.mda_id = mdas.id OR students.stc_id IN (SELECT id FROM stcs WHERE stcs.mda_id = mdas.id)").
		Joins("JOIN student_courses ON students.id = student_courses.student_id").
		Where("student_courses.course_id = ?", id).
		Group("mdas.id").
		Order("student_count DESC").
		Limit(5).
		Find(&topMDAs)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, topMDAs)
}

func GetTopSTCsByStudentCountForCourse(c *gin.Context) {
	// Get course ID from path
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "course ID path parameter not provided"})
		return
	}

	// Convert course ID to int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "course ID path parameter invalid"})
		return
	}

	// Struct to hold results
	var topSTCs []struct {
		STCName       string `json:"stc_name"`
		STCID         uint   `json:"stc_id"`
		StudentCount  int64  `json:"student_count"`
	}

	// Query
	result := config.DB.Table("stcs").
		Select(`
			stcs.name AS stc_name,
			stcs.id AS stc_id,
			COUNT(DISTINCT student_courses.student_id) AS student_count
		`).
		Joins("JOIN students ON students.stc_id = stcs.id").
		Joins("JOIN student_courses ON students.id = student_courses.student_id").
		Where("student_courses.course_id = ?", id).
		Group("stcs.id").
		Order("student_count DESC").
		Limit(5).
		Find(&topSTCs)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, topSTCs)
}




