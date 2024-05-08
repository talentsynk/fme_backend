package dashboard

import (
	"fme_backend/internal/config"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetDashSummary(c *gin.Context) {
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

	switch userRole {
	case 1:
		var totals struct {
			TotalStcs    int64
			TotalMdas    int64
			TotalStudents int64
		}
		
		
		err := config.DB.Table("stcs").Count(&totals.TotalStcs).
			Table("mdas").Count(&totals.TotalMdas).
			Table("students").Count(&totals.TotalStudents).Error

		if err != nil {
			c.JSON(http.StatusBadRequest,gin.H{"message":"unable to get information"})
			return
		}
		c.JSON(http.StatusOK,gin.H{"response":totals})

	default:
		c.JSON(http.StatusUnauthorized,gin.H{"message":"unauthorized user role"})
	}

}

func GetStudentPercentPerCourse(c *gin.Context) {
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

	var coursePercentages []CoursePercentage

	switch userRole {
	case 1:
		err := config.DB.Table("courses").
		Select("courses.name AS course_name, COUNT(student_courses.id) AS student_count, COUNT(student_courses.id) * 100.0 / SUM(COUNT(student_courses.id)) OVER() AS total_percent").
		Joins("LEFT JOIN student_courses ON courses.id = student_courses.course_id").
		Group("courses.id").
		Order("total_percent DESC").
		Scan(&coursePercentages).Error

		if err!=nil {
            c.JSON(http.StatusBadRequest,gin.H{
                "message":"error retrieving info",
            })
            return
        }

        c.JSON(http.StatusOK,gin.H{"coursePercentages":coursePercentages})

	default:
        fmt.Println(userID)
        c.JSON(http.StatusUnauthorized,gin.H{"message":"default unauthorized user"})
        return
	}


}