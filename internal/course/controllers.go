package course

import (
	"fme_backend/internal/config"
	"net/http"

	"github.com/gin-gonic/gin"
)



func CreateCourse(c *gin.Context) {
	if c.Bind(&CreateCourseSchema) != nil {
		c.JSON(http.StatusBadRequest,gin.H{
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
		Name: CreateCourseSchema.Name,
		Description: CreateCourseSchema.Description,
	}

	result:= config.DB.Create(&course)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create course",
		})
		return
	}

	c.JSON(200, gin.H{"message": "Course created successfully"})
}


func CreateSector(c *gin.Context) {
	if c.Bind(&CreateSectorSchema) != nil {
		c.JSON(http.StatusBadRequest,gin.H{
			"message": "Bad request body",
		})
	}

	var sectorCheck Course
	config.DB.Where("name= ?", CreateCourseSchema.Name).First(&sectorCheck)
	if sectorCheck.ID != 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Sector already exists",
		})
		return
	}

	sector := Sector{
		Name: CreateSectorSchema.Name,
		Description: CreateSectorSchema.Description,
	}

	result:= config.DB.Create(&sector)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create sector",
		})
		return
	}

	c.JSON(200, gin.H{"message": "Sector created successfully"})
}

func GetCourse(c *gin.Context){
	// GET ID
	id:= c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "path parameter not provided",
		})
		return
	}

	// FIND IF THE COURSE EXIST
	var instance Course
	instance_result := config.DB.Select("id","name","description").First(&instance, id)
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
	result := config.DB.Select("id","name","description").Where("name LIKE ?", "%"+name+"%").Find(&courses)
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
  result := config.DB.Select("id","name","description").Find(&courses)
  if result.Error != nil {
    c.JSON(http.StatusNotFound, gin.H{"message":"no course found"})
    return
  }

  totalCourses:= config.DB.Model(&Course{}).Count(&count)
	if totalCourses.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "unable to get total courses"})
		return
	}

  // Respond with the list of courses
  c.JSON(200, gin.H{
	"courses":courses,
	"total_courses":count,
  })
}

func GetSector(c *gin.Context) {
	// GET ID
	id:= c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "path parameter not provided",
		})
		return
	}

	// FIND IF THE COURSE EXIST
	var instance Sector
	instance_result := config.DB.Select("id","name","description").First(&instance, id)
	if instance_result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "instance does not exist",
		})
		return
	}

	c.JSON(http.StatusOK, instance)
}

func GetAllSectors(c *gin.Context) {
	var sectors []Sector
	var count int64

	// Find all courses from the database
	result := config.DB.Select("id","name","description").Find(&sectors)
	if result.Error != nil {
	  c.JSON(http.StatusNotFound, gin.H{"message":"no course found"})
	  return
	}
	totalCourses:= config.DB.Model(&Sector{}).Count(&count)
	if totalCourses.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "unable to get total sectors"})
		return
	}
  
	// Respond with the list of courses
	c.JSON(200, gin.H{
		"sectors":sectors,
		"total_sectors": count,
	})
}
