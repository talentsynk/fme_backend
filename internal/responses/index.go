package responses

import (
	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func ErrorResponse(c *fiber.Ctx, message string, statusCode int) error {
	return c.Status(statusCode).JSON(fiber.Map{
		"message": message,
	})
}

func SuccessResponse(c *fiber.Ctx, key string, data any, statusCode int) error {
	return c.Status(statusCode).JSON(fiber.Map{
		key: data,
	})
}

func CustomResponse(c *fiber.Ctx, statusCode int, data any) error {
	return c.Status(statusCode).JSON(data)
}

func LoginResponse(c *fiber.Ctx, token string, role string) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"jwt":        token,
		"token_type": "Bearer",
		"expires_in": 86400, // 24 hours in seconds
		"message":    "successful login",
		"role":       role,
	})
}

const (
	USER_EXISTS                = "user already exists"
	USER_CREATE_FAILED         = "failed to create user"
	USER_SUSPENDED_SUCCEFULLY  = "User suspended successfully"
	USER_SUSPENDED             = "suspended user"
	USER_SUSPEND_UNAUTHORIZED  = "unauthorized to suspend this user"
	USER_ACTIVATE_UNAUTHORIZED = "unauthorized to activate this user"
	UNABLE_TO_HASH_PASSWORD    = "unable to hash password"
	UNAUTHORIZED_ACCESS        = "unauthorized action"
	INCOMPLETE_DATA            = "incomplete data"
	BAD_DATA                   = "you provided the wrong data format for one or more of the fields"
	EMAIL_EXIST                = "email already in use"
	TELEPHONE_EXIST            = "telephone already in use"
	EMAIL_TELEPHONE_EXIST      = "email and telephone already in use"
	INVALID_NAME               = "invalid name format. please check your lastname and/or firstname"
	INVALID_LOGIN              = "invalid login details"
	INVALID_TELEPHONE          = "invalid telephone format"
	LOGIN_SUCCESS              = "login Successful"
	NO_ACCOUNT                 = "account not found"
	UNVERIFIED_TELEPHONE       = "telephone number is not verified"
	INVALID_COUNTRY            = "Country not supported"
	WRONG_TELEPHONE            = "Wrong telephone number"
	UPDATE_APP                 = "There is an important app update. Please update your app from your appStore/PlayStore"
	BAD_AUTHENTICATION         = "incorrect credentials"
	INVALID_OTP                = "invalid or expired Otp passed"
	TELEPHONE_VERIFIED         = "Telephone Successfully verified"
	USER_CREATED               = "Account Successfully Created"
	OTP_SENT                   = "OTP has been successfully sent to your telephone"
	PASSWORD_CHANGE_FAIL       = "unable to change password"
	PASSWORD_CHANGED           = "password changed successfully"
	WRONG_PASSWORD             = "wrong password passed"
	SOMETHING_WRONG            = "ooops! something went wrong. Please try again"
	INVALID_EMAIL_OTP          = "Wrong or expired email Otp passed"
	EMAIL_VERIFIED             = "Email successfully verified"
	BLOCKED_ACCOUNT            = "blocked account. please contact cashwise support via email at support@cashwise.finance"
	USER_FETCHED               = "User(s) fetched successfully"
	NO_USER                    = "no user found"
	SECRET_QUESTION_SET        = "2FA set successfully"
	SECRET_EXIST               = "2FA already exist for user"
	SECRET_UPDATED             = "2FA successfully updated"
	WRONG_SECRET_ANSWER        = "wrong secret answer passed"
	EMAIL_ALREADY_VERIFIED     = "Email address already verified"
	EMAIL_OTP_SENT             = "One-Time-Passord (OTP) sent to your email"
	REFRESH_TOKEN_ERROR        = "unable to generate token"
	CONTACTS_CHECKED           = "contacts successfully checked"
	COMPLETE_KYC               = "please complete your kyc"
	DATA_FETCHED               = "data successfully fetched"
	ADDRESS_ADDED              = "Address successfully added"
	MUST_BE_18                 = "User must be 18 years and above"
	KYC_LINK                   = "Kyc Link Generated"
	USER_ALREADY_VERIFIED      = "user has already been verified"
	ACCOUNT_AUTH_VERIFIED      = "user account auth has already been verified"
	PASSWORD_REUSE             = "you used this password recently. please choose a different one"
	VERIFICATION_SUCCESSFUL    = "Verification completed successfully"
	INVALID_ID_FORMAT          = "Invalid ID format"
	BVN_EXIST                  = "BVN already registered to someone else"
	USER_PENDING_APPROVAL      = "Your BVN verification is pending admin review."
	NIN_EXIST                  = "NIN already registered to someone else"
	USER_NIN_PENDING_APPROVAL  = "Your NIN verification is pending admin review."
	DATA_PROCESSED             = "data processed successfully"
	NO_SELFIE                  = "you are required to take a selfie before verifying your identity"
	ID_PENDING_APPROVAL        = "Your ID verification is pending admin review."
	USER_DELETED               = "Account Successfully Deleted"
	SHORT_ADDRESS              = "Address is too short. Please provide your full residential address"
)
