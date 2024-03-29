package myuser

var UserCreateSchema struct {
	PhoneNumber string
	Email string
	Password string
}


var LoginSchema struct {
	Email string
	Password string
}

var RequestOtpSchema struct {
	Email string
}

var VerifyOtpSchema struct {
	Email string
	Otp string
}

var ChangePasswordSchema struct {
	Email string
	Password string
}