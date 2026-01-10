package controllers

import (
	"fme_backend/internal/config"
	"fme_backend/internal/responses"
	"fme_backend/internal/schemas"
	"fme_backend/internal/servers"
	"fme_backend/internal/utilities"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type UserControllers struct{}

var userServer servers.UserServer

func (u *UserControllers) CreateFmeUser(c *fiber.Ctx) error {
	var body schemas.UserCreateSchema
	if c.BodyParser(&body) != nil {
		return responses.ErrorResponse(c, responses.BAD_DATA, 400)
	}
	if err := userServer.CreateFmeUser(config.DB, body); err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}

	return responses.SuccessResponse(c, "message", responses.USER_CREATED, 200)
}

func (u *UserControllers) SuspendUser(c *fiber.Ctx) error {
	instance := c.Params("id")
	user := c.Get("userID")
	user_id, _ := strconv.Atoi(user)
	role := c.Get("role")
	if utilities.IsEmpty(instance) || utilities.IsEmpty(user_id) {
		responses.ErrorResponse(c, responses.INCOMPLETE_DATA, 400)
	}
	if err := userServer.SuspendUser(config.DB, user_id, instance, role); err != nil {
		responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, "message", responses.USER_SUSPENDED_SUCCEFULLY, 200)

}

func (u *UserControllers) ActivateUser(c *fiber.Ctx) error {
	instance := c.Params("id")
	user := c.Get("userID")
	user_id, _ := strconv.Atoi(user)
	role := c.Get("role")
	if utilities.IsEmpty(instance) || utilities.IsEmpty(user_id) {
		responses.ErrorResponse(c, responses.INCOMPLETE_DATA, 400)
	}
	if err := userServer.ActivateUser(config.DB, user_id, instance, role); err != nil {
		responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, "message", responses.USER_SUSPENDED_SUCCEFULLY, 200)

}

func (u *UserControllers) Login(c *fiber.Ctx) error {
	var body schemas.LoginSchema
	if c.BodyParser(&body) != nil {
		return responses.ErrorResponse(c, responses.BAD_DATA, 400)
	}
	if !utilities.IsEmail(body.Email) || utilities.IsEmpty(body.Password) {
		return responses.ErrorResponse(c, responses.INCOMPLETE_DATA, 400)
	}
	resp, err := userServer.Login(config.DB, body)
	if err != nil {
		return responses.ErrorResponse(c, responses.INCOMPLETE_DATA, 400)
	}
	return responses.LoginResponse(c, resp.Token, resp.Role)
}

func (u *UserControllers) RequestOtp(c *fiber.Ctx) error {
	var body schemas.RequestOtpSchema
	if c.BodyParser(&body) != nil {
		return responses.ErrorResponse(c, responses.BAD_DATA, 400)
	}
	if !utilities.IsEmail(body.Email) {
		return responses.ErrorResponse(c, responses.INCOMPLETE_DATA, 400)
	}
	if err := userServer.RequestOtp(config.DB, body); err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, "message", responses.EMAIL_OTP_SENT, 200)
}

func (u *UserControllers) VerifyOtp(c *fiber.Ctx) error {
	var body schemas.VerifyOtpSchema
	if c.BodyParser(&body) != nil {
		return responses.ErrorResponse(c, responses.BAD_DATA, 400)
	}
	if !utilities.IsEmail(body.Email) || utilities.IsEmpty(body.Otp) {
		return responses.ErrorResponse(c, responses.INCOMPLETE_DATA, 400)
	}
	if err := userServer.VerifyOtp(config.DB, body); err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, "message", responses.VERIFICATION_SUCCESSFUL, 200)
}

func (u *UserControllers) ChangePassword(c *fiber.Ctx) error {
	var body schemas.ChangePasswordSchema
	if c.BodyParser(&body) != nil {
		return responses.ErrorResponse(c, responses.BAD_DATA, 400)
	}
	if !utilities.IsEmail(body.Email) || utilities.IsValidPassword(body.Password) {
		return responses.ErrorResponse(c, responses.INCOMPLETE_DATA, 400)
	}
	if err := userServer.ChangePassword(config.DB, body); err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, "message", responses.PASSWORD_CHANGED, 200)
}
