package middleware

import (
    "time"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)



type RequestLog struct {
    ID           uint           `gorm:"primaryKey"`
    Method       string         `gorm:"size:10"` // GET, POST, etc.
    Path         string         `gorm:"size:255"`
    StatusCode   int            `gorm:"index"`
    ResponseTime float64        `gorm:"type:decimal(10,2)"` // time taken in milliseconds
    CreatedAt    time.Time
}


func RequestLogger(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Start timer
        startTime := time.Now()

        // Process request
        c.Next()

        // End timer
        endTime := time.Now()

        // Calculate response time
        responseTime := endTime.Sub(startTime).Milliseconds()

        // Create a log entry
        logEntry := RequestLog{
            Method:       c.Request.Method,
            Path:         c.Request.URL.Path,
            StatusCode:   c.Writer.Status(),
            ResponseTime: float64(responseTime),
            CreatedAt:    time.Now(),
        }

        // Save the log to the database
        if err := db.Create(&logEntry).Error; err != nil {
            // Optionally log the error or handle it
            c.Error(err)
        }
    }
}