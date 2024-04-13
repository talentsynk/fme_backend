package course

import (
	"fme_backend/internal/config"
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
			"error": "Failed to create course",
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
			"error": "Failed to create Category",
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

	// FIND IF THE COURSE EXIST
	var instance Course
	instance_result := config.DB.Select("id", "name", "description").First(&instance, id)
	if instance_result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "instance does not exist",
		})
		return
	}

	c.JSON(http.StatusOK, instance)
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
	var courses []Course
	var count int64

	// Find all courses from the database
	result := config.DB.Select("id", "name", "description").Find(&courses)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "no course found"})
		return
	}

	totalCourses := config.DB.Model(&Course{}).Count(&count)
	if totalCourses.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "unable to get total courses"})
		return
	}

	// Respond with the list of courses
	c.JSON(200, gin.H{
		"courses":       courses,
		"total_courses": count,
	})
}

func GetCategory(c *gin.Context) {
	// GET ID
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "path parameter not provided",
		})
		return
	}

	// FIND IF THE COURSE EXIST
	var instance Category
	instance_result := config.DB.Select("id", "name", "description").First(&instance, id)
	if instance_result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Category does not exist",
		})
		return
	}

	c.JSON(http.StatusOK, instance)
}

func GetAllCategories(c *gin.Context) {
	var categories []Category
	var count int64

	// Find all courses from the database
	result := config.DB.Select("id", "name", "description").Find(&categories)
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
