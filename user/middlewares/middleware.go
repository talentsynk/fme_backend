package middleware

import (
	"fme_backend/config"
	"fmt"
	"net/http"
	"strings"
	"time"
     model "fme_backend/user/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)



// This is the middleware that accepts the authorization token and uses it to set the user id parameter if token is valid
func RequireAuth(c *gin.Context) {
	authHeader:= c.Request.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	secret := config.GetHashSecret()

	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
	
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secret), nil
	})
	if token!= nil {
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		
		var user model.User
		config.DB.First(&user,"email = ?", claims["sub"])

		if user.ID == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		c.Set("userId",user.ID)
		c.Next()
		return
		
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}}	
}