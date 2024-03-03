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