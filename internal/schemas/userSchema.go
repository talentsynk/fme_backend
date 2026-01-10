package schemas

type UserCreateSchema struct {
	Email    string
	Password string
}
type LoginSchema struct {
	Email    string
	Password string
}

type RequestOtpSchema struct {
	Email string
}

type VerifyOtpSchema struct {
	Email string
	Otp   string
}

type ChangePasswordSchema struct {
	Email    string
	Password string
}

type LoginResp struct {
	Token string
	Role  string
}
