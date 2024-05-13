package dashboard

import (
	"fme_backend/internal/config"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetDashSummary(c *gin.Context) {
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

        var summary struct {
            GraduatedCount      int
            NonGraduatedCount   int
            STCsCount          int64
        }

        err = config.DB.Table("students").
        Select("COUNT(DISTINCT CASE WHEN graduation_status = true THEN students.id END) as graduated_count, "+
            "COUNT(DISTINCT CASE WHEN graduation_status = false THEN students.id END) as non_graduated_count").
        Joins("LEFT JOIN stcs ON students.stc_id = stcs.id").
        Where("students.mda_id = ? OR stcs.mda_id = ?", userMdaId, userMdaId).
        Scan(&summary).Error

        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "message": "error retrieving data case 1",
            })
            return
        }

        err = config.DB.Table("stcs").Where("mda_id = ?", userMdaId).Count(&summary.STCsCount).Error
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "message": "error retrieving data case ",
            })
            return
        }
		c.JSON(http.StatusOK,gin.H{"response":summary})



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
        type StcStats struct {
            TotalStudents     int
            TotalCertified    int
            TotalUncertified  int
        }
        
        var stcStats StcStats
        
        err = config.DB.Table("students").
            Select("COUNT(*) as total_students, SUM(CASE WHEN student_courses.is_certified = true THEN 1 ELSE 0 END) as total_certified, SUM(CASE WHEN student_courses.is_certified = false THEN 1 ELSE 0 END) as total_uncertified").
            Joins("JOIN student_courses ON students.id = student_courses.student_id").
            Where("students.stc_id = ?", userStcId).
            Scan(&stcStats).Error

        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "message": "Error trying to  this user",
            })
            return
        }
        c.JSON(http.StatusOK,gin.H{"response":stcStats})




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
        Select("courses.name AS course_name, COUNT(student_courses.id) AS total_students, SUM(CASE WHEN student_courses.is_certified = true THEN 1 ELSE 0 END) AS certified_count, SUM(CASE WHEN student_courses.is_certified = false THEN 1 ELSE 0 END) AS uncertified_count, COUNT(student_courses.id) * 100.0 / SUM(COUNT(student_courses.id)) OVER() AS total_percent").
        Joins("LEFT JOIN student_courses ON courses.id = student_courses.course_id").
        Group("courses.id").
        Order("total_percent DESC").
        Limit(5).
        Scan(&coursePercentages).Error

		if err!=nil {
            c.JSON(http.StatusBadRequest,gin.H{
                "message":"error retrieving info",
            })
            return
        }

        c.JSON(http.StatusOK,gin.H{"coursePercentages":coursePercentages})

    case 2:
        // Get user MDA ID
		var userMdaId uint 
		err := config.DB.Table("mdas").
		Where("user_id = ?", userID).
		Pluck("id",&userMdaId).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "MdaAccount has issues",
			})
			return
		}

        err = config.DB.Table("courses").
            Select("courses.name AS course_name, COUNT(student_courses.id) AS total_students, SUM(CASE WHEN student_courses.is_certified = true THEN 1 ELSE 0 END) AS certified_count,SUM(CASE WHEN student_courses.is_certified = false THEN 1 ELSE 0 END) AS uncertified_count, COUNT(student_courses.id) * 100.0 / SUM(COUNT(student_courses.id)) OVER() AS total_percent").
            Joins("LEFT JOIN student_courses ON courses.id = student_courses.course_id").
            Joins("LEFT JOIN students ON students.id = student_courses.student_id").
            Joins("LEFT JOIN stcs ON students.stc_id = stcs.id").
            Where("students.mda_id = ? OR stcs.mda_id = ?", userMdaId, userMdaId).
            Group("courses.id").
            Order("total_percent DESC").
            Limit(5).
            Scan(&coursePercentages).Error

        if err != nil{
            c.JSON(http.StatusInternalServerError,gin.H{
                "message":"error retrieving info",
            })
            return
        }
        c.JSON(http.StatusOK,gin.H{"coursePercentages":coursePercentages})

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
        err = config.DB.Table("courses").
            Select("courses.name AS course_name, COUNT(student_courses.id) AS total_students, SUM(CASE WHEN student_courses.is_certified = true THEN 1 ELSE 0 END) AS certified_count,SUM(CASE WHEN student_courses.is_certified = false THEN 1 ELSE 0 END) AS uncertified_count, COUNT(student_courses.id) * 100.0 / SUM(COUNT(student_courses.id)) OVER() AS total_percent").
            Joins("LEFT JOIN student_courses ON courses.id = student_courses.course_id").
            Joins("LEFT JOIN students ON students.id = student_courses.student_id").
            Where("students.stc_id = ?" ,userStcId).
            Group("courses.id").
            Order("total_percent DESC").
            Limit(5).
            Scan(&coursePercentages).Error

        if err != nil{
            c.JSON(http.StatusInternalServerError,gin.H{
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

func GetTopMda(c *gin.Context) { // only fme access
    type MdaWithStudentCount struct {
        MdaID   uint
        RegisterName    string
        StudentCount int
    }
    
    var topMdaWithStudents []MdaWithStudentCount

    err := config.DB.Table("mdas").
    Select("mdas.id as mda_id, mdas.register_name, COALESCE(SUM(direct_students.student_count), 0) + COALESCE(SUM(stc_students.student_count), 0) as student_count").
    Joins("LEFT JOIN (SELECT mda_id, COUNT(*) as student_count FROM students WHERE stc_id = 0 GROUP BY mda_id) AS direct_students ON mdas.id = direct_students.mda_id").
    Joins("LEFT JOIN (SELECT stcs.mda_id, COUNT(*) as student_count FROM students INNER JOIN stcs ON students.stc_id = stcs.id GROUP BY stcs.mda_id) AS stc_students ON mdas.id = stc_students.mda_id").
    Group("mdas.id, mdas.register_name").
    Order("student_count DESC").
    Limit(5).
    Scan(&topMdaWithStudents).Error

    if err != nil {
        c.JSON(http.StatusInternalServerError,gin.H{
            "message":"error retrieving info",
        })
        return
    }
    
    
    c.JSON(http.StatusOK,gin.H{
        "top5mdas":topMdaWithStudents,
    })
}


func GetTopStc(c *gin.Context) {
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
    var top5stcs []StcCount
    switch userRole {
    case 1:
        err := config.DB.Table("stcs").
    Select("stcs.name as stc_name, stcs.id as stc_id, COUNT(*) as total_students, SUM(CASE WHEN student_courses.is_certified = true THEN 1 ELSE 0 END) as total_certified").
    Joins("JOIN students ON stcs.id = students.stc_id").
    Joins("JOIN student_courses ON students.id = student_courses.student_id").
    Group("stcs.id").
    Order("total_students DESC").
    Limit(5).
    Scan(&top5stcs).Error

    if err != nil {
        c.JSON(http.StatusInternalServerError,gin.H{
            "message":"error retrieving info",
        })
        return
    }
    c.JSON(http.StatusOK,gin.H{"top5stcs":top5stcs})

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
        err = config.DB.Table("stcs").
        Select("stcs.name as stc_name, stcs.id as stc_id, COUNT(*) as total_students, SUM(CASE WHEN student_courses.is_certified = true THEN 1 ELSE 0 END) as total_certified").
        Joins("JOIN students ON stcs.id = students.stc_id").
        Joins("JOIN student_courses ON students.id = student_courses.student_id").
        Where("stcs.mda_id = ?", userMdaId).
        Group("stcs.id").
        Order("total_students DESC").
        Limit(5).
        Scan(&top5stcs).Error

    if err != nil {
        c.JSON(http.StatusInternalServerError,gin.H{
            "message":"error retrieving info",
        })
        return
    }
    c.JSON(http.StatusOK,gin.H{"top5stcs":top5stcs})


    default:
        fmt.Println(userID)
        c.JSON(http.StatusUnauthorized,gin.H{"message":"default unauthorized user"})
        return


    }
}