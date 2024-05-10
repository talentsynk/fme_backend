// package middleware

// import (
// 	"fme_backend/internal/config"
// 	"fmt"
// 	"net/http"
// 	"strings"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"github.com/golang-jwt/jwt/v4"
//  myuser	"fme_backend/internal/user"
// )

// // This is the middleware that accepts the authorization token and uses it to set the user id parameter if token is valid
// func RequireAuth(c *gin.Context) {
// 	authHeader:= c.Request.Header.Get("Authorization")
// 	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
// 		c.AbortWithStatus(http.StatusUnauthorized)
// 	}

// 	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
// 	secret := config.GetHashSecret()

// 	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 		// Don't forget to validate the alg is what you expect:
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
// 		}

// 		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
// 		return []byte(secret), nil
// 	})
// 	if token!= nil {
// 	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
// 		if float64(time.Now().Unix()) > claims["exp"].(float64) {
// 			c.AbortWithStatus(http.StatusUnauthorized)
// 		}

// 		var user myuser.User
// 		config.DB.First(&user,"email = ?", claims["sub"])

// 		if user.ID == 0 {
// 			c.AbortWithStatus(http.StatusUnauthorized)
// 		}
// 		c.Set("userId",user.ID)
// 		c.Next()
// 		return

// 	} else {
// 		c.AbortWithStatus(http.StatusUnauthorized)
// 	}}
// }

package middleware

import (
	"fme_backend/internal/config"
	mda "fme_backend/internal/mdas"
	"fme_backend/internal/stc"
	myuser "fme_backend/internal/user"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func RequireAuth(c *gin.Context) {
    authHeader := c.GetHeader("Authorization")
    if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
        c.Abort()
        return
    }

    tokenString := strings.TrimPrefix(authHeader, "Bearer ")
    secret := config.GetHashSecret()

    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(secret), nil
    })

    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization token"})
        c.Abort()
        return
    }

    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        if float64(time.Now().Unix()) > claims["exp"].(float64) {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
            c.Abort()
            return
        }

        var user myuser.User
        config.DB.First(&user, "email = ?", claims["sub"])

        if user.ID == 0 {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
            c.Abort()
            return
        }

        c.Set("userID", user.ID)
        c.Set("userRole", user.Role) // Ensure this key matches what you use in CreateMda
        c.Next()
    } else {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization token"})
        c.Abort()
    }
}

func RequireMda(c *gin.Context) {
    authHeader := c.GetHeader("Authorization")
    if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
        c.Abort()
        return
    }

    tokenString := strings.TrimPrefix(authHeader, "Bearer ")
    secret := config.GetHashSecret()

    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(secret), nil
    })

    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization token"})
        c.Abort()
        return
    }

    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        if float64(time.Now().Unix()) > claims["exp"].(float64) {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
            c.Abort()
            return
        }

        var user myuser.User
        config.DB.First(&user, "email = ?", claims["sub"])
        if user.ID == 0 {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
            c.Abort()
            return
        }

        if user.Role != 2 {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Not an Mda user 1"})
            c.Abort()
            return
        }

        var mda mda.Mda
        config.DB.First(&mda, "user_id = ?", user.ID)
        if mda.ID == 0 {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Not an Mda user 2"})
            c.Abort()
            return
        }

        c.Set("userID", user.ID) 
        c.Set("mdaID", mda.ID)
        c.Set("userRole", user.Role) // Ensure this key matches what you use in CreateMda
        c.Next()
    } else {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization token"})
        c.Abort()
    }

}

func RequireStc(c *gin.Context) {
    authHeader := c.GetHeader("Authorization")
    if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
        c.Abort()
        return
    }

    tokenString := strings.TrimPrefix(authHeader, "Bearer ")
    secret := config.GetHashSecret()

    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(secret), nil
    })

    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization token"})
        c.Abort()
        return
    }

    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        if float64(time.Now().Unix()) > claims["exp"].(float64) {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
            c.Abort()
            return
        }

        var user myuser.User
        config.DB.First(&user, "email = ?", claims["sub"])

        if user.ID == 0 {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
            c.Abort()
            return
        }

        if user.Role !=3{
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Not an Stc user"})
            c.Abort()
            return
        }

        var stc stc.Stc
        config.DB.First(&stc, "user_id = ?", user.ID)

        if stc.ID == 0 {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "No stc account user"})
            c.Abort()
            return
        }

        c.Set("userID", user.ID) 
        c.Set("stcID", stc.ID)
        c.Next()
    } else {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization token"})
        c.Abort()
    }
}

func RequireFme(c *gin.Context) {
    authHeader := c.GetHeader("Authorization")
    if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
        c.Abort()
        return
    }

    tokenString := strings.TrimPrefix(authHeader, "Bearer ")
    secret := config.GetHashSecret()

    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(secret), nil
    })

    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization token"})
        c.Abort()
        return
    }

    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        if float64(time.Now().Unix()) > claims["exp"].(float64) {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
            c.Abort()
            return
        }

        var user myuser.User
        config.DB.First(&user, "email = ?", claims["sub"])

        if user.ID == 0 {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
            c.Abort()
            return
        }

        if user.Role != 1 {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Not an FME"})
            c.Abort()
            return
        }

        c.Set("userID", user.ID) // Ensure this key matches what you use in CreateMda
        c.Next()
    } else {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization token"})
        c.Abort()
    }
}