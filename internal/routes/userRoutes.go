package routes

import (
	"fme_backend/internal/config"
	"fme_backend/internal/controllers"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

var userController controllers.UserControllers

func UserRoutes(app *fiber.App) {
	userRoutes := app.Group("/user")
	userRoutes.Post("/createfme", userController.CreateFmeUser)
	userRoutes.Post("/login", userController.Login)
	userRoutes.Post("/otp/request", userController.RequestOtp)
	userRoutes.Post("/otp/verify", userController.VerifyOtp)
	userRoutes.Post("/changepassword", userController.ChangePassword)
	userRoutes.Get("/activate/:id", requireAuth, userController.ActivateUser)
	userRoutes.Get("/suspend/:id", requireAuth, userController.SuspendUser)

}

func requireAuth(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token required",
		})
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	secret := config.GetHashSecret()

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid authorization token",
		})
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Token expired",
			})
		}

		var userID uint
		if idFloat, ok := claims["user_id"].(float64); ok {
			userID = uint(idFloat)
		} else {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token claims"})
		}

		userRole, ok := claims["role"].(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token claims"})
		}

		email, _ := claims["sub"].(string)

		c.Locals("userID", userID)
		c.Locals("userRole", userRole)
		c.Locals("email", email)

		return c.Next()
	}

	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error": "Invalid authorization token",
	})
}
